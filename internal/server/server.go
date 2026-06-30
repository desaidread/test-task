package server

import (
	"log/slog"
	"net/http"
	"testTask/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "testTask/docs"
)

type Server struct {
	port    string
	handler *handlers.Handler
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/subscriptions", s.handler.Create)
	r.Get("/subscriptions", s.handler.List)
	r.Get("/subscriptions/total", s.handler.TotalCost)
	r.Get("/subscriptions/{id}", s.handler.GetByID)
	r.Put("/subscriptions/{id}", s.handler.Update)

	r.Delete("/subscriptions/{id}", s.handler.Delete)

	slog.Info("starting server", "port", s.port)
	return http.ListenAndServe(":"+s.port, r)
}

func New(port string, handler *handlers.Handler) *Server {
	return &Server{port: port, handler: handler}
}
