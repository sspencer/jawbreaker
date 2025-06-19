package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GameSize      = 14
	BlockSize     = 23
	WindowPadding = 30
	GridOuter     = 2
	GridInner     = 1
)

func main() {

	rl.SetConfigFlags(rl.FlagVsyncHint)

	g := NewGame(GameSize)
	var windowSize int32 = BlockSize*GameSize + (WindowPadding+GridOuter)*2

	rl.InitWindow(windowSize*2, windowSize*2, "Jawbreaker")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	g.Restart()

	for !rl.WindowShouldClose() {

		if g.gameOver {
			if rl.IsKeyPressed(rl.KeySpace) {
				g.Restart()
			}
		}

		zoom := float32(rl.GetScreenHeight()) / float32(windowSize)
		camera := rl.Camera2D{Zoom: zoom}

		// board_size := f32(NUM_BLOCKS) * BLOCK_SIZE -GRID_INNER + GRID_OUTER * 2
		boardSize := float32(GameSize)*BlockSize - GridInner + GridOuter*2
		boardRec := rl.Rectangle{WindowPadding, WindowPadding, boardSize, boardSize}

		mousePos := g.boardCoords(rl.GetMousePosition(), zoom)

		connections := g.blockSelection(mousePos)
		numConnected := len(connections)
		g.madeMove = false
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && numConnected > 1 {
			g.setUndo()
			g.madeMove = true
			fmt.Printf("Make a Move: %v = %d\n", mousePos, numConnected)
			g.makeMove(connections)
		}

		rl.BeginDrawing()

		rl.BeginMode2D(camera)
		rl.ClearBackground(Background)

		rl.DrawRectangleRounded(boardRec, 0.02, 16, rl.NewColor(20, 24, 33, 255))
		g.drawBlocks()
		g.drawSelection(connections)

		rl.EndMode2D()
		rl.EndDrawing()
	}
}

func (g *Game) setUndo() {
	dst := make([][]GameBlock, len(g.board))
	for i := range g.board {
		dst[i] = make([]GameBlock, len(g.board[i]))
		copy(dst[i], g.board[i])
	}

	g.undo = dst
	g.undoScore = g.gameScore
	g.canUndo = true
}
