package httprouter

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"
)

type Router struct {
	routers        map[string]pathRouter
	delegates      map[string]http.Handler
	defaultHandler map[string]http.Handler
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if router.delegates != nil {
		for key, handler := range router.delegates {
			if strings.HasPrefix(path, key) {
				handler.ServeHTTP(w, r)
				return
			}
		}
	}

	if h, ok := router.routers[path]; ok {
		handler, err := h.handlerByMethod(r.Method)
		if err == nil {
			handler.ServeHTTP(w, r)
			return
		}
	}

	if defaultHandler, ok := router.defaultHandler[path]; ok {
		defaultHandler.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)

	return
}

func (router *Router) Register(handler http.Handler, methods ...string) (err error) {
	// TODO: Will complete this part
	_, file, line, ok := runtime.Caller(1)
	fmt.Println(file, path.Dir(file), line, ok)
	return
}

func (router *Router) RegisterPath(handler http.Handler, path string, methods ...string) (err error) {
	if router.routers == nil {
		router.routers = map[string]pathRouter{}
	}

	p := router.routers[path]
	if len(methods) > 0 {
		for _, method := range methods {
			p.registerMethodHandler(method, handler)
		}
	} else {
		p.registerDefaultHandler(handler)
	}
	router.routers[path] = p
	return
}

func (router *Router) Delegate(handler http.Handler, directory string) {
	if router.delegates == nil {
		router.delegates = map[string]http.Handler{}
	}
	router.delegates[directory] = handler
}

func (router *Router) GET(path string, handler http.HandlerFunc) (err error) {
	err = router.RegisterPath(handler, path, "GET")
	return
}

func (router *Router) POST(path string, handler http.HandlerFunc) (err error) {
	err = router.RegisterPath(handler, path, "POST")
	return
}

func (router *Router) PUT(path string, handler http.HandlerFunc) (err error) {
	err = router.RegisterPath(handler, path, "PUT")
	return
}

func (router *Router) PATCH(path string, handler http.HandlerFunc) (err error) {
	err = router.RegisterPath(handler, path, "PATCH")
	return
}

func (router *Router) DELETE(path string, handler http.HandlerFunc) (err error) {
	err = router.RegisterPath(handler, path, "DELETE")
	return
}

func (router *Router) Default(handler http.HandlerFunc, directory string) (err error) {
	if router.defaultHandler == nil {
		router.defaultHandler = map[string]http.Handler{}
	}
	router.defaultHandler[directory] = handler
	return
}

func (router *Router) Controller(controller Controller) (err error) {
	err = router.RegisterPath(http.HandlerFunc(controller.POST), controller.Path(), "POST")
	if err != nil {
		return
	}
	err = router.RegisterPath(http.HandlerFunc(controller.GET), controller.Path(), "GET")
	if err != nil {
		return
	}
	err = router.RegisterPath(http.HandlerFunc(controller.PATCH), controller.Path(), "PATCH")
	if err != nil {
		return
	}
	err = router.RegisterPath(http.HandlerFunc(controller.PUT), controller.Path(), "PUT")
	if err != nil {
		return
	}
	err = router.RegisterPath(http.HandlerFunc(controller.DELETE), controller.Path(), "DELETE")
	return
}
