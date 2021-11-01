package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aliforever/go-httprouter"
)

func TestRegisterPath(t *testing.T) {
	data := []struct {
		route          string
		method         string
		allowedMethods []string
		response       string
		handler        func(w http.ResponseWriter, r *http.Request)
	}{
		{
			route:          "/my",
			allowedMethods: []string{"GET", "POST"},
			method:         "GET",
			response:       "This is my page",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("This is my page"))
			},
		},
		{
			route:          "/your",
			allowedMethods: nil,
			method:         "DELETE",
			response:       "This is your page",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("This is your page"))
			},
		},
		{
			route:          "/hello",
			allowedMethods: []string{"PATCH", "GET"},
			method:         "GET",
			response:       "Hello World!",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World!"))
			},
		},
		{
			route:          "/bye",
			allowedMethods: []string{"PATCH", "DELETE", "POST"},
			method:         "GET",
			response:       "404 page not found\n",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("404 page not found\n"))
			},
		},
	}

	var router httprouter.Router

	for i := range data {
		router.RegisterPathMethods(http.HandlerFunc(data[i].handler), data[i].route, data[i].allowedMethods...)
	}

	go func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	time.Sleep(time.Second * 5)

	for _, datum := range data {
		response, err := testRoute(datum.route, datum.method)
		if err != nil {
			t.Errorf("Error testing route: %s => %s", datum.route, err)
			continue
		}
		if response != datum.response {
			t.Errorf("Test Failed for %s.\nExpected: %s\nActual: %s", datum.route, datum.response, response)
			continue
		}
	}
}

func testRoute(route string, method string) (response string, err error) {
	var req *http.Request
	req, err = http.NewRequest(method, "http://localhost:8080"+route, nil)
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
