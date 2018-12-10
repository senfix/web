package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/senfix/logger"
)

var LoggerPrefix = "WEB"

type Server interface {
	Start()
	WaitShutdown(close chan bool, closed chan bool)
	Stop()
}

type server struct {
	http.Server
	router      AppRouter
	logger      logger.Log
	shutdownReq chan bool
}

func NewServer(app Config, log logger.Log, router AppRouter) Server {
	s := &server{
		Server: http.Server{
			Addr:         app.Listen,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
		logger:      log.Enable(LoggerPrefix),
		router:      router,
	}

	return s
}

func (s *server) Stop() {
	s.logger.Println("Stopping web")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		s.logger.Panic("Shutdown request error: %v", err, nil)
	}
}

func (s *server) WaitShutdown(close chan bool, closed chan bool) {
	<-close
	s.Stop()
	closed <- true
}

func (s *server) Start() {
	s.logger.Printf("Listening on %v\n", s.Server.Addr)

	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	handler := s.router.Setup()
	handler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		if err2 == nil {
			for _, m := range met {
				s.logger.Printf("%v %v\n", m, tpl)
			}
		}

		return nil
	})
	s.Handler = handlers.CORS(originsOk, headersOk, methodsOk)(handler)
	s.ListenAndServe()
}
