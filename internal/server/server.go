package server

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router      *chi.Mux
	db          *sql.DB
	port        string
	staticFiles *embed.FS
}

func New(db *sql.DB, staticFiles *embed.FS) *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := &Server{
		router:      chi.NewRouter(),
		db:          db,
		port:        port,
		staticFiles: staticFiles,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Server starting on port %s\n", s.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
