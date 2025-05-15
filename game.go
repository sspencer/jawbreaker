package jawbreaker

import (
	"encoding/base64"
	"errors"
	"fmt"
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
	Board           Board
	Score           int
	Bonus           int
	RemainingPieces int
	GameOver        bool
}

// Piece constants used to represent different colored board in the game.
// Each color is represented by a single byte character.
const (
	PieceSpace            = 32
	White                 = 0 // "empty" or removed piece
	Purple                = 32
	Blue                  = 64
	Green                 = 96
	Red                   = 128
	Yellow                = 160
	PowerUpX              = 1
	PowerUpPlus           = 2
	PowerUpCircle         = 3
	PowerUpRect           = 4
	PowerExtraFill        = 5
	PowerExtraRotateRight = 6
	PowerExtraRotateLeft  = 7
)

var (
	// colorPieces is a slice of all available colored board used for random piece generation.
	colorPieces    = []Piece{Purple, Blue, Green, Red, Yellow}
	allColorPieces = []Piece{White, Purple, Blue, Green, Red, Yellow}
	powerPieces    = []Piece{PowerUpX, PowerUpPlus, PowerUpCircle, PowerUpRect}
	extraPieces    = []Piece{PowerExtraFill, PowerExtraRotateRight, PowerExtraRotateLeft}
	allPowerPieces = append(powerPieces, extraPieces...)

	// ErrGameSize is returned when attempting to restore a game with a board size
	// that doesn't match the expected dimensions.
	ErrGameSize = errors.New("game board does not match expected size")
)

// Color returns each piece as a string representing its color.
// It returns "white" for any other piece.
func (c Piece) Color() int {
	return int(c/PieceSpace) * PieceSpace
}

func (c Piece) ColorName() string {
	switch c.Color() {
	case Purple:
		return "purple"
	case Blue:
		return "blue"
	case Green:
		return "green"
	case Red:
		return "red"
	case Yellow:
		return "yellow"
	default:
		if c.PowerUp() == 0 {
			return "white"
		}
		return "gray"
	}
}

func (c Piece) PowerUp() int {
	return int(c) % PieceSpace
}

func (c Piece) PowerUpName() string {
	// https://www.w3schools.com/charsets/ref_utf_symbols.asp
	switch c.PowerUp() {
	case PowerUpX:
		return "&#10005;"
	case PowerUpPlus:
		return "&#43;" //"&#9532;"
	case PowerUpCircle:
		return "&#1054;"
	case PowerUpRect:
		return "&#127020;"
	case PowerExtraFill:
		return "&#9734;"
	case PowerExtraRotateLeft:
		return "&#8617;"
	case PowerExtraRotateRight:
		return "&#8618;"
	default:
		return ""
	}
}

func (c Piece) IsExtraPowerUp() bool {
	powerUp := c.PowerUp()
	return powerUp == PowerExtraFill || powerUp == PowerExtraRotateRight || powerUp == PowerExtraRotateLeft
}

// NewGame initializes a new Game instance with the given rows and columns, populating the board with random board.
func NewGame(rows, cols int) *Game {
	return NewGameWithOptions(rows, cols, nil)
}

type GameOptions struct {
	powerUps bool
}

func (o *GameOptions) PowerUps() *GameOptions {
	o.powerUps = true
	return o
}

// NewGameWithOptions initializes a new Game instance with the given rows and columns, populating the board with random board.
func NewGameWithOptions(rows, cols int, opts *GameOptions) *Game {
	size := rows * cols
	board := make(Board, size)

	for i := range board {
		board[i] = colorPieces[rand.IntN(len(colorPieces))]
	}

	// Number of PowerUps to add (current 27) must be less than board size
	powerUpSize := len(colorPieces)*len(powerPieces) + len(allPowerPieces)

	if opts != nil && opts.powerUps && powerUpSize < size {
		all := AllIndices(size)
		Shuffle(all)
		index := 0

		for a := range allColorPieces {
			for p := range allPowerPieces {
				i := all[index]
				color := allColorPieces[a]
				powerUp := allPowerPieces[p]

				if powerUp.IsExtraPowerUp() && color != White {
					continue
				}

				board[i] = color + powerUp
				index++
			}
		}
	}

	return &Game{board: board, rows: rows, cols: cols}
}

// RestoreGame restores a game state based on the given board string, dimensions, and score.
// It returns a pointer to a Game instance or an error if the board size does not match the expected size.
func RestoreGame(encodedBoard string, rows, cols, score int) (*Game, error) {
	b, err := base64ToBytes(encodedBoard)
	if err != nil {
		return nil, err

	}

	if len(b) != rows*cols {
		return nil, ErrGameSize
	}

	board := convertBytesToBoard(b)

	return &Game{board: board, rows: rows, cols: cols, score: score}, nil
}

// Board returns the current state of the board as bytes
func (g *Game) Board() Board {
	return g.board
}

func (b Board) Base64() string {
	return bytesToBase64(convertBoardToBytes(b))
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
	bonusScore := 0
	remainingPieces := 0

	if gameOver {
		for _, p := range g.board {
			if p != White {
				remainingPieces++
			}
		}
		bonusScore = calculateBonusScore(remainingPieces)
		g.score += bonusScore
	}

	return Status{
		Board:           g.board,
		Bonus:           bonusScore,
		Score:           g.score,
		RemainingPieces: remainingPieces,
		GameOver:        gameOver,
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

// calculateBonusScore computes a bonus score at the end of the game
// based on the number of pieces remaining on the board.
func calculateBonusScore(remainingPieces int) int {
	const threshold = 10 // Bonus if 10 or fewer pieces remain
	if remainingPieces == 0 {
		// Bonus for clearing the board
		return 2000
	} else if remainingPieces <= threshold {
		// Scaled bonus
		return (threshold - remainingPieces + 1) * 100
	}

	return 0
}

// GetConnectedPieces returns a slice of all indices that are connected to the piece at the given index.
// It is a non-destructive operation that doesn't modify the game state.
func (g *Game) GetConnectedPieces(index int) []int {
	if index < 0 || index >= len(g.board) {
		return nil
	}

	target := g.board[index]
	if target == White {
		return nil
	}

	var connectedIndices []int
	stack := []int{index}
	// Pre-allocate the visited map with a reasonable capacity
	visited := make([]bool, g.rows*g.cols)

	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i < 0 || i >= len(g.board) || g.board[i].Color() != target.Color() || visited[i] {
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

	if len(connectedIndices) < 2 {
		return nil
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

// AllIndices returns all possible indices into the game board for
// the specified size.
func AllIndices(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}

	return slice
}

// Shuffle randomizes the order of elements in a slice of integers
func Shuffle(slice []int) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func base64ToBytes(encoded string) ([]byte, error) {
	// Try standard decoding first
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err == nil {
		return decoded, nil
	}

	// If standard fails, try URL encoding
	decoded, err = base64.URLEncoding.DecodeString(encoded)
	if err == nil {
		return decoded, nil
	}

	// If both fail, return the error
	return nil, fmt.Errorf("input is not valid base64 (std or url)")
}

func bytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// convertBytesToBoard manually converts a []byte to a Board.
// It creates a new Board and copies elements, converting each byte to a Piece.
func convertBytesToBoard(byteSlice []byte) Board {
	if byteSlice == nil {
		// Depending on desired semantics, you might return Board{} for nil input
		// but returning nil is often preferred if the input can be nil.
		return nil
	}

	// Create a new Board with the same length as the byteSlice.
	gameBoard := make(Board, len(byteSlice))

	// Iterate over the byteSlice, convert each byte to Piece, and assign to gameBoard.
	for i, bVal := range byteSlice {
		gameBoard[i] = Piece(bVal) // Convert byte to Piece
	}
	return gameBoard
}

// convertBoardToBytes manually converts a Board to a []byte.
// It creates a new []byte and copies elements, converting each Piece to a byte.
func convertBoardToBytes(gameBoard Board) []byte {
	if gameBoard == nil {
		return nil
	}

	// Create a new []byte with the same length as the gameBoard.
	byteSlice := make([]byte, len(gameBoard))

	// Iterate over the gameBoard, convert each Piece to byte, and assign to byteSlice.
	for i, pVal := range gameBoard {
		byteSlice[i] = byte(pVal) // Convert Piece to byte
	}
	return byteSlice
}
