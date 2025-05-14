package main

import (
	"embed"
	"html/template"
)

//go:embed tmpl/*.html
var templateFS embed.FS

//go:embed static/jawbreaker.js
var gameCode []byte

var tmpl *template.Template

func init() {
	// Parse all embedded templates during startup
	tmpl = template.Must(template.ParseFS(templateFS, "tmpl/*.html"))
}
