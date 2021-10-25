package httprouter

import (
	"errors"
	"net/http"
)

type pathRouter struct {
	methodHandler  map[string]http.Handler
	defaultHandler http.Handler
}

func (pr pathRouter) handlerByMethod(method string) (handler http.Handler, err error) {
	if pr.methodHandler == nil && pr.defaultHandler == nil {
		err = errors.New("not_found")
		return
	}

	if pr.methodHandler == nil {
		return pr.defaultHandler, nil
	}

	for meth, handler := range pr.methodHandler {
		if meth == method {
			return handler, nil
		}
	}

	return nil, errors.New("not_found")
}

func (pr *pathRouter) registerDefaultHandler(handler http.Handler) {
	pr.defaultHandler = handler
}

func (pr *pathRouter) registerMethodHandler(method string, handler http.Handler) {
	if pr.methodHandler == nil {
		pr.methodHandler = map[string]http.Handler{}
	}
	pr.methodHandler[method] = handler
}
