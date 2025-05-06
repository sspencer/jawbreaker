package main

import _ "embed"

var (
	//go:embed tmpl/index.gohtml
	indexHTML []byte

	//go:embed tmpl/js.gohtml
	jsHTML []byte

	//go:embed tmpl/canvas.gohtml
	canvasHTML []byte

	//go:embed tmpl/gameover.html
	gameoverHTML []byte

	//go:embed tmpl/help.html
	helpHTML []byte

	//go:embed tmpl/sidebar.html
	sidebarHTML []byte
)
