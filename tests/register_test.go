package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aliforever/go-httprouter"
)

func TestRegister(t *testing.T) {
	var router httprouter.Router

	err := router.Register(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))
	if err != nil {
		fmt.Println(err)
		return
	}

}
