package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type pageData struct {
	GameCode   template.JS
	GameSrc    string
	StyleCSS   string
	Size       int
	Block      int
	MSize      int
	MBlock     int
	CookieName string
}

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	data := pageData{
		GameSrc:    fsys.HashName("static/jawbreaker.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		CookieName: cookieName,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) oneHandler(w http.ResponseWriter, r *http.Request) {

	gameBytes, err := os.ReadFile("cmd/web/static/jawbreaker.js")
	if err != nil {
		slog.Error("jawbreaker.js not found", "error", err)
		gameBytes = gameCode
	}

	cfg := app.cfg
	data := pageData{
		GameCode:   template.JS(gameBytes),
		StyleCSS:   fsys.HashName("static/style.css"),
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		CookieName: cookieName,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "one", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
