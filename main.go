package main

import (
	"RMicroService/cmd"
	"RMicroService/internal"
	"RMicroService/internal/crud"
	"fmt"
)

func main() {
	if err := internal.InitDbConnection(); err != nil {
		fmt.Println(err)
	}
	cmd.AddHandler("/getdocuments", crud.ReadDocuments)
	cmd.AddHandler("/getonedocument", crud.ReadOneDocument)
	cmd.AddHandler("/deletedocument", crud.DeleteDocument)
	cmd.RunMicroservice()
}
