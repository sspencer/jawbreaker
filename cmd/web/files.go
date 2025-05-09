package main

import (
	"embed"
	"html/template"
)

//go:embed tmpl/*.html
var templateFS embed.FS

var tmpl *template.Template

func init() {
	// Parse all embedded templates once at startup
	tmpl = template.Must(template.ParseFS(templateFS,
		"tmpl/datastar.html",
		"tmpl/gameover.html",
		"tmpl/help.html",
		"tmpl/index.html",
		"tmpl/sidebar.html",
	))
}
