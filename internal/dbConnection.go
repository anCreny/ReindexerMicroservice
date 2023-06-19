package internal

import (
	"github.com/restream/reindexer/v3"
	"math/rand"
)
import _ "github.com/restream/reindexer/v3/bindings/cproto"

var databaseInstance *reindexer.Reindexer = nil

type DbInitError struct {
	errorMessage string
}

func (d *DbInitError) Error() string {
	return d.errorMessage
}

func Database() (*reindexer.Reindexer, error) {
	if databaseInstance != nil {
		return databaseInstance, nil
	}

	return nil, &DbInitError{errorMessage: "Database is not initialized yet"}
}

func InitDbConnection() error {
	var dbPort = Db_port()
	var dbUsername = Db_username()
	var dbPassword = Db_password()
	var dbName = Db_name()

	var connectionString = "cproto://" + dbUsername + ":" + dbPassword + "@reindexer_db:" + dbPort + "/" + dbName

	var db = reindexer.NewReindex(connectionString, reindexer.WithCreateDBIfMissing())
	if err := db.Ping(); err != nil {
		return err
	}

	databaseInstance = db

	if err := db.OpenNamespace("Documents", reindexer.DefaultNamespaceOptions(), Document{}); err != nil {
		return err
	}

	if _, found := db.Query("Documents").Get(); !found {
		if err := fillNamespace("Documents"); err != nil {
			return err
		}
	}

	return nil
}

func fillNamespace(namespace string) error {

	for i := 0; i < 100; i++ {
		var DocumentsB = []DocumentB{}

		for j := 0; j < 100; j++ {
			DocumentsB = append(DocumentsB, DocumentB{
				Title:  "Random title",
				Text:   "Random text",
				Author: "Random author",
			})
		}

		if err := databaseInstance.Upsert(namespace, Document{
			DocumentsBList: DocumentsB,
			Sort:           rand.Intn(100),
		}, "id=serial()"); err != nil {
			return err
		}
	}

	return nil
}

type Document struct {
	ID             int         `reindex:"id,,pk"`
	DocumentsBList []DocumentB `reindex:"documents_B_list"`
	Sort           int         `reindex:"sort,tree"`
}

type DocumentB struct {
	Title  string `reindex:"title" json:"title"`
	Text   string `reindex:"text" json:"text"`
	Author string `reindex:"author" json:"author"`
}

type DocumentJson struct {
	ID             int         `json:"id"`
	DocumentsBList []DocumentB `json:"documents-b-list"`
}
