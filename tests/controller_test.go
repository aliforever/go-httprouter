package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aliforever/go-httprouter"
)

type HelloController struct{}

func (c HelloController) Path() string {
	return "/hello"
}

func (c HelloController) GET(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("GET Handler Called"))
}

func (c HelloController) POST(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("POST Handler Called"))
}

func (c HelloController) PATCH(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("PATCH Handler Called"))
}

func (c HelloController) PUT(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("PUT Handler Called"))
}

func (c HelloController) DELETE(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("DELETE Handler Called"))
}

func TestController(t *testing.T) {
	var router httprouter.Router
	router.Controller(HelloController{})

	go func() {
		err := http.ListenAndServe("localhost:8080", router)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	time.Sleep(time.Second * 5)

	data := []struct {
		Method   string
		Response string
	}{
		{
			Method:   "GET",
			Response: "GET Handler Called",
		},
		{
			Method:   "POST",
			Response: "POST Handler Called",
		},
		{
			Method:   "PUT",
			Response: "PUT Handler Called",
		},
		{
			Method:   "PATCH",
			Response: "PATCH Handler Called",
		},
		{
			Method:   "DELETE",
			Response: "DELETE Handler Called",
		},
	}
	for _, datum := range data {
		response, err := testController(HelloController{}, datum.Method)
		if err != nil {
			t.Errorf("Error testing route: %s => %s", HelloController{}.Path(), err)
			continue
		}
		if response != datum.Response {
			t.Errorf("Test Failed for %s.\nExpected: %s\nActual: %s", HelloController{}.Path(), datum.Response, response)
			continue
		}
	}
}

func testController(controller httprouter.Controller, method string) (response string, err error) {
	var req *http.Request
	req, err = http.NewRequest(method, "http://localhost:8080"+controller.Path(), nil)
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = string(b)
	return
}
