package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const (
	defPort    = 8000
	defRows    = 11
	defCols    = 13
	defBlock   = 40
	defMRows   = 7
	defMCols   = 9
	defMBlock  = 30
	cookieName = "score"
)

type Config struct {
	port    int
	rows    int
	cols    int
	block   int
	mblock  int
	mrows   int
	mcols   int
	animate bool
}

type Application struct {
	envFile string
	cfg     Config
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var watchFiles bool
	flag.BoolVar(&watchFiles, "w", false, "Watch .env file for changes")
	flag.Parse()

	fn, _ := filepath.Abs(".env")
	app := Application{
		envFile: fn,
		cfg:     Config{},
	}

	app.loadConfig()
	app.cfg.animate = true

	port := app.cfg.port
	slog.Info("Starting server", "port", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), app.routes())
	if err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
