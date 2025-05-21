package main

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/benbjohnson/hashfs"
)

var (
	//go:embed tmpl/*.gohtml
	templateFS embed.FS
	tmpl       *template.Template

	//go:embed static
	staticFiles embed.FS
	fsys        = hashfs.NewFS(staticFiles)

	//go:embed static/jawbreaker/*.js
	jawbreakerFS embed.FS
	jsys         = hashfs.NewFS(NewConcatFS(jawbreakerFS, "js/jawbreaker.js"))
)

func init() {
	// Parse all embedded templates during startup
	tmpl = template.Must(template.ParseFS(templateFS, "tmpl/*.gohtml"))
}

// bundleLiveJavascript is used in dev mode and reloads every
// JavaScript file per page load.  Will not work when deployed
// unless the static dir is copied over in the same hierarchy
// relative to the executable.
func bundleLiveJavascript() ([]byte, error) {
	dir := "cmd/web/static/jawbreaker"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var code []byte
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".js") {
			path := filepath.Join(dir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			code = append(code, data...)
			code = append(code, '\n')
		}
	}

	return code, nil
}
