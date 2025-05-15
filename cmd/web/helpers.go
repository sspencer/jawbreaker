package main

import (
	"errors"
	"fmt"
	"math"
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
		sb.WriteString(pieceColor(p))
		if connected {
			sb.WriteString(" connected")
		}
		sb.WriteString("\">")
		sb.WriteString(pieceIcon(p))
		sb.WriteString("</div>")
	}

	return fmt.Sprintf("<div id=\"game\">%s</div>", sb.String())
}

func pieceColor(c jawbreaker.Piece) string {
	switch c.Color() {
	case jawbreaker.Purple:
		return "purple"
	case jawbreaker.Blue:
		return "blue"
	case jawbreaker.Green:
		return "green"
	case jawbreaker.Red:
		return "red"
	case jawbreaker.Yellow:
		return "yellow"
	default:
		if c.Power() == 0 {
			return "white"
		}
		return "gray"
	}
}

func pieceIcon(c jawbreaker.Piece) string {
	// https://www.w3schools.com/charsets/ref_utf_symbols.asp
	switch c.Power() {
	case jawbreaker.PowerX:
		return "&#10005;"
	case jawbreaker.PowerPlus:
		return "&#43;" //"&#9532;"
	case jawbreaker.PowerCircle:
		return "&#1054;"
	case jawbreaker.PowerRect:
		return "&#127020;"
	case jawbreaker.PowerFill:
		return "&#9734;"
	case jawbreaker.PowerRotateLeft:
		return "&#8617;"
	case jawbreaker.PowerRotateRight:
		return "&#8618;"
	default:
		return ""
	}
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

func perfectSquare(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("cannot calculate square root of a negative number")
	}

	sqrt := int(math.Sqrt(float64(n)))

	// Check if the square of `sqrt` is equal to the original number `n`
	if sqrt*sqrt == n {
		return sqrt, nil
	}

	return 0, errors.New("the given number is not a perfect square")
}

func restoreGame(pieces string, score int) (*jawbreaker.Game, error) {
	size, err := perfectSquare(len(pieces))
	if err != nil {
		return nil, err
	}

	game, err := jawbreaker.RestoreGame(pieces, size, size, score)
	if err != nil {
		return nil, err
	}

	return game, nil
}
