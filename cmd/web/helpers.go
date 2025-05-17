package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

const (
	piecePrefix = "piece"
	gameSize    = `#game {
    grid-template-columns: repeat(%d, %dpx);
    grid-template-rows: repeat(%d, %dpx);
    gap: 2px;}`
)

func gameSizeCSS(size, block int) string {
	return fmt.Sprintf(gameSize, size, block, size, block)
}

// extractNumberFromPieceId checks if a string starts with "piece" followed by
// a non-negative integer. If so, it returns the integer. Otherwise, it returns -1.
func extractNumberFromPieceId(input string) int {
	if !strings.HasPrefix(input, piecePrefix) {
		return -1 // Prefix doesn't match.
	}

	numberPart := input[len(piecePrefix):]
	if numberPart == "" {
		return -1
	}

	num, err := strconv.Atoi(numberPart)
	if err != nil || num < 0 {
		return -1
	}

	return num
}

func (s ScoreData) serializeCookie() string {
	return fmt.Sprintf("%d|%d", s.LastScore, s.BestScore)
}

func deserializeCookie(s *ScoreData, data string) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		return
	}

	lastScore, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	bestScore, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	s.LastScore = lastScore
	s.BestScore = bestScore
}

func generateSessionID() (string, error) {
	b := make([]byte, 24)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// registerClient adds a new client channel with a name
func (app *Application) registerClient(name string) chan Action {
	slog.Info("Registering client", "name", name)
	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()

	clientChan := make(chan Action)
	app.clients[name] = clientChan
	return clientChan
}

// unregisterClient removes a client channel
func (app *Application) unregisterClient(name string) {
	slog.Info("Unregistering client", "name", name)
	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()
	app.gamesMux.Lock()
	defer app.gamesMux.Unlock()

	if clientChan, ok := app.clients[name]; ok {
		close(clientChan)
		delete(app.clients, name)
	}

	if _, ok := app.games[name]; ok {
		delete(app.games, name)
	}
}
