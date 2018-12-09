package web

import "net/http"

type RouterMiddleware interface {
	Set(next http.Handler) http.Handler
}
