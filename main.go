package main

import (
	"github.com/anCreny/ReindexerMicroservice/cmd"
	"github.com/anCreny/ReindexerMicroservice/internal"
	"github.com/anCreny/ReindexerMicroservice/internal/crud"
	"fmt"
)

func main() {
	if err := internal.InitDbConnection(); err != nil {
		fmt.Println(err)
		return
	}

	cmd.AddHandler("/getdocuments", crud.ReadDocuments)     // http://localhost/getdocument
	cmd.AddHandler("/getonedocument", crud.ReadOneDocument) // http://localhost/getonedocument?id='number'

	cmd.AddHandler("/deletedocument", crud.DeleteDocument) // http://localhost/deletedocument?id='number'

	cmd.AddHandler("/createdocument", crud.CreateDocument) // http://localhost/createdocument

	cmd.AddHandler("/updatedocument", crud.UpdateDocument) // http://localhost/updatedocument

	cmd.RunMicroservice()
}
