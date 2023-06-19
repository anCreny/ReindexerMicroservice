package crud

import (
	"encoding/json"
	"github.com/anCreny/ReindexerMicroservice/internal"
	"github.com/restream/reindexer/v3"
	"io"
	"net/http"
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

	var requestDocument internal.DocumentJson
	if err := json.Unmarshal(body, &requestDocument); err != nil {
		w.WriteHeader(400)
		panic(err)
	}

	updateQuery := db.Query("Documents").Where("id", reindexer.EQ, requestDocument.ID).Set("documents_B_list", requestDocument.DocumentsBList)

	if updateErr := updateQuery.Update(); updateErr != nil {
		w.WriteHeader(400)
		panic(updateErr)
	}

	updateCachedDocumentIfExists(requestDocument)
	w.WriteHeader(200)

}
