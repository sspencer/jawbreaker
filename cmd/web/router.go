package main

import (
	"embed"
	"log/slog"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS
var fsys = hashfs.NewFS(staticFiles)

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(middleware.RequestID)
	r.Use(slogMiddleware(slog.Default()))
	r.Use(middleware.Recoverer)

	// Brotli compression middleware
	r.Use(CompressionMiddleware)

	// Routes
	r.Get("/", indexHandler)
	r.Get("/js", jsHandler)
	r.Post("/click/{id}", clickHandler)
	r.Post("/mouse/{id}", mouseHandler)
	r.Post("/new", newGameHandler)

	slog.Info("Serving static files", "source", "filesystem")
	r.Handle("/static/*", hashfs.FileServer(fsys))

	return r
}
