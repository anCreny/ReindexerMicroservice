package crud

import (
	"encoding/json"
	"github.com/anCreny/ReindexerMicroservice/internal"
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

	var reqDocument internal.DocumentJson
	if unmErr := json.Unmarshal(body, &reqDocument); unmErr != nil {
		w.WriteHeader(400)
		panic(unmErr)
	}

	var newDoc = convertFromJsonToDoc(reqDocument)

	if status, err := db.Insert("Documents", newDoc, "id=serial()"); status == 0 && err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func convertFromJsonToDoc(docJson internal.DocumentJson) internal.Document {
	var result = internal.Document{
		ID:             docJson.ID,
		DocumentsBList: docJson.DocumentsBList,
		Sort:           rand.Intn(100),
	}

	return result
}
