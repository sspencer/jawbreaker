package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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

func (app *Application) loadConfig() {
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
	app.cfg.mrows = envInt("MROWS", defMRows)
	app.cfg.mcols = envInt("MCOLS", defMCols)
	app.cfg.mblock = envInt("MBLOCK", defMBlock)

	// Log all configuration values as structured data
	slog.Info("Configuration values",
		"port", app.cfg.port,
		"block", app.cfg.block,
		"rows", app.cfg.rows,
		"cols", app.cfg.cols,
		"mblock", app.cfg.mblock,
		"mrows", app.cfg.mrows,
		"mcols", app.cfg.mcols,
	)

}
