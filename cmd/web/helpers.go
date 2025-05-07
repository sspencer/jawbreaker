package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/sspencer/jawbreaker"
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

// gameToHTML generates the HTML fragment for the game board with the click handler.
func gameToHTML(g *jawbreaker.Game, connections []int) string {
	var sb strings.Builder

	set := make(map[int]bool)
	if len(connections) > 1 {
		for _, c := range connections {
			set[c] = true
		}
	}

	for i, p := range g.Board() {
		connected := set[i]

		// id="piece123"
		sb.WriteString("<div id=\"")
		sb.WriteString(piecePrefix)
		sb.WriteString(strconv.Itoa(i))

		// class="piece red connected"
		sb.WriteString("\" class=\"piece ")
		sb.WriteString(p.Color())
		if connected {
			sb.WriteString(" connected")
		}
		sb.WriteString("\"></div>")
	}

	return fmt.Sprintf("<div id=\"game\">%s</div>", sb.String())
}

func (s ScoreData) serialize() string {
	return fmt.Sprintf("%d|%d", s.LastScore, s.BestScore)
}

func deserializeScoreData(s *ScoreData, data string) {
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

func getBlockSize(r *http.Request) int {
	if isMobile(r) {
		return mobileBlock
	}

	return block
}

func getGameSize(rows, cols, blockSize, gapSize int) string {
	return fmt.Sprintf(gameSize, cols, blockSize, rows, blockSize, gapSize)
}

// isMobile checks if the request is coming from a mobile device
// by examining the User-Agent header
func isMobile(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")

	// Common patterns for mobile devices
	mobilePatterns := []string{
		"Android",
		"iPhone",
		"iPad",
		"iPod",
		"BlackBerry",
		"Windows Phone",
		"Mobile",
		"webOS",
		"Opera Mini",
	}

	for _, pattern := range mobilePatterns {
		matched, _ := regexp.MatchString(pattern, userAgent)
		if matched {
			return true
		}
	}

	return false
}
