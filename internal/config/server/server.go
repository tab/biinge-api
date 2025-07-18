package server

import (
	"context"
	"net/http"
	"time"

	"biinge-api/internal/config"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, appRouter http.Handler) Server {
	return &server{
		httpServer: &http.Server{
			Addr:         cfg.AppAddr,
			Handler:      appRouter,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
