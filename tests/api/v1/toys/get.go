package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aliforever/go-httprouter"
)

func main() {
	var router httprouter.Router

	expected := "This is toys handler"

	err := router.Register(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(expected))
	}))
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		err := http.ListenAndServe("localhost:8080", router)
		if err != nil {
			return
		}
	}()

	time.Sleep(time.Second * 2)

	resp, err := testToys()
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp != expected {
		fmt.Printf("Test Failed.\nExpected: %s\nActual: %s", expected, resp)
		return
	}

	fmt.Printf("Test Passed.")
}

func testToys() (response string, err error) {
	var resp *http.Response
	resp, err = http.Get("http://localhost:8080/v1/toys")
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
