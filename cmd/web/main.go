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
	port            = 8000
	numRows         = 12
	numCols         = 12
	blockSize       = 40
	mobileBlockSize = 30
	gapSize         = 2
	cookieName      = "scores"
	titleTime       = 2_000 // 2 seconds
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
	blockSize = envInt("JB_BLOCK_SIZE", blockSize)
	mobileBlockSize = envInt("JB_MOBILE_BLOCK_SIZE", mobileBlockSize)
	gapSize = envInt("JB_GAP_SIZE", gapSize)
	cookieName = envString("JB_COOKIE_NAME", cookieName)
	titleTime = envInt("JB_TITLE_TIME", titleTime)

	// Log all configuration values as structured data
	slog.Info("Configuration values",
		"port", port,
		"rows", numRows,
		"cols", numCols,
		"block", blockSize,
		"mobile", mobileBlockSize,
		"gap", gapSize,
		"cookie", cookieName,
		"timer", titleTime)
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
