package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	piecePrefix = "piece"
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
func gameToHTML(g *Game, connections []int) string {
	var sb strings.Builder

	set := make(map[int]bool)
	if len(connections) > 1 {
		for _, c := range connections {
			set[c] = true
		}
	}

	board := g.String()

	for i, p := range board {
		connected := set[i]

		// id="piece123"
		sb.WriteString("<div id=\"")
		sb.WriteString(piecePrefix)
		sb.WriteString(strconv.Itoa(i))

		// class="piece red connected"
		sb.WriteString("\" class=\"piece ")
		sb.WriteString(PieceFromChar(p).String())
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
