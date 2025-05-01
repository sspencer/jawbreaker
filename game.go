package jawbreaker

import (
	"errors"
	"math/rand/v2"
)

type Piece byte

type Board []Piece

func (b Board) String() string {
	return string(b)
}

// Game represents the state of a Jawbreaker game.
// It contains the game board, dimensions, and current score.
type Game struct {
	board Board // board represents the game board as a flat array of color bytes
	rows  int   // number of rows in the game board
	cols  int   // number of columns in the game board
	score int   // current score of the game
}

type Status struct {
	Board    Board
	Score    int
	GameOver bool
}

// Piece constants used to represent different colored board in the game.
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
	// colorPieces is a slice of all available colored board used for random piece generation.
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

// Color returns each piece as a string representing its color.
// It returns "white" for any other piece.
func (c Piece) Color() string {
	if name, ok := colorMap[c]; ok {
		return name
	}
	return "white"
}

func fromChar(c rune) Piece {
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

// NewGame initializes a new Game instance with the given rows and columns, populating the board with random board.
func NewGame(rows, cols int) *Game {
	pieces := make(Board, rows*cols)

	for i := 0; i < len(pieces); i++ {
		pieces[i] = colorPieces[rand.IntN(len(colorPieces))]
	}

	return &Game{board: pieces, rows: rows, cols: cols}
}

// RestoreGame restores a game state based on the given board string, dimensions, and score.
// It returns a pointer to a Game instance or an error if the board size does not match the expected size.
func RestoreGame(board string, rows, cols, score int) (*Game, error) {
	if len(board) != rows*cols {
		return nil, ErrGameSize
	}

	pieces := make(Board, rows*cols)
	for i, c := range board {
		pieces[i] = fromChar(c)
	}

	return &Game{board: pieces, rows: rows, cols: cols, score: score}, nil
}

// Board returns the current state of the board as a string,
// where each Piece in the string represents a color.  Pieces
// are represented as a single character, where 'w' represents
// an empty space, 'p' represents purple, 'b' represents blue,
// 'g' represents green, 'r' represents red, and 'y' represents yellow.
// The slice is returned in row-major order, with the first row
// at the beginning of the string.
func (g *Game) Board() Board {
	return g.board
}

// Score returns the current score of the game.
func (g *Game) Score() int {
	return g.score
}

// Move removes connected pieces of the same color from the board and updates the score.
// It takes the index of the clicked piece and removes all connected pieces of the same color.
// If less than 2 connected pieces are found, no pieces are removed and the score remains unchanged.
// After removing board, gravity is applied to make pieces fall down, and empty columns are shifted right.
// If the game is over after the move, a bonus score is added based on the number of remaining pieces.
func (g *Game) Move(index int) Status {
	n := g.floodFill(index)
	if n < 2 {
		return Status{
			Board:    g.board,
			Score:    g.score,
			GameOver: g.IsGameOver(),
		}
	}

	g.applyGravityAndShiftRight()
	g.score += calculateMoveScore(n)
	gameOver := g.IsGameOver()

	if gameOver {
		remainingPieces := 0
		for _, p := range g.board {
			if p != White {
				remainingPieces++
			}
		}
		g.score += calculateRemainingPiecesScore(remainingPieces)
	}

	return Status{
		Board:    g.board,
		Score:    g.score,
		GameOver: gameOver,
	}
}

// IsGameOver checks if the game is over by determining if there are any valid moves left.
// A valid move requires at least two connected pieces of the same color.
// Returns true if the game is over (no valid moves), false otherwise.
func (g *Game) IsGameOver() bool {
	for i := 0; i < len(g.board); i++ {
		if g.board[i] == White {
			continue
		}
		if len(g.GetConnectedPieces(i)) > 1 {
			return false
		}
	}

	return true
}

// calculateMoveScore computes the score for removing a group of connected board.
// The score is calculated as n * (n-1), where n is the number of pieces removed.
// If less than 2 pieces are removed, the score is 0.
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
// If more than the threshold number of pieces remains, no bonus is awarded.
func calculateRemainingPiecesScore(remainingPieces int) int {
	const threshold = 20
	if remainingPieces <= threshold {
		p := threshold - remainingPieces
		return p * p
	}

	return 0
}

// GetConnectedPieces returns a slice of all indices that are connected to the piece at the given index.
// It is a non-destructive operation that doesn't modify the game state.
func (g *Game) GetConnectedPieces(index int) []int {
	if index < 0 || index >= len(g.board) {
		return []int{}
	}

	target := g.board[index]
	if target == White {
		return []int{}
	}

	var connectedIndices []int
	stack := []int{index}
	// Pre-allocate the visited map with a reasonable capacity
	visited := make([]bool, g.rows*g.cols)

	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i < 0 || i >= len(g.board) || g.board[i] != target || visited[i] {
			continue
		}
		visited[i] = true
		connectedIndices = append(connectedIndices, i)

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
	return connectedIndices
}

// floodFill performs a flood-fill operation starting at the given index,
// marking connected board of the same color as White.
// Returns the number of connected pieces modified.
// If less than two connected pieces are found, no changes are made.
func (g *Game) floodFill(index int) int {
	connectedPieces := g.GetConnectedPieces(index)
	if len(connectedPieces) < 2 {
		return 0
	}

	for _, i := range connectedPieces {
		g.board[i] = White
	}

	return len(connectedPieces)
}

// ApplyGravityAndShiftRight modifies the board in-place to apply vertical gravity
// and then right-align non-empty columns.
func (g *Game) applyGravityAndShiftRight() {
	// Gravity Phase: shift non-'w' characters down in each column
	for col := 0; col < g.cols; col++ {
		writeRow := g.rows - 1
		for row := g.rows - 1; row >= 0; row-- {
			index := row*g.cols + col
			if g.board[index] != White {
				g.board[writeRow*g.cols+col] = g.board[index]
				if writeRow != row {
					g.board[index] = White
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
			if g.board[row*g.cols+col] != White {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			if writeCol != col {
				// Copy column to new position
				for row := 0; row < g.rows; row++ {
					g.board[row*g.cols+writeCol] = g.board[row*g.cols+col]
					g.board[row*g.cols+col] = White
				}
			}
			writeCol--
		}
	}
}
