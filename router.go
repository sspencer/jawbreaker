package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed index.html
var indexHTML []byte

//go:embed static
var staticFiles embed.FS

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Brotli compression middleware
	r.Use(CompressionMiddleware)

	// Routes
	r.Get("/", indexHandler)
	r.Post("/move/{id}", moveHandler)
	r.Post("/new", newGameHandler)

	// Serve static files
	if os.Getenv("STATIC") == "1" {
		// Serve static files from the filesystem
		log.Println("Serving static files from the filesystem")
		fileServer := http.FileServer(http.Dir("static"))
		r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	} else {
		// Serve static files from the embedded filesystem
		// The embedded filesystem includes the "static" directory itself
		log.Println("Serving static files from embedded filesystem")
		staticFS, err := fs.Sub(staticFiles, "static")
		if err != nil {
			panic(err)
		}
		fileServer := http.FileServer(http.FS(staticFS))
		r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	}

	return r
}
