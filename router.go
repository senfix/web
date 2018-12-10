package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	StaticDir = "/static/"
)

type RouterPath interface {
	Register(root *mux.Router) (p *mux.Router)
}

type AppRouter interface {
	Setup() *mux.Router
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

func (r *appRouter) Setup() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix(StaticDir).Handler(http.StripPrefix(StaticDir, http.FileServer(http.Dir("./website/"+StaticDir))))

	for _, a := range r.routes {
		a.Register(router)
	}

	return router
}
