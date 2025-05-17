package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sspencer/jawbreaker"
)

const (
	defPort    = 8000
	defSize    = 16
	defBlock   = 40
	defMSize   = 12
	defMBlock  = 30
	cookieName = "score"
)

type Config struct {
	port    int
	size    int
	block   int
	mblock  int
	msize   int
	animate bool
}
type Action struct {
	Action  string
	Index   int
	Session string
	Board   string
}

func (a Action) String() string {
	return fmt.Sprintf("%s[%d] session=%s board-len=%d",
		a.Action, a.Index, a.Session, len(a.Board))
}

type Application struct {
	envFile    string
	cfg        Config
	clients    map[string]chan Action
	clientsMux sync.Mutex
	games      map[string]*jawbreaker.Game
	gamesMux   sync.Mutex
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
		clients: make(map[string]chan Action),
		games:   make(map[string]*jawbreaker.Game),
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
