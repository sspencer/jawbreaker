package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"net/http"

	"github.com/fsnotify/fsnotify"
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
	cfgLock sync.RWMutex
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

	if watchFiles {
		app.watcher()
	}

	port := app.cfg.port
	slog.Info("Starting server", "port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

// watcher monitors the .env file for changes and updates the application config
func (app *application) watcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to create file watcher", "error", err)
		return
	}

	defer watcher.Close()

	slog.Info("Started watching .env file for changes", "path", app.envFile)

	// Create a debounce timer to prevent multiple reloads for a single change
	var debounceTimer *time.Timer
	var debounceTimerMutex sync.Mutex

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					debounceTimerMutex.Lock()
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
						slog.Info("Detected change in .env file, reloading configuration")
						app.reloadConfig()
					})
					debounceTimerMutex.Unlock()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("Error watching .env file", "error", err)
			}
		}
	}()

	// Watch the directory, not the file itself as editors
	// save files in different ways.
	err = watcher.Add(filepath.Dir(app.envFile))
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
