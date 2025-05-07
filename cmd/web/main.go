package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"net/http"

	"github.com/joho/godotenv"
)

const (
	defPort        = 8000
	defNumRows     = 12
	defNumCols     = 12
	defBlock       = 40
	defBorder      = 6
	defMobileBlock = 30
	defGap         = 1
	defCookieName  = "jawbreaker"
)

type config struct {
	port        int
	numRows     int
	numCols     int
	block       int
	mobileBlock int
	border      int
	gap         int
	cookieName  string
}
type application struct {
	cfg config
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config{}
	loadEnv(&cfg)

	app := application{
		cfg: cfg,
	}

	slog.Info("Starting server", "port", cfg.port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), app.routes())
	if err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func loadEnv(cfg *config) {
	err := godotenv.Load()

	if err == nil {
		slog.Info("Environment configuration loaded", "source", ".env file")
	}

	cfg.port = envInt("PORT", defPort)
	cfg.numRows = envInt("JB_ROWS", defNumRows)
	cfg.numCols = envInt("JB_COLS", defNumCols)
	cfg.block = envInt("JB_BLOCK", defBlock)
	cfg.border = envInt("JB_BORDER", defBorder)
	cfg.mobileBlock = envInt("JB_MOBILE_BLOCK", defMobileBlock)
	cfg.gap = envInt("JB_GAP", defGap)
	cfg.cookieName = envString("JB_COOKIE_NAME", defCookieName)

	// Log all configuration values as structured data
	slog.Info("Configuration values",
		"port", cfg.port,
		"rows", cfg.numRows,
		"cols", cfg.numCols,
		"block", cfg.block,
		"mobile", cfg.mobileBlock,
		"border", cfg.border,
		"gap", cfg.gap,
		"cookie", cfg.cookieName)
}

func envString(key, defValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defValue
	}
	return v
}

func envInt(key string, defValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defValue
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return defValue
	}

	return n
}
