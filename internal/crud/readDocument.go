package crud

import (
	"RMicroService/internal"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var cachedDocuments = make(map[int]cachedDocument)

type cachedDocument struct {
	document internal.DocumentAJson
	timer    *time.Timer
}

func cacheDocument(document internal.DocumentAJson) {
	if cachedDoc, found := cachedDocuments[document.ID]; found {
		cachedDoc.timer.Stop()
	}
	var timeout = time.NewTimer(15 * time.Minute)
	cachedDocuments[document.ID] = cachedDocument{document, timeout}
	go func() {
		<-timeout.C
		tryDeleteCachedDocument(document.ID)
	}()
}

func tryUpdateCachedDocument(document internal.DocumentAJson) {
	if value, found := cachedDocuments[document.ID]; found {
		cachedDocuments[document.ID] = cachedDocument{document, value.timer}
	}
}

func tryDeleteCachedDocument(id int) {
	delete(cachedDocuments, id)
}

func tryGetCachedDocument(id int) (internal.DocumentAJson, bool) {
	if value, found := cachedDocuments[id]; found {
		return value.document, true
	}

	return internal.DocumentAJson{}, false
}

// ReadOneDocument http://localhost/getonedocument?id='int'
func ReadOneDocument(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var response internal.DocumentAJson

	var intId, err = strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	var isCached bool

	if value, found := tryGetCachedDocument(intId); found {
		response = value
		isCached = true
	} else {
		var query = db.Query("DocumentsA").
			Where("id", reindexer.EQ, id).Limit(1).
			Join(db.Query("DocumentsB"), "documents_B_list").
			On("documents_B_ids", reindexer.EQ, "id").
			Sort("sort", true)

		if result, err1 := query.Exec().FetchOne(); err1 == nil {

			var doc = *result.(*internal.DocumentA)

			var resChan = make(chan internal.DocumentAJson)
			go processDocument(doc, resChan, nil)
			response = <-resChan

		} else {
			w.WriteHeader(404)
		}
	}

	if jsonResponse, respErr := json.Marshal(response); respErr == nil {
		if !isCached {
			cacheDocument(response)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		panic(err)
	}

}

// ReadDocuments http://localhost/getdocuments
func ReadDocuments(w http.ResponseWriter, r *http.Request) {
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	queryA := db.Query("DocumentsA").
		Join(db.Query("DocumentsB"), "documents_B_list").
		On("documents_B_ids", reindexer.EQ, "id").
		Sort("sort", true)

	if documents, err2 := queryA.Exec().FetchAll(); err2 == nil {
		var length = len(documents)
		var group sync.WaitGroup
		var resultChan = make(chan internal.DocumentAJson, length)
		group.Add(length)
		for _, value := range documents {
			go processDocument(*value.(*internal.DocumentA), resultChan, &group)
		}
		group.Wait()
		close(resultChan)
		var response = []internal.DocumentAJson{}
		for formattedDocument := range resultChan {
			response = append(response, formattedDocument)
		}

		if jsonResponse, err := json.Marshal(response); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		} else {
			panic(err)
		}

	} else {
		fmt.Fprint(w, err2)
	}
}

func processDocument(documentA internal.DocumentA, output chan internal.DocumentAJson, group *sync.WaitGroup) {
	defer func() {
		if group != nil {
			group.Done()
		}
	}()

	var documentAJson = internal.DocumentAJson{ID: documentA.ID}

	for _, value := range documentA.DocumentsBList {
		var documentsCJsonTemp = []internal.DocumentCJson{}
		for _, value2 := range value.DocumentsCList {
			documentsCJsonTemp = append(documentsCJsonTemp, internal.DocumentCJson{Text: value2.Text})
		}
		documentAJson.DocumentsBList = append(documentAJson.DocumentsBList, internal.DocumentBJson{documentsCJsonTemp})
	}

	output <- documentAJson
}
