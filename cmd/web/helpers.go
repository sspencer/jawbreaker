package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	piecePrefix = "piece"
	gameSize    = `#game {
    grid-template-rows: repeat(%d, %dpx);
    grid-template-columns: repeat(%d, %dpx);
    gap: 2px;}`
)

func (app *Application) gameSizeCSS() string {
	return fmt.Sprintf(gameSize, app.cfg.rows, app.cfg.block, app.cfg.cols, app.cfg.block)
}

func (app *Application) getCookieName() string {
	return fmt.Sprintf("jawbreaker_score_%dx%d", app.cfg.rows, app.cfg.cols)
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
