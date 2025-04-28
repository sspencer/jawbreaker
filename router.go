package main

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed index.html
var indexHTML []byte

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
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return r
}
