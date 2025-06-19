package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameBlock byte

type Vec2i struct {
	X int
	Y int
}

type PowerUpSelection struct {
	Start Vec2i
	Dir   Vec2i
}

const (
	Empty GameBlock = iota
	Purple
	Blue
	Green
	Red
	Yellow
	PowerPlus
	PowerMinus
	PowerTimes
	PowerRect
	PowerCircle
	PowerPipe
)

var (
	Background = rl.NewColor(43, 60, 80, 255)

	blockColorValues = map[GameBlock]rl.Color{
		Empty:       rl.NewColor(35, 38, 44, 255),
		Purple:      rl.DarkPurple,
		Blue:        rl.Blue,
		Green:       rl.Lime,
		Red:         rl.Red,
		Yellow:      rl.Orange,
		PowerPlus:   rl.DarkGray,
		PowerMinus:  rl.DarkGray,
		PowerTimes:  rl.DarkGray,
		PowerRect:   rl.DarkGray,
		PowerCircle: rl.DarkGray,
		PowerPipe:   rl.DarkGray,
	}

	highlightColorValues = map[GameBlock]rl.Color{
		Empty:       rl.NewColor(38, 55, 75, 204),
		Purple:      rl.Purple,
		Blue:        rl.SkyBlue,
		Green:       rl.Green,
		Red:         rl.Pink,
		Yellow:      rl.Yellow,
		PowerPlus:   rl.Gray,
		PowerMinus:  rl.Gray,
		PowerTimes:  rl.Gray,
		PowerRect:   rl.Gray,
		PowerCircle: rl.Gray,
		PowerPipe:   rl.Gray,
	}

	timesSelection = []PowerUpSelection{
		{Vec2i{-1, -1}, Vec2i{-1, -1}},
		{Vec2i{1, -1}, Vec2i{1, -1}},
		{Vec2i{-1, 1}, Vec2i{-1, 1}},
		{Vec2i{1, 1}, Vec2i{1, 1}},
	}

	plusSelection = []PowerUpSelection{
		{Vec2i{0, -1}, Vec2i{0, -1}},
		{Vec2i{0, 1}, Vec2i{0, 1}},
		{Vec2i{-1, 0}, Vec2i{-1, 0}},
		{Vec2i{1, 0}, Vec2i{1, 0}},
	}

	minusSelection = []PowerUpSelection{
		{Vec2i{-1, 0}, Vec2i{-1, 0}},
		{Vec2i{1, 0}, Vec2i{1, 0}},
	}

	pipeSelection = []PowerUpSelection{ // vertical selection
		{Vec2i{0, -1}, Vec2i{0, -1}},
		{Vec2i{0, 1}, Vec2i{0, 1}},
	}

	rectSelection = []PowerUpSelection{
		{Vec2i{-2, -2}, Vec2i{0, 1}},
		{Vec2i{-2, 2}, Vec2i{1, 0}},
		{Vec2i{2, 2}, Vec2i{0, -1}},
		{Vec2i{2, -2}, Vec2i{-1, 0}},
	}

	circleSelection = []PowerUpSelection{
		{Vec2i{-3, -1}, Vec2i{0, 1}},
		{Vec2i{-1, 3}, Vec2i{1, 0}},
		{Vec2i{3, 1}, Vec2i{0, -1}},
		{Vec2i{1, -3}, Vec2i{-1, 0}},
	}
)

type Game struct {
	size      int
	board     [][]GameBlock
	undo      [][]GameBlock
	undoScore int
	gameOver  bool
	canUndo   bool
	gameScore int
	gameBonus int
	lastScore int
	bestScore int
	madeMove  bool
}

func NewGame(n int) *Game {
	g := &Game{
		size:      n,
		board:     make([][]GameBlock, n),
		undo:      make([][]GameBlock, n),
		gameOver:  false,
		canUndo:   false,
		gameScore: 0,
		gameBonus: 0,
		lastScore: 0,
	}

	for i := 0; i < n; i++ {
		g.board[i] = make([]GameBlock, n)
		g.undo[i] = make([]GameBlock, n)
	}

	return g
}

func (g *Game) Restart() {
	colors := []GameBlock{Purple, Blue, Green, Red, Yellow}
	for y := 0; y < g.size; y++ {
		for x := 0; x < g.size; x++ {
			g.board[y][x] = colors[rand.IntN(len(colors))]
		}
	}

	powerUps := []GameBlock{PowerPlus, PowerMinus, PowerTimes, PowerRect, PowerCircle, PowerPipe}
	randomPositions := g.getRandomPositions(len(powerUps))
	for i, p := range randomPositions {
		g.board[p.Y][p.X] = powerUps[i]
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

func (g *Game) getRandomPositions(n int) []Vec2i {
	pos := make([]Vec2i, 0, n)
	used := make(map[Vec2i]bool)

	for i := 0; i < n; i++ {
		x := rand.IntN(g.size)
		y := rand.IntN(g.size)
		if _, ok := used[Vec2i{x, y}]; ok {
			continue
		}
		pos = append(pos, Vec2i{x, y})
		used[Vec2i{x, y}] = true
	}

	return pos
}

func (g *Game) blockSelection(start Vec2i) []Vec2i {
	var sel []Vec2i

	if start.X < 0 || start.X >= g.size || start.Y < 0 || start.Y >= g.size {
		return sel
	}

	target := g.board[start.Y][start.X]

	switch target {
	case Empty:
		return sel
	case PowerPlus, PowerMinus, PowerTimes, PowerRect, PowerCircle, PowerPipe:
		return g.powerUpBlocksSelection(start)
	default:
		return g.colorBlocksSelection(start)
	}
}

func (g *Game) colorBlocksSelection(start Vec2i) []Vec2i {
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

			if next.X < 0 || next.X >= g.size || next.Y < 0 || next.Y >= g.size {
				continue
			}

			if g.board[next.Y][next.X] == target && !(visited[next]) {
				stack = append(stack, next)
			}
		}
	}

	return sel
}

func (g *Game) powerUpBlocksSelection(start Vec2i) []Vec2i {
	target := g.board[start.Y][start.X]
	switch target {
	case PowerPlus:
		return g.powerSelection(start, plusSelection, g.size)
	case PowerMinus:
		return g.powerSelection(start, minusSelection, g.size)
	case PowerPipe:
		return g.powerSelection(start, pipeSelection, g.size)
	case PowerTimes:
		iters := int(math.Ceil(float64(g.size) * math.Sqrt(2)))
		return g.powerSelection(start, timesSelection, iters)
	case PowerRect:
		return g.powerSelection(start, rectSelection, 4)
	case PowerCircle:
		s1 := g.powerSelection(start, rectSelection, 1)
		s2 := g.powerSelection(start, circleSelection, 3)
		return append(s1, s2...)
	default:
		return nil
	}
}

func (g *Game) powerSelection(start Vec2i, selections []PowerUpSelection, maxIters int) []Vec2i {
	var sel []Vec2i

	startRow := start.Y
	startCol := start.X

	for iter := 0; iter < maxIters; iter++ {
		for _, s := range selections {
			y := startRow + s.Start.Y + s.Dir.Y*iter
			x := startCol + s.Start.X + s.Dir.X*iter
			if y < 0 || x < 0 || y >= g.size || x >= g.size {
				continue
			}

			cur := g.board[y][x]
			//if cur == Empty || blockColorValues[cur] == rl.DarkGray {
			if cur == Empty {
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

func (g *Game) boardCoords(mouse rl.Vector2, zoom float32) Vec2i {
	gridX := mouse.X/zoom - WindowPadding - GridOuter
	gridY := mouse.Y/zoom - WindowPadding - GridOuter

	// Check if within board bounds
	if gridX < 0 ||
		gridX >= float32(BlockSize*g.size) ||
		gridY < 0 ||
		gridY >= float32(BlockSize*g.size) {
		return Vec2i{-1, -1}
	}

	boardX := gridX / BlockSize
	boardY := gridY / BlockSize

	if boardX >= float32(g.size) || boardY >= float32(g.size) {
		return Vec2i{-1, -1}
	}

	return Vec2i{int(math.Floor(float64(boardX))), int(math.Floor(float64(boardY)))}
}

func (g *Game) makeMove(sel []Vec2i) {
	numConnected := len(sel)
	if numConnected < 2 {
		return
	}

	g.gameScore += numConnected * (numConnected - 1)
	for _, s := range sel {
		g.board[s.Y][s.X] = Empty
	}

	g.applyGravity()

	g.gameOver = g.isGameOver()

	if g.gameOver {
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
	for col := 0; col < g.size; col++ {
		writePos := g.size - 1

		for row := g.size - 1; row >= 0; row -= 1 {
			if g.board[row][col] != Empty {
				if writePos != row {
					g.board[writePos][col] = g.board[row][col]
					g.board[row][col] = Empty
				}
				writePos -= 1
			}
		}
	}

	// Pass 2: Move entire columns right to fill gaps from empty columns
	writeCol := g.size - 1

	for col := g.size - 1; col >= 0; col -= 1 {
		// Check if this column has any non-empty pieces
		columnHasPieces := false
		for row := 0; row < g.size; row++ {
			if g.board[row][col] != Empty {
				columnHasPieces = true
				break
			}
		}

		// If column has pieces, move it to the write position
		if columnHasPieces {
			if writeCol != col {
				// Move entire column
				for row := 0; row < g.size; row++ {
					g.board[row][writeCol] = g.board[row][col]
					g.board[row][col] = Empty
				}
			}
			writeCol -= 1
		}
	}
}

func (g *Game) isGameOver() bool {
	for y := 0; y < len(g.board); y++ {
		for x := 0; x < len(g.board[y]); x++ {
			if g.board[y][x] != Empty {
				sel := g.blockSelection(Vec2i{x, y})
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
	pieces := g.countPieces()
	if pieces == 0 {
		return 2000
	} else if pieces <= threshold {
		return (threshold - pieces + 1) * 100
	}

	return 0
}

func (g *Game) countPieces() int {
	pieces := 0
	for y := 0; y < len(g.board); y++ {
		for x := 0; x < len(g.board[y]); x++ {
			if g.board[y][x] != Empty {
				pieces++
			}
		}
	}

	return pieces
}
