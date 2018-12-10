package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	StaticDir = "/static/"
)

type RouteMatch struct {
	*mux.RouteMatch
}

type Router struct {
	*mux.Router
}

type RouterPath interface {
	Register(root *Router)
}

type AppRouter interface {
	Setup() *Router
}

type appRouter struct {
	routes []RouterPath
}

func NewRouter(routes []RouterPath) AppRouter {
	r := &appRouter{
		routes: routes,
	}
	return r
}

func (r *appRouter) Setup() *Router {
	router := &Router{mux.NewRouter()}

	router.PathPrefix(StaticDir).Handler(http.StripPrefix(StaticDir, http.FileServer(http.Dir("./website/"+StaticDir))))

	for _, a := range r.routes {
		a.Register(router)
	}

	return router
}
