package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	jb "github.com/sspencer/jawbreaker"
)

var (
	tileColor = map[jb.Tile]rl.Color{
		jb.TileEmpty:  rl.NewColor(35, 38, 44, 255),
		jb.TilePurple: rl.DarkPurple,
		jb.TileBlue:   rl.Blue,
		jb.TileGreen:  rl.Lime,
		jb.TileRed:    rl.Red,
		jb.TileYellow: rl.Orange,
		jb.TilePlus:   rl.DarkGray,
		jb.TileMinus:  rl.DarkGray,
		jb.TileTimes:  rl.DarkGray,
		jb.TileRect:   rl.DarkGray,
		jb.TileCircle: rl.DarkGray,
		jb.TilePipe:   rl.DarkGray,
	}

	tileHighlight = map[jb.Tile]rl.Color{
		jb.TileEmpty:  rl.NewColor(38, 55, 75, 204),
		jb.TilePurple: rl.Purple,
		jb.TileBlue:   rl.SkyBlue,
		jb.TileGreen:  rl.Green,
		jb.TileRed:    rl.Pink,
		jb.TileYellow: rl.Yellow,
		jb.TilePlus:   rl.Gray,
		jb.TileMinus:  rl.Gray,
		jb.TileTimes:  rl.Gray,
		jb.TileRect:   rl.Gray,
		jb.TileCircle: rl.Gray,
		jb.TilePipe:   rl.Gray,
	}

	shapeColor = rl.White
	textColor  = rl.White
	scoreSize  = int32(12)
	overSize   = int32(16)
)

func drawBlocks(g *jb.Game) {
	board := g.Board()
	for y := range board {
		for x := range board[y] {
			rec := rl.Rectangle{
				X:      float32(x*blockSize + windowPadding + gridOuter),
				Y:      float32(y*blockSize + windowPadding + gridOuter),
				Width:  blockSize - gridInner,
				Height: blockSize - gridInner,
			}

			tile := board[y][x]
			lo := tileColor[tile]
			hi := tileHighlight[tile]

			if tile == jb.TileEmpty {
				rl.DrawRectangleRec(rec, hi)
			} else {
				rl.DrawRectangleGradientEx(rec, hi, lo, lo, lo)
			}

			drawShape(tile, rec)
		}
	}
}

func drawShape(tile jb.Tile, rec rl.Rectangle) {
	switch tile {
	case jb.TilePlus:
		drawPlus(rec)
	case jb.TilePipe:
		drawPipe(rec)
	case jb.TileMinus:
		drawMinus(rec)
	case jb.TileTimes:
		drawTimes(rec)
	case jb.TileRect:
		drawRectangle(rec)
	case jb.TileCircle:
		drawCircle(rec)
	default:
		// ignore
	}
}

func drawSelection(sel []jb.Vec2i) {
	if len(sel) < 2 {
		return
	}

	var (
		color     = rl.NewColor(255, 255, 255, 128)
		xfg       = shapeColor
		xbg       = rl.Red
		shrink    = float32(8.0) // shrink X shape from corners
		thickness = float32(2.0)
		radius    = float32(blockSize-gridInner) * 0.35
	)

	for _, s := range sel {
		rec := rl.Rectangle{
			X:      float32(s.X*blockSize + windowPadding + gridOuter),
			Y:      float32(s.Y*blockSize + windowPadding + gridOuter),
			Width:  blockSize - gridInner,
			Height: blockSize - gridInner,
		}

		rl.DrawRectangleRec(rec, color)

		var (
			upperLeft  = rl.Vector2{X: rec.X + shrink, Y: rec.Y + shrink}
			lowerRight = rl.Vector2{X: rec.X + rec.Width - shrink, Y: rec.Y + rec.Height - shrink}
			lowerLeft  = rl.Vector2{X: rec.X + shrink, Y: rec.Y + rec.Height - shrink}
			upperRight = rl.Vector2{X: rec.X + rec.Width - shrink, Y: rec.Y + shrink}
			centerX    = int32(rec.X + rec.Width/2)
			centerY    = int32(rec.Y + rec.Height/2)
		)

		rl.DrawCircle(centerX, centerY, radius, xbg)
		rl.DrawLineEx(upperLeft, lowerRight, thickness, xfg)
		rl.DrawLineEx(lowerLeft, upperRight, thickness, xfg)
	}
}

func drawPlus(rec rl.Rectangle) {
	centerX := rec.X + rec.Width/2
	centerY := rec.Y + rec.Height/2

	// Calculate dimensions for the plus shape
	thickness := rec.Width / 8  // Thickness of the lines
	length := rec.Width * 3 / 5 // Length of each arm

	// Vertical bar
	vRec := rl.Rectangle{
		X:      centerX - thickness/2,
		Y:      centerY - length/2,
		Width:  thickness,
		Height: length,
	}

	// Horizontal bar
	hRec := rl.Rectangle{
		X:      centerX - length/2,
		Y:      centerY - thickness/2,
		Width:  length,
		Height: thickness,
	}

	rl.DrawRectangleRec(vRec, shapeColor)
	rl.DrawRectangleRec(hRec, shapeColor)
}

func drawPipe(rec rl.Rectangle) {
	centerX := rec.X + rec.Width/2
	centerY := rec.Y + rec.Height/2

	// Calculate dimensions for the plus shape
	thickness := rec.Width / 8  // Thickness of the lines
	length := rec.Width * 3 / 5 // Length of each arm

	// Vertical bar
	vRec := rl.Rectangle{
		X:      centerX - thickness/2,
		Y:      centerY - length/2,
		Width:  thickness,
		Height: length,
	}

	rl.DrawRectangleRec(vRec, shapeColor)
}

func drawMinus(rec rl.Rectangle) {
	centerX := rec.X + rec.Width/2
	centerY := rec.Y + rec.Height/2

	// Calculate dimensions for the plus shape
	thickness := rec.Width / 8  // Thickness of the lines
	length := rec.Width * 3 / 5 // Length of each arm

	// Horizontal bar
	hRec := rl.Rectangle{
		X:      centerX - length/2,
		Y:      centerY - thickness/2,
		Width:  length,
		Height: thickness,
	}

	rl.DrawRectangleRec(hRec, shapeColor)
}

func drawTimes(rec rl.Rectangle) {
	var (
		shrink     = float32(4.0)
		thickness  = float32(3.0)
		upperLeft  = rl.Vector2{X: rec.X + shrink, Y: rec.Y + shrink}
		lowerRight = rl.Vector2{X: rec.X + rec.Width - shrink, Y: rec.Y + rec.Height - shrink}
		lowerLeft  = rl.Vector2{X: rec.X + shrink, Y: rec.Y + rec.Height - shrink}
		upperRight = rl.Vector2{X: rec.X + rec.Width - shrink, Y: rec.Y + shrink}
	)

	rl.DrawLineEx(upperLeft, lowerRight, thickness, shapeColor)
	rl.DrawLineEx(lowerLeft, upperRight, thickness, shapeColor)
}

func drawRectangle(rect rl.Rectangle) {
	thickness := rect.Width / 8

	x := rect.X + thickness + 1
	y := rect.Y + thickness + 1

	width := rect.Width - 2*thickness - 2
	height := rect.Height - 2*thickness - 2

	rl.DrawRectangleLinesEx(rl.Rectangle{X: x, Y: y, Width: width, Height: height}, thickness, shapeColor)
}

func drawCircle(rect rl.Rectangle) {
	centerX := rect.X + rect.Width/2
	centerY := rect.Y + rect.Height/2
	radius := min(rect.Width, rect.Height) * 0.3

	rl.DrawRing(rl.Vector2{X: centerX, Y: centerY}, radius-1, radius+1.5, 0, 360, 24, shapeColor)
}

func drawScore(g *jb.Game, boardRec rl.Rectangle) {
	str := fmt.Sprintf("Score: %d  /  Last: %d  / Best: %d", g.Score(), g.LastScore(), g.BestScore())
	width := rl.MeasureText(str, scoreSize)
	rl.DrawText(str,
		int32(boardRec.Y+boardRec.Width/2.0-float32(width)/2.0),
		int32(boardRec.Y)-scoreSize-6,
		scoreSize,
		textColor,
	)
}

func drawGameOver(g *jb.Game, boardRec rl.Rectangle) {
	var text []string
	if g.Bonus() == 0 {
		text = []string{
			"Game Over!",
			fmt.Sprintf("You scored %d points.", g.Score()),
			"Press <SPACE> to restart",
		}
	} else {
		text = []string{
			"Game Over!",
			fmt.Sprintf("You scored %d points which", g.Score()),
			fmt.Sprintf("includes %d bonus points.", g.Bonus()),
			"Press <SPACE> to restart",
		}
	}

	pad := float32(34)
	rec := rl.Rectangle{
		X:      boardRec.X + pad,
		Y:      boardRec.Y + 2*pad,
		Width:  boardRec.Width - pad*2,
		Height: boardRec.Height - pad*4,
	}

	rl.DrawRectangleRec(rec, rl.NewColor(20, 24, 33, 220))

	gap := int32(pad)
	for i, t := range text {
		width := rl.MeasureText(t, overSize)
		x := int32(rec.X + (rec.Width-float32(width))/2)
		y := int32(rec.Y) + int32(i+1)*gap // tweak factor to taste for vertical placement
		rl.DrawText(t, x, y, overSize, rl.White)
	}
}
