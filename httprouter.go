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
