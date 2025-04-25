package main

import (
	"errors"
	"math/rand/v2"
)

type Piece byte

// Game represents the state of a Jawbreaker game.
// It contains the game board, dimensions, and current score.
type Game struct {
	pieces []Piece // pieces represents the game board as a flat array of color bytes
	rows   int     // number of rows in the game board
	cols   int     // number of columns in the game board
	score  int     // current score of the game
}

// Piece constants used to represent different colored pieces in the game.
// Each color is represented by a single byte character.
const (
	White  Piece = 'w' // "empty" or removed piece
	Purple Piece = 'p' // purple piece
	Blue   Piece = 'b' // blue piece
	Green  Piece = 'g' // green piece
	Red    Piece = 'r' // red piece
	Yellow Piece = 'y' // yellow piece
)

var (
	// colorPieces is a slice of all available colored pieces used for random piece generation.
	colorPieces = []Piece{Purple, Blue, Green, Red, Yellow}

	// colorMap maps color byte values to their string representation for HTML class names.
	colorMap = map[Piece]string{
		Purple: "purple",
		Blue:   "blue",
		Green:  "green",
		Red:    "red",
		Yellow: "yellow",
		White:  "white",
	}

	// ErrGameSize is returned when attempting to restore a game with a board size
	// that doesn't match the expected dimensions.
	ErrGameSize = errors.New("game board does not match expected size")
)

// String converts a game piece to its color name
func (c Piece) String() string {
	if name, ok := colorMap[c]; ok {
		return name
	}
	return "white"
}

func PieceFromChar(c rune) Piece {
	switch c {
	case 'p':
		return Purple
	case 'b':
		return Blue
	case 'g':
		return Green
	case 'r':
		return Red
	case 'y':
		return Yellow
	default:
		return White
	}
}

// NewGame initializes a new Game instance with the given rows and columns, populating the board with random pieces.
func NewGame(rows, cols int) *Game {
	pieces := make([]Piece, rows*cols)

	for i := 0; i < len(pieces); i++ {
		pieces[i] = colorPieces[rand.IntN(len(colorPieces))]
	}

	return &Game{pieces: pieces, rows: rows, cols: cols}
}

// RestoreGame restores a game state based on the given board string, dimensions, and score.
// It returns a pointer to a Game instance or an error if the board size does not match the expected size.
func RestoreGame(board string, rows, cols, score int) (*Game, error) {
	if len(board) != rows*cols {
		return nil, ErrGameSize
	}

	pieces := make([]Piece, rows*cols)
	for i, c := range board {
		pieces[i] = PieceFromChar(c)
	}

	return &Game{pieces: pieces, rows: rows, cols: cols, score: score}, nil
}

// String returns the current state of the board as a string,
// where each character in the string represents a color
func (g *Game) String() string {
	return string(g.pieces)
}

// Score returns the current score of the game.
func (g *Game) Score() int {
	return g.score
}

// Move removes connected pieces of the same color from the board and updates the score.
// It takes the index of the clicked piece and removes all connected pieces of the same color.
// If fewer than 2 connected pieces are found, no pieces are removed and the score remains unchanged.
// After removing pieces, gravity is applied to make pieces fall down, and empty columns are shifted right.
// If the game is over after the move, a bonus score is added based on the number of remaining pieces.
func (g *Game) Move(index int) {
	n := g.countConnectedPieces(index)
	if n < 2 {
		return
	}

	n = g.floodFill(index)
	g.applyGravityAndShiftRight()

	g.score += calculateMoveScore(n)

	if g.IsGameOver() {
		remainingPieces := 0
		for _, p := range g.pieces {
			if p != White {
				remainingPieces++
			}
		}
		g.score += calculateRemainingPiecesScore(remainingPieces)
	}
}

// IsGameOver checks if the game is over by determining if there are any valid moves left.
// A valid move requires at least two connected pieces of the same color.
// Returns true if the game is over (no valid moves), false otherwise.
func (g *Game) IsGameOver() bool {
	for i := 0; i < len(g.pieces); i++ {
		if g.pieces[i] == White {
			continue
		}
		if g.countConnectedPieces(i) > 1 {
			return false
		}
	}

	return true
}

// calculateMoveScore computes the score for removing a group of connected pieces.
// The score is calculated as n * (n-1), where n is the number of pieces removed.
// If fewer than 2 pieces are removed, the score is 0.
// This scoring system rewards removing larger groups of pieces with a quadratic score increase.
func calculateMoveScore(piecesRemoved int) int {
	if piecesRemoved < 2 {
		return 0
	}
	return piecesRemoved * (piecesRemoved - 1)
}

// calculateRemainingPiecesScore computes a bonus score at the end of the game
// based on the number of pieces remaining on the board.
// If the number of remaining pieces is less than or equal to the threshold (20),
// a bonus score is awarded. The bonus is calculated as (threshold - remainingPieces)².
// This rewards players who clear most of the board with an exponentially increasing bonus.
// If more than the threshold number of pieces remain, no bonus is awarded.
func calculateRemainingPiecesScore(remainingPieces int) int {
	const threshold = 20
	if remainingPieces <= threshold {
		p := threshold - remainingPieces
		return p * p
	}

	return 0
}

// countConnectedPieces counts how many pieces of the same color are connected to the piece at the given index.
// This is a non-destructive version of floodFill that doesn't modify the game state.
func (g *Game) countConnectedPieces(index int) int {
	if index < 0 || index >= len(g.pieces) {
		return 0
	}

	target := g.pieces[index]
	if target == White {
		return 0
	}

	count := 0
	stack := []int{index}
	// Pre-allocate the visited map with a reasonable capacity
	visited := make([]bool, g.rows*g.cols)

	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i < 0 || i >= len(g.pieces) || g.pieces[i] != target || visited[i] {
			continue
		}
		visited[i] = true
		count++

		row, col := i/g.cols, i%g.cols
		// Check all four adjacent positions (up, down, left, right)
		if row > 0 {
			stack = append(stack, i-g.cols) // Up
		}
		if row < g.rows-1 {
			stack = append(stack, i+g.cols) // Down
		}
		if col > 0 {
			stack = append(stack, i-1) // Left
		}
		if col < g.cols-1 {
			stack = append(stack, i+1) // Right
		}
	}
	return count
}

func (g *Game) floodFill(index int) int {
	target := g.pieces[index]
	if target == White {
		return 0
	}

	count := 0
	stack := []int{index}
	visited := make(map[int]bool)

	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i < 0 || i >= len(g.pieces) || g.pieces[i] != target || visited[i] {
			continue
		}
		g.pieces[i] = White
		visited[i] = true
		count++

		row, col := i/g.cols, i%g.cols
		if row > 0 {
			stack = append(stack, i-g.cols)
		}
		if row < g.rows-1 {
			stack = append(stack, i+g.cols)
		}
		if col > 0 {
			stack = append(stack, i-1)
		}
		if col < g.cols-1 {
			stack = append(stack, i+1)
		}
	}
	return count
}

// ApplyGravityAndShiftRight modifies the board in-place to apply vertical gravity
// and then right-align non-empty columns.
func (g *Game) applyGravityAndShiftRight() {
	// Gravity Phase: shift non-'w' characters down in each column
	for col := 0; col < g.cols; col++ {
		writeRow := g.rows - 1
		for row := g.rows - 1; row >= 0; row-- {
			index := row*g.cols + col
			if g.pieces[index] != White {
				g.pieces[writeRow*g.cols+col] = g.pieces[index]
				if writeRow != row {
					g.pieces[index] = White
				}
				writeRow--
			}
		}
	}

	// Right Shift Phase: move non-empty columns (columns that have any non-'w') to the right
	writeCol := g.cols - 1
	for col := g.cols - 1; col >= 0; col-- {
		isEmpty := true
		for row := 0; row < g.rows; row++ {
			if g.pieces[row*g.cols+col] != White {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			if writeCol != col {
				// Copy column to new position
				for row := 0; row < g.rows; row++ {
					g.pieces[row*g.cols+writeCol] = g.pieces[row*g.cols+col]
					g.pieces[row*g.cols+col] = White
				}
			}
			writeCol--
		}
	}
}
