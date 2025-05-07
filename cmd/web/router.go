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

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(middleware.RequestID)
	r.Use(slogMiddleware(slog.Default()))
	r.Use(middleware.Recoverer)

	// Brotli compression middleware
	r.Use(CompressionMiddleware)

	// Routes
	r.Get("/", app.indexHandler)
	r.Get("/js", app.jsHandler)
	r.Post("/click/{id}", app.clickHandler)
	r.Post("/mouse/{id}", app.mouseHandler)
	r.Post("/new", app.newGameHandler)

	slog.Info("Serving static files", "source", "filesystem")
	r.Handle("/static/*", hashfs.FileServer(fsys))

	return r
}
