package httprouter

import (
	"errors"
	"net/http"
	"path"
	"runtime"
	"strings"
)

type Router struct {
	routers   map[string]pathRouter
	delegates map[string]http.Handler
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Path

	var (
		hasOtherMethods bool
		handler         http.Handler
		err             error
	)

	if h, ok := router.routers[route]; ok {
		handler, hasOtherMethods, err = h.handlerByMethod(r.Method)
		if err == nil {
			handler.ServeHTTP(w, r)
			return
		}
	}

	if router.delegates != nil {
		for key, handler := range router.delegates {
			if strings.HasPrefix(route, key) {
				handler.ServeHTTP(w, r)
				return
			}
		}
	}

	if hasOtherMethods {
		http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)

	return
}

func (router *Router) register(handler http.Handler, methods ...string) (err error) {
	var directory string
	directory, err = router.getPath()
	if err != nil {
		return
	}
	router.registerPath(handler, directory, methods...)
	return
}

func (router *Router) Register(handler http.Handler) (err error) {
	err = router.register(handler)
	return
}

func (router *Router) RegisterDELETE(handler http.Handler) (err error) {
	err = router.register(handler, "DELETE")
	return
}

func (router *Router) RegisterGET(handler http.Handler) (err error) {
	err = router.register(handler, "GET")
	return
}

func (router *Router) RegisterPATCH(handler http.Handler) (err error) {
	err = router.register(handler, "PATCH")
	return
}

func (router *Router) RegisterPOST(handler http.Handler) (err error) {
	err = router.register(handler, "POST")
	return
}

func (router *Router) RegisterPUT(handler http.Handler) (err error) {
	err = router.register(handler, "PUT")
	return
}

func (router *Router) RegisterMethods(handler http.Handler, methods ...string) (err error) {
	err = router.register(handler, methods...)
	return
}

func (router *Router) registerPath(handler http.Handler, path string, methods ...string) (err error) {
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

func (router *Router) RegisterPath(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path)
	return
}

func (router *Router) RegisterPathDELETE(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path, "DELETE")
	return
}

func (router *Router) RegisterPathGET(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path, "GET")
	return
}

func (router *Router) RegisterPathPATCH(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path, "PATCH")
	return
}

func (router *Router) RegisterPathPOST(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path, "POST")
	return
}

func (router *Router) RegisterPathPUT(handler http.Handler, path string) (err error) {
	router.registerPath(handler, path, "PUT")
	return
}

func (router *Router) RegisterPathMethods(handler http.Handler, path string, methods ...string) (err error) {
	router.registerPath(handler, path, methods...)
	return
}

func (router *Router) delegate(handler http.Handler, directory string) (err error) {
	if router.delegates == nil {
		router.delegates = map[string]http.Handler{}
	}
	router.delegates[directory] = handler
	return
}

func (router *Router) getPath() (directory string, err error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		err = errors.New("cant_detect_package")
		return
	}
	index := strings.Index(file, "/api/")
	if index == -1 {
		err = errors.New("api_path_not_found")
		return
	}

	directory = path.Dir(file[index+len("/api"):])
	return
}

func (router *Router) RegisterDelegate(handler http.Handler) (err error) {
	var directory string
	directory, err = router.getPath()
	if err != nil {
		return
	}
	err = router.delegate(handler, directory)
	return
}

func (router *Router) RegisterDelegatePath(handler http.Handler, directory string) (err error) {
	err = router.delegate(handler, directory)
	return
}

func (router *Router) RegisterController(controller Controller) (err error) {
	var directory string
	directory, err = router.getPath()
	if err != nil {
		return
	}

	err = router.RegisterPathMethods(http.HandlerFunc(controller.POST), directory, "POST")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.GET), directory, "GET")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.PATCH), directory, "PATCH")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.PUT), directory, "PUT")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.DELETE), directory, "DELETE")
	return
}

func (router *Router) RegisterControllerPath(controller Controller, directory string) (err error) {
	err = router.RegisterPathMethods(http.HandlerFunc(controller.POST), directory, "POST")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.GET), directory, "GET")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.PATCH), directory, "PATCH")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.PUT), directory, "PUT")
	if err != nil {
		return
	}
	err = router.RegisterPathMethods(http.HandlerFunc(controller.DELETE), directory, "DELETE")
	return
}
