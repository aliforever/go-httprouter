package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aliforever/go-httprouter"
)

func TestDelegate(t *testing.T) {
	var router httprouter.Router

	router.Delegate(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("Delegate Handler Called for Path: %s", request.URL.Path)))
	}), "/images/")

	go func() {
		err := http.ListenAndServe("localhost:8080", router)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	time.Sleep(time.Second * 2)

	data := []struct {
		Path     string
		Response string
	}{
		{
			Path:     "/images/logo.png",
			Response: "Delegate Handler Called for Path: /images/logo.png",
		},
		{
			Path:     "/images/slightly-smiling.png",
			Response: "Delegate Handler Called for Path: /images/slightly-smiling.png",
		},
		{
			Path:     "/images/team/joe.jpeg",
			Response: "Delegate Handler Called for Path: /images/team/joe.jpeg",
		},
		{
			Path:     "/images/team/marry.jpeg",
			Response: "Delegate Handler Called for Path: /images/team/marry.jpeg",
		},
		{
			Path:     "/images/webbug/1x1.gif",
			Response: "Delegate Handler Called for Path: /images/webbug/1x1.gif",
		},
		{
			Path:     "/image/webbug/1x1.gif",
			Response: "404 page not found\n",
		},
	}
	for _, datum := range data {
		response, err := testDelegate(datum.Path)
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

func testDelegate(path string) (response string, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", "http://localhost:8080"+path, nil)
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
