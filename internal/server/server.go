package server

import (
	"context"
	"net/http"
	"time"

	"go-final-progect/internal/handlers"
	"go-final-progect/internal/storage/sqlite"
)

type AppServer struct {
	Server *http.Server
}

func NewServer(port string, webDir string, store *sqlite.Storage) *AppServer {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	handlers.RegisterRoutes(mux, store)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &AppServer{
		Server: srv,
	}
}

func (a *AppServer) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}
