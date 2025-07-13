package server

import (
	"io"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"minibb/internal/handlers"
)

func (s *Server) setupRoutes() {
	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.HealthCheck)
		// TODO: Add other API endpoints
	})

	// Static file serving for production
	if !isDevelopment() {
		s.setupStaticFileServing()
	}
}

func (s *Server) setupStaticFileServing() {
	if s.staticFiles == nil {
		// No static files to serve (dev mode)
		return
	}

	// Serve embedded static files
	distFS, err := fs.Sub(*s.staticFiles, "web/dist")
	if err != nil {
		// If no embedded files, serve nothing (dev mode)
		return
	}

	fileServer := http.FileServer(http.FS(distFS))

	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file
		if r.URL.Path != "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		// For root path, serve index.html
		file, err := distFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/html")
		http.ServeContent(w, r, "index.html", time.Time{}, file.(io.ReadSeeker))
	})
}
