package main

import (
	"RMicroService/cmd"
	"RMicroService/internal"
	"fmt"
)

func main() {
	if err := internal.InitDbConnection(); err != nil {
		fmt.Println(err)
	}

	cmd.RunMicroservice()
}
