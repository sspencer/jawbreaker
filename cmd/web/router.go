package main

import (
	"embed"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS
var fsys = hashfs.NewFS(staticFiles)

func (app *application) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Timeout(30 * time.Second))
	mux.Use(slogMiddleware)
	mux.Use(recoverMiddleware)
	compressionMiddleware(mux)

	mux.Get("/", app.indexHandler)
	mux.Get("/star", app.datastarHandler)
	mux.Post("/click/{id}", app.clickHandler)
	mux.Post("/mouse/{id}", app.mouseHandler)
	mux.Post("/new", app.newGameHandler)

	mux.Handle("/static/*", hashfs.FileServer(fsys))

	return mux
}
