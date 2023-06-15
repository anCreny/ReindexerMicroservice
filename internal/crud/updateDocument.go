package crud

import (
	"github.com/anCreny/ReindexerMicroservice/internal"
	"encoding/json"
	"github.com/restream/reindexer/v3"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

// UpdateDocument http://localhost/updatedocument
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var body, readErr = io.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(400)
		panic(readErr)
	}

	var requestDocument internal.DocumentAJson
	if err := json.Unmarshal(body, &requestDocument); err != nil {
		w.WriteHeader(400)
		panic(err)
	}

	query := db.Query("DocumentsA").Where("id", reindexer.EQ, requestDocument.ID).Limit(1).
		Join(db.Query("DocumentsB"), "documents_B_list").
		On("documents_B_ids", reindexer.EQ, "id").
		Sort("sort", true)

	if qResult, err := query.Exec().FetchOne(); err == nil {
		var doc = *qResult.(*internal.DocumentA)
		var AId = doc.ID
		var docsB = doc.DocumentsBList
		var newDocsBJson = requestDocument.DocumentsBList

		deleteDocA(strconv.Itoa(doc.ID))

		var newDocsB = []internal.DocumentB{}
		var newDocsBIds = []int{}

		for i := 0; i < len(newDocsBJson); i++ {
			var newId int
			if i+1 > len(docsB) {
				newId = getUniqId("DocumentsB")
			} else {
				newId = docsB[i].ID
			}

			newDocsB = append(newDocsB, internal.DocumentB{
				ID:             newId,
				DocumentsCList: convertJsonsToDocsC(newDocsBJson[i].DocumentsCList),
				Sort:           rand.Intn(len(newDocsBJson)),
			})

			newDocsBIds = append(newDocsBIds, newId)
		}

		for _, newDocB := range newDocsB {
			if status, insErr := db.Insert("DocumentsB", newDocB); insErr != nil && status == 0 {
				panic(insErr)
			}
		}

		if status, insErr := db.Insert("DocumentsA", internal.DocumentA{
			ID:             AId,
			DocumentsBList: nil,
			DocumentsB_IDs: newDocsBIds,
		}); insErr != nil && status == 0 {
			panic(insErr)
		}

		cacheQuery := db.Query("DocumentsA").
			Where("id", reindexer.EQ, AId).Limit(1).
			Join(db.Query("DocumentsB"), "documents_B_list").
			On("documents_B_ids", reindexer.EQ, "id").
			Sort("sort", true)

		if result, cacheErr := cacheQuery.Exec().FetchOne(); cacheErr == nil {
			var docA = *result.(*internal.DocumentA)
			var respChan = make(chan internal.DocumentAJson)

			go processDocument(docA, respChan, nil)

			var docAJson = <-respChan

			tryUpdateCachedDocument(docAJson)
		} else {
			panic(cacheErr)
		}

	} else {
		w.WriteHeader(400)
		panic(err)
	}

	w.WriteHeader(200)

}

func convertJsonsToDocsC(documentsCJson []internal.DocumentCJson) []internal.DocumentC {
	var result = []internal.DocumentC{}

	for _, documentCJson := range documentsCJson {
		result = append(result, internal.DocumentC{Text: documentCJson.Text})
	}

	return result
}

func getUniqId(namespace string) int {
	var result int

	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	query := db.Query(namespace)

	if collection, err := query.Exec().FetchAll(); err == nil {
		result = len(collection) + 1
		for {
			foundQ := db.Query(namespace).Where("id", reindexer.EQ, result)
			if _, foundErr := foundQ.Exec().FetchOne(); foundErr != nil {
				break
			}
			result += 1
		}
	} else {
		panic(err)
	}

	return result
}
