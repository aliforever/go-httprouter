package httprouter

import "net/http"

type Controller interface {
	GET(http.ResponseWriter, *http.Request)
	POST(http.ResponseWriter, *http.Request)
	PATCH(http.ResponseWriter, *http.Request)
	PUT(http.ResponseWriter, *http.Request)
	DELETE(http.ResponseWriter, *http.Request)
}
