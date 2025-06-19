package jawbreaker

import (
	"math"
	"math/rand/v2"
)

type Tile byte

type Vec2i struct {
	X, Y int
}

const (
	defaultSize   = 14
	defaultExtras = true
)

type Game struct {
	width     int
	height    int
	board     [][]Tile
	undo      [][]Tile
	undoScore int
	gameOver  bool
	canUndo   bool
	gameScore int
	gameBonus int
	lastScore int
	bestScore int
	extras    bool
	//madeMove  bool
}

type TileSelection struct {
	Pos Vec2i
	Dir Vec2i
}

const (
	TileEmpty Tile = iota
	TilePurple
	TileBlue
	TileGreen
	TileRed
	TileYellow
	TileTimes
	TilePlus
	TileMinus
	TilePipe
	TileRect
	TileCircle
)

var (
	tileSelection = map[Tile][]TileSelection{
		TileTimes: {
			{Vec2i{-1, -1}, Vec2i{-1, -1}},
			{Vec2i{1, -1}, Vec2i{1, -1}},
			{Vec2i{-1, 1}, Vec2i{-1, 1}},
			{Vec2i{1, 1}, Vec2i{1, 1}},
		},
		TilePlus: {
			{Vec2i{0, -1}, Vec2i{0, -1}},
			{Vec2i{0, 1}, Vec2i{0, 1}},
			{Vec2i{-1, 0}, Vec2i{-1, 0}},
			{Vec2i{1, 0}, Vec2i{1, 0}},
		},
		TileMinus: {
			{Vec2i{-1, 0}, Vec2i{-1, 0}},
			{Vec2i{1, 0}, Vec2i{1, 0}},
		},
		TilePipe: {
			{Vec2i{0, -1}, Vec2i{0, -1}},
			{Vec2i{0, 1}, Vec2i{0, 1}},
		},
		TileRect: {
			{Vec2i{-2, -2}, Vec2i{0, 1}},
			{Vec2i{-2, 2}, Vec2i{1, 0}},
			{Vec2i{2, 2}, Vec2i{0, -1}},
			{Vec2i{2, -2}, Vec2i{-1, 0}},
		},
		TileCircle: {
			{Vec2i{-3, -1}, Vec2i{0, 1}},
			{Vec2i{-1, 3}, Vec2i{1, 0}},
			{Vec2i{3, 1}, Vec2i{0, -1}},
			{Vec2i{1, -3}, Vec2i{-1, 0}},
		},
	}
)

type Option func(*Game)

func WithExtras() Option {
	return func(g *Game) {
		g.extras = true
	}
}

func WithoutExtras() Option {
	return func(g *Game) {
		g.extras = false
	}
}

func WithWidth(width int) Option {
	return func(g *Game) {
		g.width = width
	}
}

func WithHeight(height int) Option {
	return func(g *Game) {
		g.height = height
	}
}

func WithSize(size int) Option {
	return func(g *Game) {
		g.width = size
		g.height = size
	}
}

func WithScore(score int) Option {
	return func(g *Game) {
		g.gameScore = score
	}
}

func WithLastScore(score int) Option {
	return func(g *Game) {
		g.lastScore = score
	}
}

func WithBestScore(score int) Option {
	return func(g *Game) {
		g.bestScore = score
	}
}

func New(opts ...Option) *Game {
	g := &Game{
		width:  defaultSize,
		height: defaultSize,
		extras: defaultExtras,
	}

	for _, opt := range opts {
		opt(g)
	}

	g.board = make([][]Tile, g.height)
	g.undo = make([][]Tile, g.height)

	for y := 0; y < g.height; y++ {
		g.board[y] = make([]Tile, g.width)
		g.undo[y] = make([]Tile, g.width)
	}

	g.Restart()

	return g
}

//func NewGame(width, height int, extras bool) *Game {
//	g := &Game{
//		width:     width,
//		height:    height,
//		board:     make([][]Tile, height),
//		undo:      make([][]Tile, height),
//		gameOver:  false,
//		gameScore: 0,
//		gameBonus: 0,
//		lastScore: 0,
//		extras:    extras,
//	}
//
//	for y := 0; y < height; y++ {
//		g.board[y] = make([]Tile, width)
//		g.undo[y] = make([]Tile, width)
//	}
//
//	g.Restart()
//
//	return g
//}

func (g *Game) Restart() {
	tiles := []Tile{TilePurple, TileBlue, TileGreen, TileRed, TileYellow}
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.board[y][x] = tiles[rand.IntN(len(tiles))]
		}
	}

	if g.extras {
		extras := []Tile{TilePlus, TileMinus, TileTimes, TileRect, TileCircle, TilePipe}
		pos := g.getRandomPositions(len(extras))
		for i, p := range pos {
			g.board[p.Y][p.X] = extras[i]
		}
	}

	g.lastScore = g.gameScore
	if g.gameScore > g.bestScore {
		g.bestScore = g.gameScore
	}

	g.gameScore = 0
	g.gameOver = false
	g.gameBonus = 0
	g.canUndo = false
}

func (g *Game) GameOver() bool {
	return g.gameOver
}

func (g *Game) Score() int {
	return g.gameScore
}

func (g *Game) Bonus() int {
	return g.gameBonus
}

func (g *Game) RemainingTiles() int {
	return g.countTiles()
}

func (g *Game) Board() [][]Tile {
	return g.board
}

func (g *Game) CanUndo() bool {
	return g.canUndo
}

func (g *Game) Undo() {
	if g.canUndo {
		for y := 0; y < g.height; y++ {
			for x := 0; x < g.width; x++ {
				g.board[y][x] = g.undo[y][x]
			}
		}
	}

	g.gameScore = g.undoScore
	g.canUndo = false
}

func (g *Game) LastScore() int {
	return g.lastScore
}

func (g *Game) BestScore() int {
	return g.bestScore
}

func (g *Game) getRandomPositions(n int) []Vec2i {
	pos := make([]Vec2i, 0, n)
	used := make(map[Vec2i]bool)

	for i := 0; i < n; i++ {
		x := rand.IntN(g.width)
		y := rand.IntN(g.height)
		if _, ok := used[Vec2i{x, y}]; ok {
			continue
		}
		pos = append(pos, Vec2i{x, y})
		used[Vec2i{x, y}] = true
	}

	return pos
}

func (g *Game) Move(sel []Vec2i) {
	numConnected := len(sel)
	if numConnected < 2 {
		return
	}

	// UNDO
	g.canUndo = true
	g.undoScore = g.gameScore
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.undo[y][x] = g.board[y][x]
		}
	}

	// MOVE
	g.gameScore += numConnected * (numConnected - 1)
	for _, s := range sel {
		g.board[s.Y][s.X] = TileEmpty
	}

	g.applyGravity()

	g.gameOver = g.isGameOver()

	if g.gameOver {
		g.canUndo = false
		g.gameBonus = g.calculateBonus()
		g.gameScore += g.gameBonus
		g.lastScore = g.gameScore
		if g.gameScore > g.bestScore {
			g.bestScore = g.gameScore
		}

		// save high score
	}
}

func (g *Game) applyGravity() {
	// Pass 1: Move all non-zero values down within each column
	for col := 0; col < g.width; col++ {
		writePos := g.height - 1

		for row := g.height - 1; row >= 0; row -= 1 {
			if g.board[row][col] != TileEmpty {
				if writePos != row {
					g.board[writePos][col] = g.board[row][col]
					g.board[row][col] = TileEmpty
				}
				writePos -= 1
			}
		}
	}

	// Pass 2: Move entire columns right to fill gaps from empty columns
	writeCol := g.width - 1

	for col := g.width - 1; col >= 0; col -= 1 {
		// Check if this column has any non-empty pieces
		columnHasPieces := false
		for row := 0; row < g.height; row++ {
			if g.board[row][col] != TileEmpty {
				columnHasPieces = true
				break
			}
		}

		// If column has pieces, move it to the write position
		if columnHasPieces {
			if writeCol != col {
				// Move entire column
				for row := 0; row < g.height; row++ {
					g.board[row][writeCol] = g.board[row][col]
					g.board[row][col] = TileEmpty
				}
			}
			writeCol -= 1
		}
	}
}

func (g *Game) isGameOver() bool {
	for y := 0; y < len(g.board); y++ {
		for x := 0; x < len(g.board[y]); x++ {
			if g.board[y][x] != TileEmpty {
				sel := g.Connections(Vec2i{x, y})
				if len(sel) > 1 {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) calculateBonus() int {
	const threshold = 10 // Bonus if 10 or fewer pieces remain
	pieces := g.countTiles()
	if pieces == 0 {
		return 2000
	} else if pieces <= threshold {
		return (threshold - pieces + 1) * 100
	}

	return 0
}

func (g *Game) countTiles() int {
	tiles := 0
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.board[y][x] != TileEmpty {
				tiles++
			}
		}
	}

	return tiles
}

func (g *Game) Connections(start Vec2i) []Vec2i {
	var sel []Vec2i

	if start.X < 0 || start.X >= g.width || start.Y < 0 || start.Y >= g.height {
		return sel
	}

	target := g.board[start.Y][start.X]

	switch target {
	case TileEmpty:
		return sel
	case TilePlus, TileMinus, TileTimes, TileRect, TileCircle, TilePipe:
		return g.connectedExtras(start)
	default:
		return g.connectedColors(start)
	}
}

func (g *Game) connectedColors(start Vec2i) []Vec2i {
	var sel []Vec2i
	target := g.board[start.Y][start.X]
	visited := make(map[Vec2i]bool)

	var stack []Vec2i
	stack = append(stack, start)

	directions := [4]Vec2i{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1] // pop

		if visited[current] {
			continue
		}

		visited[current] = true
		sel = append(sel, current)

		for _, dir := range directions {
			next := Vec2i{current.X + dir.X, current.Y + dir.Y}

			if next.X < 0 || next.X >= g.width || next.Y < 0 || next.Y >= g.height {
				continue
			}

			if g.board[next.Y][next.X] == target && !(visited[next]) {
				stack = append(stack, next)
			}
		}
	}

	if len(sel) > 1 {
		return sel
	}

	return nil
}

func (g *Game) connectedExtras(start Vec2i) []Vec2i {
	target := g.board[start.Y][start.X]
	selection, ok := tileSelection[target]
	if !ok {
		return nil
	}

	maxIters := 0

	switch target {
	case TileRect:
		return g.computeConnectedExtras(start, selection, 4)
	case TileCircle:
		maxIters = 3
		s1 := g.computeConnectedExtras(start, selection, 3)
		sel2, _ := tileSelection[TileRect]
		s2 := g.computeConnectedExtras(start, sel2, 1)
		return append(s1, s2...)
	default:
		maxIters = int(math.Ceil(math.Sqrt(float64(g.width*g.width + g.height*g.height))))
		return g.computeConnectedExtras(start, selection, maxIters)
	}
}

func (g *Game) computeConnectedExtras(start Vec2i, selections []TileSelection, maxIters int) []Vec2i {
	var sel []Vec2i

	startRow := start.Y
	startCol := start.X

	for iter := 0; iter < maxIters; iter++ {
		for _, s := range selections {
			y := startRow + s.Pos.Y + s.Dir.Y*iter
			x := startCol + s.Pos.X + s.Dir.X*iter
			if y < 0 || x < 0 || y >= g.height || x >= g.width {
				continue
			}

			cur := g.board[y][x]
			if cur == TileEmpty {
				continue
			}

			sel = append(sel, Vec2i{x, y})
		}
	}

	if len(sel) > 0 {
		sel = append(sel, start)
	}

	return sel
}

func (g *Game) Rotate(clockwise bool) {
	if clockwise {
		rotateClockwise(g.board)
	} else {
		rotateCounterClockwise(g.board)
	}
}

// rotateClockwise rotates a square 2D slice 90 degrees clockwise in-place
func rotateClockwise(matrix [][]Tile) {
	n := len(matrix)
	if n == 0 || n != len(matrix[0]) {
		return
	}

	// Step 1: Transpose the matrix (swap elements across diagonal)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	// Step 2: Reverse each row
	for i := 0; i < n; i++ {
		for j := 0; j < n/2; j++ {
			matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
		}
	}
}

// RotateCounterClockwise rotates a square 2D slice 90 degrees counter-clockwise in-place
func rotateCounterClockwise(matrix [][]Tile) {
	n := len(matrix)
	if n == 0 || n != len(matrix[0]) {
		panic("matrix must be square and non-empty")
	}

	// Step 1: Reverse each row
	for i := 0; i < n; i++ {
		for j := 0; j < n/2; j++ {
			matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
		}
	}

	// Step 2: Transpose the matrix (swap elements across diagonal)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
}
