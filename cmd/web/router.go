package main

import (
	"embed"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS
var fsys = hashfs.NewFS(staticFiles)

func (app *Application) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(slogMiddleware)
	mux.Use(recoverMiddleware)
	compressionMiddleware(mux)

	mux.Get("/", app.indexHandler)
	mux.Get("/one", app.oneHandler)
	mux.Get("/datastar", app.datastarHandler)
	mux.Post("/click/{id}", app.clickHandler)
	mux.Post("/mouse/{id}", app.mouseHandler)
	mux.Post("/mouse/", app.mouseHandler)
	mux.Post("/new", app.newGameHandler)

	mux.Handle("/static/*", hashfs.FileServer(fsys))

	return mux
}
