package crud

import (
	"encoding/json"
	"github.com/anCreny/ReindexerMicroservice/internal"
	"github.com/restream/reindexer/v3"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var cachedDocuments = make(map[int]cachedDocument)

type cachedDocument struct {
	document internal.DocumentJson
	timer    *time.Timer
}

func cacheDocument(document internal.DocumentJson) {
	if cachedDoc, found := cachedDocuments[document.ID]; found {
		cachedDoc.timer.Stop()
	}
	var timeout = time.NewTimer(15 * time.Minute)
	cachedDocuments[document.ID] = cachedDocument{document, timeout}
	go func() {
		<-timeout.C
		deleteCachedDocument(document.ID)
	}()
}

func updateCachedDocumentIfExists(document internal.DocumentJson) {
	if value, found := cachedDocuments[document.ID]; found {
		cachedDocuments[document.ID] = cachedDocument{document, value.timer}
	}
}

func deleteCachedDocument(id int) {
	delete(cachedDocuments, id)
}

func tryGetCachedDocument(id int) (internal.DocumentJson, bool) {
	if value, found := cachedDocuments[id]; found {
		return value.document, true
	}

	return internal.DocumentJson{}, false
}

// ReadOneDocument http://localhost/getonedocument?id='int'
func ReadOneDocument(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var intId, err = strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	if value, found := tryGetCachedDocument(intId); found {
		var response = value

		if jsonResponse, respErr := json.Marshal(response); respErr == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		} else {
			panic(err)
		}

	} else {
		var readQuery = db.Query("DocumentsA").
			Where("id", reindexer.EQ, id).Limit(1)

		if result, readErr := readQuery.Exec().FetchOne(); readErr == nil {

			var response = convertFromDocToJson(result.(internal.Document))

			if jsonResponse, respErr := json.Marshal(response); respErr == nil {
				cacheDocument(response)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)
			} else {
				panic(err)
			}

		} else {
			w.WriteHeader(404)
			return
		}
	}
}

// ReadDocuments http://localhost/getdocuments
func ReadDocuments(w http.ResponseWriter, r *http.Request) {
	var pageNumberStr = r.URL.Query().Get("page")
	var limitStr = r.URL.Query().Get("limit")

	var pageNumber, _ = strconv.Atoi(pageNumberStr)
	var limit, _ = strconv.Atoi(limitStr)

	var offset int

	pageNumber--
	if pageNumber <= 0 {
		pageNumber = 1
		offset = 0
	} else {
		offset = limit * pageNumber
	}

	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var readQuery *reindexer.Query

	if limit > 0 {
		readQuery = db.Query("Documents").Offset(offset).Limit(limit).Sort("sort", true)
	} else {
		readQuery = db.Query("Documents").Sort("sort", true)
	}

	if qResult, readErr := readQuery.Exec().FetchAll(); readErr == nil {

		var length = len(qResult)

		var docsJson = make([]internal.DocumentJson, length)
		var waitGroup = sync.WaitGroup{}

		waitGroup.Add(length)
		for i, elem := range qResult {
			go func(index int, doc internal.Document) {
				defer waitGroup.Done()
				docsJson[index] = convertFromDocToJson(doc)
			}(i, *elem.(*internal.Document))
		}
		waitGroup.Wait()

		if response, marshErr := json.Marshal(docsJson); marshErr == nil {
			if _, wErr := w.Write(response); wErr != nil {
				w.WriteHeader(404)
				panic(wErr)
			}
		} else {
			panic(marshErr)
		}

	} else {
		w.WriteHeader(404)
		panic(readErr)
	}
}

func convertFromDocToJson(document internal.Document) internal.DocumentJson {
	var result = internal.DocumentJson{
		ID:             document.ID,
		DocumentsBList: document.DocumentsBList,
	}

	return result
}
