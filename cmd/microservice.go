package cmd

import (
	"RMicroService/internal"
	"fmt"
	"net/http"
)

func AddHandler(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func RunMicroservice() {
	var port = internal.Port()
	fmt.Println("Microservice now listening on :" + port)
	panic(http.ListenAndServe(":"+port, nil))
}
