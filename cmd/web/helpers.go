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
    gap: %dpx;}`
)

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

func (s ScoreData) serializeScoreCookie() string {
	return fmt.Sprintf("%d|%d", s.LastScore, s.BestScore)
}

func deserializeScoreCookie(s *ScoreData, data string) {
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

type Action struct {
	Action string
	Index  int
	Score  int
	Board  string
}

func serializeAction(action string, index, score int, board string) string {
	return fmt.Sprintf("%s|%d|%d|%s", action, index, score, board)
}

func deserializeAction(data string) *Action {
	parts := strings.Split(data, "|")
	if len(parts) != 4 {
		return nil
	}

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	score, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	return &Action{
		Action: parts[0],
		Index:  index,
		Score:  score,
		Board:  parts[3],
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 42)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//func perfectSquare(n int) (int, error) {
//	if n < 0 {
//		return 0, errors.New("cannot calculate square root of a negative number")
//	}
//
//	sqrt := int(math.Sqrt(float64(n)))
//
//	// Check if the square of `sqrt` is equal to the original number `n`
//	if sqrt*sqrt == n {
//		return sqrt, nil
//	}
//
//	return 0, errors.New("the given number is not a perfect square")
//}
//
//func restoreGame(pieces string, score int) (*jawbreaker.Game, error) {
//	board := base64ToBytes(pieces)
//	size, err := perfectSquare(len(pieces))
//	if err != nil {
//		return nil, err
//	}
//
//	game, err := jawbreaker.RestoreGame(pieces, size, size, score)
//	if err != nil {
//		return nil, err
//	}
//
//	return game, nil
//}

// registerClient adds a new client channel with a name
func (app *application) registerClient(name string) chan string {
	slog.Info("Registering client", "name", name)
	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()

	clientChan := make(chan string)
	app.clients[name] = clientChan
	return clientChan
}

// unregisterClient removes a client channel
func (app *application) unregisterClient(name string) {
	slog.Info("Unregistering client", "name", name)
	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()

	if clientChan, ok := app.clients[name]; ok {
		close(clientChan)
		delete(app.clients, name)
	}
}
