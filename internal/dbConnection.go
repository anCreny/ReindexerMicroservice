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

	if _, found := db.Query("DocumentsA").Get(); !found {
		if err := fillNamespace("Documents"); err != nil {
			return err
		}
	}

	return nil
}

func fillNamespace(namespace string) error {
	var documentsB_Ids = []int{}

	var documentsC = []DocumentC{}
	for i := 0; i < 10; i++ {
		var docC = DocumentC{
			Text: "Some text",
		}
		documentsC = append(documentsC, docC)
	}

	for i := 0; i < 100; i++ {
		var docB = &DocumentB{
			ID:             i,
			Sort:           rand.Intn(100),
			DocumentsCList: documentsC,
		}
		if err := databaseInstance.Upsert(namespace+"B", docB); err != nil {
			return err
		}
		documentsB_Ids = append(documentsB_Ids, i)
	}

	var docAId int
	var tempDocsBIds = []int{}

	for i := 0; i < 100; i++ {

		if i%10 == 0 && i != 0 {
			if err := databaseInstance.Upsert(namespace+"A", &DocumentA{
				ID:             docAId,
				DocumentsB_IDs: tempDocsBIds,
			}); err != nil {
				return err
			}
			docAId++
			tempDocsBIds = []int{}
		}
		tempDocsBIds = append(tempDocsBIds, documentsB_Ids[i])
	}

	return nil
}

type DocumentA struct {
	ID             int          `reindex:"id,,pk" json:"ID"`
	DocumentsBList []*DocumentB `reindex:"documents_B_list,,joined"`
	DocumentsB_IDs []int        `reindex:"documents_B_ids"`
}

type DocumentB struct {
	ID             int         `reindex:"id,,pk"`
	DocumentsCList []DocumentC `reindex:"documents_C_list,,joined"`
	Sort           int         `reindex:"sort,tree"`
}

type DocumentC struct {
	Text string `reindex:"text"`
}

type DocumentAJson struct {
	ID             int             `json:"id"`
	DocumentsBList []DocumentBJson `json:"documents-b-list"`
}

type DocumentBJson struct {
	DocumentsCList []DocumentCJson `json:"documents-c-list"`
}

type DocumentCJson struct {
	Text string `json:"text"`
}
