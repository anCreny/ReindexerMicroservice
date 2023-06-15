package crud

import (
	"RMicroService/internal"
	"github.com/restream/reindexer/v3"
	"net/http"
	"strconv"
)

func DeleteDocument(w http.ResponseWriter, r *http.Request) {

	var id = r.URL.Query().Get("id")

	deleteDocA(id)

	if id, err := strconv.Atoi(id); err == nil {
		tryDeleteCachedDocument(id)
	} else {
		panic(err)
	}

	w.WriteHeader(200)
}

func deleteDocA(id string) {
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var deleteId = id

	query := db.Query("DocumentsA").
		Where("id", reindexer.EQ, deleteId).Limit(1)

	if result, err := query.Exec().FetchOne(); err == nil {
		var doc = *result.(*internal.DocumentA)
		for _, docBId := range doc.DocumentsB_IDs {
			if _, err := db.Query("DocumentsB").Where("id", reindexer.EQ, docBId).Delete(); err != nil {
				panic(err)
			}
		}
	} else {
		panic(err)
	}

	deleteQuery := db.Query("DocumentsA").Where("id", reindexer.EQ, id)

	if _, err := deleteQuery.Delete(); err != nil {
		panic(err)
	}

}
