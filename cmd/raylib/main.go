package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	jb "github.com/sspencer/jawbreaker"
)

const (
	gameSize      = 14
	blockSize     = 23
	windowPadding = 30
	gridOuter     = 2
	gridInner     = 1
)

var (
	background = rl.NewColor(43, 60, 80, 255)
)

func main() {
	g := jb.New(jb.WithSize(gameSize), jb.WithoutExtras())

	rl.SetConfigFlags(rl.FlagVsyncHint)
	var windowSize int32 = blockSize*gameSize + (windowPadding+gridOuter)*2

	rl.InitWindow(windowSize*2, windowSize*2, "Jawbreaker")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		if g.GameOver() {
			if rl.IsKeyPressed(rl.KeySpace) {
				g.Restart()
			}
		} else {
			if g.CanUndo() && (rl.IsKeyPressed(rl.KeyZ) || rl.IsKeyPressed(rl.KeyU)) {
				fmt.Println("Undo")
				g.Undo()
			}

			if rl.IsKeyPressed(rl.KeyR) {
				fmt.Println("Restart")
				g.Restart()
			}
		}

		zoom := float32(rl.GetScreenHeight()) / float32(windowSize)
		camera := rl.Camera2D{Zoom: zoom}

		boardSize := float32(gameSize)*blockSize - gridInner + gridOuter*2
		boardRec := rl.Rectangle{X: windowPadding, Y: windowPadding, Width: boardSize, Height: boardSize}

		mousePos := boardCoords(rl.GetMousePosition(), zoom)
		connections := g.Connections(mousePos)
		numConnected := len(connections)
		//madeMove := false
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && numConnected > 1 {
			//g.setUndo()
			//madeMove = true
			fmt.Printf("Make a Move: %v = %d\n", mousePos, numConnected)
			g.Move(connections)
		}

		rl.BeginDrawing()

		rl.BeginMode2D(camera)
		rl.ClearBackground(background)

		rl.DrawRectangleRounded(boardRec, 0.02, 16, rl.NewColor(20, 24, 33, 255))

		drawBlocks(g)
		drawSelection(connections)
		drawScore(g, boardRec)

		if g.GameOver() {
			drawGameOver(g, boardRec)
		}

		rl.EndMode2D()
		rl.EndDrawing()
	}
}

func boardCoords(mouse rl.Vector2, zoom float32) jb.Vec2i {
	gridX := mouse.X/zoom - windowPadding - gridOuter
	gridY := mouse.Y/zoom - windowPadding - gridOuter

	// Check if within board bounds
	if gridX < 0 ||
		gridX >= float32(blockSize*gameSize) ||
		gridY < 0 ||
		gridY >= float32(blockSize*gameSize) {
		return jb.Vec2i{X: -1, Y: -1}
	}

	boardX := gridX / blockSize
	boardY := gridY / blockSize

	if boardX >= float32(gameSize) || boardY >= float32(gameSize) {
		return jb.Vec2i{X: -1, Y: -1}
	}

	return jb.Vec2i{X: int(math.Floor(float64(boardX))), Y: int(math.Floor(float64(boardY)))}
}
