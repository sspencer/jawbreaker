package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (g *Game) drawBlocks() {
	for y := range g.board {
		for x := range g.board[y] {
			rec := rl.Rectangle{
				X:      float32(x*BlockSize + WindowPadding + GridOuter),
				Y:      float32(y*BlockSize + WindowPadding + GridOuter),
				Width:  BlockSize - GridInner,
				Height: BlockSize - GridInner,
			}

			p := g.board[y][x]
			lo := blockColorValues[p]
			hi := highlightColorValues[p]

			if p == Empty {
				rl.DrawRectangleRec(rec, hi)
			} else {
				rl.DrawRectangleGradientEx(rec, hi, lo, lo, lo)
			}

			g.drawShape(p, rec)
		}
	}
}

func (g *Game) drawShape(p GameBlock, rec rl.Rectangle) {
	switch p {
	case PowerPlus:
		g.drawPlus(rec)
	case PowerPipe:
		g.drawPipe(rec)
	case PowerMinus:
		g.drawMinus(rec)
	case PowerTimes:
		g.drawTimes(rec)
	case PowerRect:
		g.drawRectangle(rec)
	case PowerCircle:
		g.drawCircle(rec)
	default:
		// ignore
	}
}

func (g *Game) drawSelection(sel []Vec2i) {
	if len(sel) < 2 {
		return
	}

	var (
		color     = rl.NewColor(255, 255, 255, 128)
		xfg       = rl.White
		xbg       = rl.Red
		shrink    = float32(8.0) // shrink X shape from corners
		thickness = float32(2.0)
		radius    = float32(BlockSize-GridInner) * 0.35
	)

	for _, s := range sel {
		rec := rl.Rectangle{
			X:      float32(s.X*BlockSize + WindowPadding + GridOuter),
			Y:      float32(s.Y*BlockSize + WindowPadding + GridOuter),
			Width:  BlockSize - GridInner,
			Height: BlockSize - GridInner,
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

func (g *Game) drawPlus(rec rl.Rectangle) {
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

	rl.DrawRectangleRec(vRec, rl.White)
	rl.DrawRectangleRec(hRec, rl.White)
}

func (g *Game) drawPipe(rec rl.Rectangle) {
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

	rl.DrawRectangleRec(vRec, rl.White)
}

func (g *Game) drawMinus(rec rl.Rectangle) {
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

	rl.DrawRectangleRec(hRec, rl.White)
}

func (g *Game) drawTimes(rec rl.Rectangle) {
	var (
		f          float32 = 4.0
		thickness  float32 = 3.0
		color              = rl.White
		upperLeft          = rl.Vector2{X: rec.X + f, Y: rec.Y + f}
		lowerRight         = rl.Vector2{X: rec.X + rec.Width - f, Y: rec.Y + rec.Height - f}
		lowerLeft          = rl.Vector2{X: rec.X + f, Y: rec.Y + rec.Height - f}
		upperRight         = rl.Vector2{X: rec.X + rec.Width - f, Y: rec.Y + f}
	)

	rl.DrawLineEx(upperLeft, lowerRight, thickness, color)
	rl.DrawLineEx(lowerLeft, upperRight, thickness, color)

}

func (g *Game) drawRectangle(rect rl.Rectangle) {
	thickness := rect.Width / 8

	x := rect.X + thickness + 1
	y := rect.Y + thickness + 1

	width := rect.Width - 2*thickness - 2
	height := rect.Height - 2*thickness - 2

	rl.DrawRectangleLinesEx(rl.Rectangle{x, y, width, height}, thickness, rl.White)
}

func (g *Game) drawCircle(rect rl.Rectangle) {
	centerX := rect.X + rect.Width/2
	centerY := rect.Y + rect.Height/2
	radius := min(rect.Width, rect.Height) * 0.3

	rl.DrawRing(rl.Vector2{centerX, centerY}, radius-1, radius+1.5, 0, 360, 24, rl.White)
}
