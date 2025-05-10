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
	defPort   = 8000
	defRows   = 12
	defCols   = 12
	defBlock  = 40
	defMBlock = 30
	defBorder = 6
	defGap    = 1
	defCookie = "jawbreaker"
)

type config struct {
	port   int
	rows   int
	cols   int
	block  int
	mblock int
	border int
	gap    int
	cookie string
}
type application struct {
	envFile string
	cfg     config
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var watchFiles bool
	flag.BoolVar(&watchFiles, "w", false, "Watch .env file for changes")
	flag.Parse()

	fn, _ := filepath.Abs(".env")
	app := application{
		envFile: fn,
		cfg:     config{},
	}

	app.loadConfig()

	port := app.cfg.port
	slog.Info("Starting server", "port", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), app.routes())
	if err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
