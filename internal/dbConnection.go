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

	if err := db.OpenNamespace("DocumentsA", reindexer.DefaultNamespaceOptions(), DocumentA{}); err != nil {
		return err
	}

	if err := db.OpenNamespace("DocumentsB", reindexer.DefaultNamespaceOptions(), DocumentB{}); err != nil {
		return err
	}

	if err := db.OpenNamespace("DocumentsC", reindexer.DefaultNamespaceOptions(), DocumentC{}); err != nil {
		return err
	}

	if _, found := db.Query("DocumentsA").Get(); !found {
		if err := fillNamespace("Documents"); err != nil {
			return err
		}
	}

	return nil
}

func fillNamespace(namespace string) error {
	var documentsC = []*DocumentC{}
	var documentsB = []*DocumentB{}

	for i := 0; i < 100; i++ {
		var docC = &DocumentC{
			ID:   i,
			Text: "Some text",
		}
		if err := databaseInstance.Upsert(namespace+"C", docC); err != nil {
			return err
		}
		documentsC = append(documentsC, docC)
	}

	for i := 0; i < 100; i++ {
		var docB = &DocumentB{
			ID:         i,
			DocumentsC: documentsC,
			Sort:       rand.Int(),
		}
		if err := databaseInstance.Upsert(namespace+"B", docB); err != nil {
			return err
		}
		documentsB = append(documentsB, docB)
	}

	for i := 0; i < 100; i++ {
		if err := databaseInstance.Upsert(namespace+"A", &DocumentA{
			ID:         i,
			DocumentsB: documentsB,
		}); err != nil {
			return err
		}
	}
	return nil
}

type DocumentA struct {
	ID         int          `reindex:"id,,pk" json:"ID"`
	DocumentsB []*DocumentB `reindex:"documentsB,,joined" json:"documentsb"`
}

type DocumentB struct {
	ID         int          `reindex:"id,,pk"`
	DocumentsC []*DocumentC `reindex:"documentsC,,joined" json:"documentsc"`
	Sort       int          `reindex:"sort,tree"`
}

type DocumentC struct {
	ID   int    `reindex:"id,,pk"`
	Text string `reindex:"text" json:"text"`
}
