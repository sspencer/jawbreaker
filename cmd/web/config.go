package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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

func (app *application) loadConfig() {
	err := godotenv.Load(app.envFile)

	if err == nil {
		slog.Info("Environment configuration loaded", "file", app.envFile)
	} else {
		slog.Info("No .env file found, using default configuration")
	}

	app.cfg.port = envInt("PORT", defPort)
	app.cfg.rows = envInt("ROWS", defRows)
	app.cfg.cols = envInt("COLS", defCols)
	app.cfg.block = envInt("BLOCK", defBlock)
	app.cfg.mblock = envInt("M_BLOCK", defMBlock)
	app.cfg.border = envInt("BORDER", defBorder)
	app.cfg.gap = envInt("GAP", defGap)
	app.cfg.cookie = envString("COOKIE", defCookie)

	// Log all configuration values as structured data
	slog.Info("Configuration values",
		"port", app.cfg.port,
		"rows", app.cfg.rows,
		"cols", app.cfg.cols,
		"block", app.cfg.block,
		"mblock", app.cfg.mblock,
		"border", app.cfg.border,
		"gap", app.cfg.gap,
		"cookie", app.cfg.cookie)
}
