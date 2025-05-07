package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"net/http"

	"github.com/joho/godotenv"
)

var (
	port        = 8000
	numRows     = 12
	numCols     = 12
	block       = 40
	border      = 6
	mobileBlock = 30
	gap         = 1
	cookieName  = "jawbreaker"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	loadEnv()
	slog.Info("Starting server", "port", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), newRouter())
	if err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func loadEnv() {
	err := godotenv.Load()

	if err == nil {
		slog.Info("Environment configuration loaded", "source", ".env file")
	}
	port = envInt("PORT", port)
	numRows = envInt("JB_ROWS", numRows)
	numCols = envInt("JB_COLS", numCols)
	block = envInt("JB_BLOCK", block)
	border = envInt("JB_BORDER", border)
	mobileBlock = envInt("JB_MOBILE_BLOCK", mobileBlock)
	gap = envInt("JB_GAP", gap)
	cookieName = envString("JB_COOKIE_NAME", cookieName)

	// Log all configuration values as structured data
	slog.Info("Configuration values",
		"port", port,
		"rows", numRows,
		"cols", numCols,
		"block", block,
		"mobile", mobileBlock,
		"border", border,
		"gap", gap,
		"cookie", cookieName)
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
