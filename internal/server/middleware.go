package server

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) setupMiddleware() {
	// Basic logging
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// CORS for development mode
	if isDevelopment() {
		s.router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// Request timeout
	s.router.Use(middleware.Timeout(30 * time.Second))
}

func isDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
