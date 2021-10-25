package httprouter

import (
	"net/http"
)

type Router struct {
	routers map[string]pathRouter
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if routers, ok := router.routers[path]; !ok {
		http.NotFound(w, r)
		return
	} else {
		h, err := routers.handlerByMethod(r.Method)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func (router *Router) Register(handler http.Handler, path string, methods ...string) (err error) {
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

func (router *Router) GET(path string, handler http.HandlerFunc) (err error) {
	err = router.Register(handler, path, "GET")
	return
}

func (router *Router) POST(path string, handler http.HandlerFunc) (err error) {
	err = router.Register(handler, path, "POST")
	return
}

func (router *Router) PUT(path string, handler http.HandlerFunc) (err error) {
	err = router.Register(handler, path, "PUT")
	return
}

func (router *Router) PATCH(path string, handler http.HandlerFunc) (err error) {
	err = router.Register(handler, path, "PATCH")
	return
}

func (router *Router) DELETE(path string, handler http.HandlerFunc) (err error) {
	err = router.Register(handler, path, "DELETE")
	return
}

func (router *Router) Controller(controller Controller) (err error) {
	err = router.Register(http.HandlerFunc(controller.POST), controller.Path(), "POST")
	if err != nil {
		return
	}
	err = router.Register(http.HandlerFunc(controller.GET), controller.Path(), "GET")
	if err != nil {
		return
	}
	err = router.Register(http.HandlerFunc(controller.PATCH), controller.Path(), "PATCH")
	if err != nil {
		return
	}
	err = router.Register(http.HandlerFunc(controller.PUT), controller.Path(), "PUT")
	if err != nil {
		return
	}
	err = router.Register(http.HandlerFunc(controller.DELETE), controller.Path(), "DELETE")
	return
}
