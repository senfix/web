package web

import (
	"fmt"
	"net/http"
	"strings"

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

func (r *Router) Wildcart(domain string) *mux.Route {
	return r.PathPrefix("/").Subrouter().MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {

		fmt.Printf("%+v\n", domain)
		fmt.Printf("%+v\n", r.Host)

		uri := r.RequestURI
		if strings.HasPrefix(uri, "/socket") {
			return false
		}
		return true
	})
}

func (r *appRouter) Setup() *Router {
	router := &Router{mux.NewRouter()}

	router.PathPrefix(StaticDir).Handler(http.StripPrefix(StaticDir, http.FileServer(http.Dir("./website/"+StaticDir))))

	for _, a := range r.routes {
		a.Register(router)
	}

	return router
}
