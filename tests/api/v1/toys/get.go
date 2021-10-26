package main

import (
	"fmt"
	"net/http"

	"github.com/aliforever/go-httprouter"
)

func main() {
	var router httprouter.Router

	err := router.Register(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))
	if err != nil {
		fmt.Println(err)
		return
	}

}
