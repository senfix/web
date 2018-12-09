package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	STATIC_DIR = "/static/"
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

	router.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("./website/"+STATIC_DIR))))

	for _, a := range r.routes {
		a.Register(router)
	}

	return router
}
