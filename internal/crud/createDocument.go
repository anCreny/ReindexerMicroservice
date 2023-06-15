package crud

import (
	"github.com/anCreny/ReindexerMicroservice/internal"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

// CreateDocument http://localhost/createdocument
func CreateDocument(w http.ResponseWriter, r *http.Request) {
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(db)
	}

	var body, bodyErr = io.ReadAll(r.Body)
	if bodyErr != nil {
		w.WriteHeader(400)
		panic(bodyErr)
	}

	var reqDocument internal.DocumentAJson
	if unmErr := json.Unmarshal(body, &reqDocument); unmErr != nil {
		w.WriteHeader(400)
		panic(unmErr)
	}

	var newDocA internal.DocumentA
	var newDocsB = []internal.DocumentB{}
	var newDocsBIds = []int{}

	for _, docBJson := range reqDocument.DocumentsBList {
		var newDocB internal.DocumentB
		newDocB = internal.DocumentB{
			ID:             getUniqId("DocumentsB"),
			DocumentsCList: convertJsonsToDocsC(docBJson.DocumentsCList),
			Sort:           rand.Intn(len(reqDocument.DocumentsBList)),
		}
		newDocsB = append(newDocsB, newDocB)
		newDocsBIds = append(newDocsBIds, newDocB.ID)
	}

	for _, newDocB := range newDocsB {
		if status, err := db.Insert("DocumentsB", newDocB); status == 0 && err != nil {
			panic(err)
		}
	}

	var newAId = getUniqId("DocumentsA")

	newDocA = internal.DocumentA{
		ID:             newAId,
		DocumentsBList: nil,
		DocumentsB_IDs: newDocsBIds,
	}

	if status, err := db.Insert("DocumentsA", newDocA); status == 0 && err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}
