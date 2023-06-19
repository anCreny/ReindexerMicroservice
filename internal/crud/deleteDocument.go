package crud

import (
	"github.com/anCreny/ReindexerMicroservice/internal"
	"github.com/restream/reindexer/v3"
	"net/http"
	"strconv"
)

// DeleteDocument http://localhost/deletedocument?id='int'
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	var db, initErr = internal.Database()
	if initErr != nil {
		panic(initErr)
	}

	var delId = r.URL.Query().Get("id")

	var delQuery = db.Query("Documents").Where("id", reindexer.EQ, delId)

	if _, delErr := delQuery.Delete(); delErr != nil {
		panic(delErr)
	}

	if id, err := strconv.Atoi(delId); err == nil {
		deleteCachedDocument(id)
	} else {
		panic(err)
	}

	w.WriteHeader(200)
}
