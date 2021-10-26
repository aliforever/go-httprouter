package tests

import (
	"net/http"
	"testing"

	"github.com/aliforever/go-httprouter"
)

func TestRegister(t *testing.T) {
	var router httprouter.Router

	router.Register(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))

}
