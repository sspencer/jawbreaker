package jawbreaker

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
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
	board           Board // board represents the game board as a flat array of color bytes
	rows            int   // number of rows in the game board
	cols            int   // number of columns in the game board
	score           int   // current score of the game
	lastScore       int   // stored for convenience
	bestScore       int   // stored for convenience
	bonus           int
	remainingPieces int
}

type Status struct {
	Board     Board
	Score     int
	Bonus     int
	Last      int
	Best      int
	Remaining int
	GameOver  bool
}

func (s Status) String() string {
	return fmt.Sprintf("Score: %d, Bonus: %d, Last: %d, Best: %d, Remaining: %d, GameOver: %t", s.Score, s.Bonus, s.Last, s.Best, s.Remaining, s.GameOver)
}

type direction struct {
	sr int // start row
	sc int // start column
	dc int // direction column
	dr int // direction row
}

var xDirection = []direction{
	{sr: -1, sc: -1, dr: -1, dc: -1},
	{sr: -1, sc: 1, dr: -1, dc: 1},
	{sr: 1, sc: -1, dr: 1, dc: -1},
	{sr: 1, sc: 1, dr: 1, dc: 1},
}

var plusDirection = []direction{
	{sr: -1, sc: 0, dr: -1, dc: 0},
	{sr: 1, sc: 0, dr: 1, dc: 0},
	{sr: 0, sc: -1, dr: 0, dc: -1},
	{sr: 0, sc: 1, dr: 0, dc: 1},
}

var rectDirection = []direction{
	{sr: -2, sc: -2, dc: 0, dr: 1},
	{sr: 2, sc: -2, dc: 1, dr: 0},
	{sr: 2, sc: 2, dc: 0, dr: -1},
	{sr: -2, sc: 2, dc: -1, dr: 0},
}

var circleDirection = []direction{
	{sr: -1, sc: -3, dc: 0, dr: 1},
	{sr: 3, sc: -1, dc: 1, dr: 0},
	{sr: 1, sc: 3, dc: 0, dr: -1},
	{sr: -3, sc: 1, dc: -1, dr: 0},
}

// Piece constants used to represent different colored board in the game.
// Each color is represented by a single byte character.
const (
	PieceSpace       = 32
	White            = 0 // "empty" or removed piece
	Purple           = 32
	Blue             = 64
	Green            = 96
	Red              = 128
	Yellow           = 160
	PowerX           = 1
	PowerPlus        = 2
	PowerCircle      = 3
	PowerRect        = 4
	PowerFill        = 5
	PowerRotateRight = 6
	PowerRotateLeft  = 7
	PieceSelected    = 16
)

var (
	// colorPieces is a slice of all available colored board used for random piece generation.
	colorPieces    = []Piece{Purple, Blue, Green, Red, Yellow}
	allColorPieces = []Piece{White, Purple, Blue, Green, Red, Yellow}
	powerPieces    = []Piece{PowerX, PowerPlus, PowerCircle, PowerRect}
	extraPieces    = []Piece{PowerFill, PowerRotateRight, PowerRotateLeft}
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

func (c Piece) Power() int {
	power := c % PieceSpace
	if power >= PieceSelected {
		return int(power - PieceSelected)
	}

	return int(power)
}

func (c Piece) isSelected() bool {
	power := c % PieceSpace
	return power >= PieceSelected
}

func (c Piece) isPowerExtra() bool {
	power := c.Power()
	return power == PowerFill || power == PowerRotateRight || power == PowerRotateLeft
}

func (c Piece) isPowerDirection() bool {
	power := c.Power()
	return power == PowerX || power == PowerPlus || power == PowerCircle || power == PowerRect
}

func (c Piece) isPowerGlobDirection() bool {
	power := c.Power()
	return c.Color() == White && (power == PowerX || power == PowerPlus || power == PowerCircle || power == PowerRect)
}

type GameOptions struct {
	powerUps   bool
	startEmpty bool
}

func (o *GameOptions) PowerUps() *GameOptions {
	o.powerUps = true
	return o
}

func (o *GameOptions) StartEmpty() *GameOptions {
	o.startEmpty = true
	return o
}

// NewGame initializes a new Game instance with the given rows and columns, populating the board with random board.
func NewGame(rows, cols int) *Game {
	return NewGameWithOptions(rows, cols, nil)
}

// NewGameWithOptions initializes a new Game instance with the given rows and columns, populating the board with random board.
func NewGameWithOptions(rows, cols int, opts *GameOptions) *Game {
	rows = max(rows, 6)
	cols = max(cols, 6)

	size := rows * cols
	board := make(Board, size)
	if opts == nil {
		opts = &GameOptions{}
	}

	game := &Game{board: board, rows: rows, cols: cols}
	if opts.startEmpty {
		return game
	}

	// fill game board with random pieces
	game.fillEmpties()

	// fill the game board with random power-ups
	if opts.powerUps {
		all := ShuffledIndices(size)
		index := 0

		// there are 27 power ups (20 for the colors + 7 specials)
		for a := range allColorPieces {
			for p := range allPowerPieces {
				i := all[index]
				color := allColorPieces[a]
				powerUp := allPowerPieces[p]

				if powerUp.isPowerExtra() && color != White {
					continue
				}

				game.board[i] = color + powerUp
				index++
			}
		}
	}

	return game
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
	for i, p := range board {
		if p.isSelected() {
			board[i] -= PieceSelected
		}
	}

	return &Game{board: board, rows: rows, cols: cols, score: score}, nil
}

// Board returns the current state of the board as bytes
func (g *Game) Board() Board {
	return g.board
}

// BoardWithConnections returns the board with the selected piece and all of its
// connections highlighted.
func (g *Game) BoardWithConnections(index int) Board {
	connected := g.GetConnectedPieces(index)

	board := make(Board, len(g.board))
	copy(board, g.board)

	for _, idx := range connected {
		board[idx] += PieceSelected
	}

	return board
}

func (b Board) Base64() string {
	return bytesToBase64(convertBoardToBytes(b))
}

// Score returns the current score of the game.
func (g *Game) Score() int {
	return g.score
}

func (g *Game) SetLastScore(score int) {
	g.lastScore = score
}

func (g *Game) SetBestScore(score int) {
	g.bestScore = score
}

func (g *Game) powerMove(index int) Status {
	powerUp := g.board[index].Power()

	g.board[index] = White

	switch powerUp {
	case PowerFill:
		g.fillEmpties()
	case PowerRotateLeft:
		g.rotateLeft()
	case PowerRotateRight:
		g.rotateRight()
	default:
		log.Fatalf("Unknown power move: %d\n", powerUp)
	}

	g.applyGravityAndShiftRight()

	return Status{
		Board:    g.board,
		Score:    g.score,
		GameOver: g.IsGameOver(),
	}
}

// Move removes connected pieces of the same color from the board and updates the score.
// It takes the index of the clicked piece and removes all connected pieces of the same color.
// If less than 2 connected pieces are found, no pieces are removed and the score remains unchanged.
// After removing board, gravity is applied to make pieces fall down, and empty columns are shifted right.
// If the game is over after the move, a bonus score is added based on the number of remaining pieces.
func (g *Game) Move(index int) Status {
	if g.board[index].isPowerExtra() {
		return g.powerMove(index)
	}

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
	g.bonus = 0
	g.remainingPieces = 0

	if gameOver {
		for _, p := range g.board {
			if p != White {
				g.remainingPieces++
			}
		}
		g.bonus = calculateBonusScore(g.remainingPieces)
		g.score += g.bonus
		g.lastScore = g.score
		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}

	return Status{
		Board:     g.board,
		Bonus:     g.bonus,
		Score:     g.score,
		Last:      g.lastScore,
		Best:      g.bestScore,
		Remaining: g.remainingPieces,
		GameOver:  gameOver,
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

	connectedIndices := g.getPowerUpConnections(index)

	stack := []int{index}
	// Pre-allocate the visited map with a reasonable capacity
	visited := make([]bool, g.rows*g.cols)

	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i < 0 || i >= len(g.board) || g.board[i].Color() != target.Color() || visited[i] || g.board[i] == White {
			continue
		}
		visited[i] = true
		connectedIndices[i] = true

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
		if !(len(connectedIndices) == 1 && connectedIndices[index] && g.board[index].isPowerExtra()) {
			return nil
		}
	}

	indices := make([]int, 0, len(connectedIndices))
	for i := range connectedIndices {
		indices = append(indices, i)
	}

	return indices
}

func (g *Game) getPowerUpConnections(index int) map[int]bool {
	connectedIndices := make(map[int]bool)
	powerUp := g.board[index].Power()
	switch powerUp {
	case PowerX:
		return g.getXConnections(index)
	case PowerPlus:
		return g.getPlusConnections(index)
	case PowerRect:
		return g.getRectConnections(index)
	case PowerCircle:
		return g.getCircularConnections(index)
	case PowerFill, PowerRotateRight, PowerRotateLeft:
		connectedIndices[index] = true
		return connectedIndices
	default:
		return connectedIndices
	}
}

func (g *Game) getXConnections(index int) map[int]bool {
	maxIters := int(math.Ceil(math.Sqrt(float64(g.rows*g.rows + g.cols*g.cols))))
	return g.getConnectionsWithDirections(index, xDirection, maxIters)
}

func (g *Game) getPlusConnections(index int) map[int]bool {
	maxIters := max(g.rows, g.cols)
	return g.getConnectionsWithDirections(index, plusDirection, maxIters)
}

func (g *Game) getRectConnections(index int) map[int]bool {
	return g.getConnectionsWithDirections(index, rectDirection, 4)
}

func (g *Game) getCircularConnections(index int) map[int]bool {
	c1 := g.getConnectionsWithDirections(index, rectDirection, 1)
	c2 := g.getConnectionsWithDirections(index, circleDirection, 3)
	for k, v := range c1 {
		c2[k] = v
	}
	return c2
}

func (g *Game) getConnectionsWithDirections(index int, directions []direction, maxIters int) map[int]bool {
	target := g.board[index]
	connectedIndices := make(map[int]bool)
	startRow, startCol := g.pointFromIndex(index)
	for iter := 0; iter < maxIters; iter++ {
		for _, dir := range directions {
			r := startRow + dir.sr + dir.dr*iter
			c := startCol + dir.sc + dir.dc*iter
			if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
				currentIndex := r*g.rows + c
				piece := g.board[currentIndex]
				color := piece.Color()
				if color != White && (color == target.Color() || target.isPowerGlobDirection()) {
					connectedIndices[currentIndex] = true
				}
			}
		}
	}

	return connectedIndices
}

func (g *Game) pointFromIndex(index int) (int, int) {
	col := index % g.rows
	row := index / g.rows

	return row, col
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

func ShuffledIndices(size int) []int {
	all := AllIndices(size)
	shuffle(all)
	return all
}

func AllIndices(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}

	return slice
}

// Shuffle randomizes the order of elements in a slice of integers
func shuffle(slice []int) {
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

	// If the standard fails, try URL encoding
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
		// Depending on desired semantics, you might return Board{} for nil input,
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

// AnimateBoard applies up to numValues from the game board at given indices into a new slice.
// Both the game board and indices must have the same length.
func (g *Game) AnimateBoard(indices []int, numValues int) (Board, error) {
	size := len(g.board)
	if len(indices) != size {
		return nil, errors.New("goal and indices slices must have the same length")
	}

	// Create a zero-initialized result slice
	result := make(Board, size)

	// Apply up to numValues from goal using indices
	count := 0
	for i := 0; i < size && count < numValues; i++ {
		idx := indices[i]
		if idx >= 0 && idx < size {
			result[idx] = g.board[idx]
			count++
		}
	}

	return result, nil
}

func (g *Game) fillEmpties() {
	for i := range g.board {
		if g.board[i] == White {
			g.board[i] = colorPieces[rand.IntN(len(colorPieces))]
		}
	}
}

func (g *Game) rotateRight() {
	board, err := rotateSlice(g.board, g.rows, g.cols, -90)
	if err != nil {
		slog.Error("rotateRight", "error", err.Error())
		return
	}

	g.board = board
}

func (g *Game) rotateLeft() {
	board, err := rotateSlice(g.board, g.rows, g.cols, 90)
	if err != nil {
		slog.Error("rotateLeft", "error", err.Error())
		return
	}

	g.board = board
}

// RotateSlice rotates a 1D slice (row-major) by 90 or -90 degrees.
// `rows` and `cols` are the dimensions of the input matrix.
// Positive degree means clockwise, negative means counter-clockwise.
func rotateSlice(data Board, rows, cols, degree int) (Board, error) {
	if len(data) != rows*cols {
		return nil, errors.New("invalid dimensions for the given slice length")
	}

	result := make(Board, len(data))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			val := data[r*cols+c]
			if degree > 0 {
				// 90 degrees
				// New row becomes the column index from the bottom
				newRow := c
				newCol := rows - 1 - r
				result[newRow*rows+newCol] = val
			} else {
				// -90 degrees
				newRow := cols - 1 - c
				newCol := r
				result[newRow*rows+newCol] = val
			}
		}
	}

	// New dimensions are transposed
	return result, nil
}
