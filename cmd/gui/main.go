package main

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	jb "github.com/sspencer/jawbreaker"
	"github.com/sspencer/jawbreaker/util"
)

// Game constants
const (
	scoreHeight = 100
	blockSize   = 50
	borderSize  = 3 // Size of the black border between blocks
	gridSize    = 12
)

// Game configuration (can be modified by command line args)
var (
	screenWidth  int
	screenHeight int
)

// Colors for the blocks - start and end colors for gradients
var (
	// Purple: linear-gradient(135deg, #8a2be2, #d442ff)
	colorPurpleStart = color.RGBA{R: 212, G: 66, B: 255, A: 255} // #d442ff
	colorPurpleEnd   = color.RGBA{R: 138, G: 43, B: 226, A: 255} // #8a2be2

	// Blue: linear-gradient(135deg, #00f0ff, #0070ff)
	colorBlueStart = color.RGBA{R: 0, G: 240, B: 255, A: 255} // #00f0ff
	colorBlueEnd   = color.RGBA{R: 0, G: 72, B: 255, A: 255}  // #0070ff

	// Green: linear-gradient(135deg, #00ff99, #00cc66)
	colorGreenStart = color.RGBA{R: 51, G: 255, B: 153, A: 255} // #00ff99
	colorGreenEnd   = color.RGBA{R: 0, G: 164, B: 72, A: 255}   // #00cc66

	// Red: linear-gradient(135deg, #ff5500, #ff0055)
	colorRedStart = color.RGBA{R: 255, G: 105, B: 30, A: 255} // #ff5500
	colorRedEnd   = color.RGBA{R: 255, G: 0, B: 0, A: 255}    // #ff0055

	// Yellow: linear-gradient(135deg, #ffee00, #ff8800)
	colorYellowStart = color.RGBA{R: 255, G: 238, B: 0, A: 255} // #ffee00
	colorYellowEnd   = color.RGBA{R: 215, G: 106, B: 0, A: 255} // #ff8800

	colorBlackStart = color.RGBA{R: 44, G: 44, B: 44, A: 255} // Nearly Black
	colorBlackEnd   = color.RGBA{R: 22, G: 22, B: 22, A: 255} // Dark background for contrast
)

// Game implements ebiten.Game interface
type Game struct {
	jawbreaker   *jb.Game
	currentScore int
	lastScore    int
	bestScore    int
	connections  []jb.Vec2i
	images       map[jb.Tile]*ebiten.Image // Map to store texture images for each color and state
	smallFont    font.Face
	normalFont   font.Face
	titleFont    font.Face
	valueFont    font.Face
}

// loadFonts loads and initializes font faces of different sizes and weights
func loadFonts() (font.Face, font.Face, font.Face, font.Face, error) {
	// Load OpenType fonts
	regularTT, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	boldTT, err := opentype.Parse(gobold.TTF)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Create different font sizes and weights
	smallFont, err := opentype.NewFace(regularTT, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	normalFont, err := opentype.NewFace(regularTT, &opentype.FaceOptions{
		Size: 16,
		DPI:  72,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	titleFont, err := opentype.NewFace(boldTT, &opentype.FaceOptions{
		Size: 18,
		DPI:  72,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	valueFont, err := opentype.NewFace(boldTT, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return smallFont, normalFont, titleFont, valueFont, nil
}

// NewGame creates a new game
func NewGame() *Game {
	// Try to load custom fonts
	smallFont, normalFont, titleFont, valueFont, err := loadFonts()
	if err != nil {
		// Fallback to basic fonts if loading fails
		basicFace := basicfont.Face7x13
		smallFont = basicFace
		normalFont = basicFace
		titleFont = basicFace
		valueFont = basicFace
	}

	lastScore, bestScore, _ := util.LoadScores()
	game := &Game{
		jawbreaker: jb.New(jb.WithSize(gridSize), jb.WithLastScore(lastScore), jb.WithBestScore(bestScore), jb.WithoutExtras()),
		smallFont:  smallFont,
		normalFont: normalFont,
		titleFont:  titleFont,
		valueFont:  valueFont,
		images:     make(map[jb.Tile]*ebiten.Image),
	}

	// Create all the texture images we'll need
	game.generateTextures()

	return game
}

// generateTextures creates all necessary texture images
func (g *Game) generateTextures() {
	// Generate stained glass textures
	g.images[jb.TilePurple] = createStainedGlassTexture(colorPurpleStart, colorPurpleEnd, false)
	g.images[jb.TileBlue] = createStainedGlassTexture(colorBlueStart, colorBlueEnd, false)
	g.images[jb.TileGreen] = createStainedGlassTexture(colorGreenStart, colorGreenEnd, false)
	g.images[jb.TileRed] = createStainedGlassTexture(colorRedStart, colorRedEnd, false)
	g.images[jb.TileYellow] = createStainedGlassTexture(colorYellowStart, colorYellowEnd, false)
	g.images[jb.TileEmpty] = createStainedGlassTexture(colorBlackEnd, colorBlackStart, false)
}

// createStainedGlassTexture creates a modern gradient texture for a block
func createStainedGlassTexture(startColor, endColor color.RGBA, _ bool) *ebiten.Image {
	blockRect := image.Rect(0, 0, blockSize-borderSize*2, blockSize-borderSize*2)
	img := ebiten.NewImage(blockRect.Dx(), blockRect.Dy())

	opacity := uint8(230)

	width := blockRect.Dx()
	height := blockRect.Dy()

	// Create a 135-degree linear gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate position along the 135-degree diagonal
			// This maps (0,0) to 0.0 and (width,height) to 1.0
			progress := (float64(x) + float64(y)) / float64(width+height)

			// Apply ripple-like brightness to simulate glass distortion
			brightness := 0.9 + 0.1*math.Sin(float64(x*y%7))

			// Interpolate between start and end colors with brightness
			r := clamp(uint8(float64(startColor.R)*(1.0-progress) + float64(endColor.R)*progress*brightness))
			g := clamp(uint8(float64(startColor.G)*(1.0-progress) + float64(endColor.G)*progress*brightness))
			b := clamp(uint8(float64(startColor.B)*(1.0-progress) + float64(endColor.B)*progress*brightness))

			// Apply opacity
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: opacity})
		}
	}

	// Note: The box-shadow effect will be applied when drawing the block
	// since we can't directly add shadows to the image itself

	return img
}

// clamp ensures a value stays within 0-255 range
func clamp(value uint8) uint8 {
	if value > 255 {
		return 255
	}
	return value
}

// handleInput processes user input
func (g *Game) handleInput() error {
	// Get current mouse position for hover effect
	x, y := ebiten.CursorPosition()

	// This approach is less efficient (O(n) time) but may be preferred if:
	// You want to keep the exact same map allocation (though in practice this is rarely needed)
	// The map is shared and you don't want to change the reference
	//for key := range myMap {
	//	delete(myMap, key)
	//}

	// Check if mouse is within grid
	if y < screenHeight-scoreHeight {
		// Adjust for the double border (2*borderSize on each side)
		adjustedX := x - 2*borderSize
		adjustedY := y - 2*borderSize

		// Calculate grid coordinates
		gridX := adjustedX / blockSize
		gridY := adjustedY / blockSize

		g.connections = g.jawbreaker.Connections(jb.Vec2i{X: gridX, Y: gridY})
	}

	// Check for mouse click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.handleMouseClick(x, y)
	}

	// Reset the game if R key is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetGame()
	}

	// Undo the last move if U key is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		if g.jawbreaker.CanUndo() {
			g.jawbreaker.Undo()
			g.currentScore = g.jawbreaker.Score()
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyQ) || inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		_ = util.SaveScores(g.currentScore, g.bestScore)
		return ebiten.Termination
	}

	return nil
}

// handleMouseClick processes mouse clicks
func (g *Game) handleMouseClick(x, y int) {
	// Check if click is within grid
	if y < screenHeight-scoreHeight {
		// Adjust for the double border (2*borderSize on each side)
		adjustedX := x - 2*borderSize
		adjustedY := y - 2*borderSize

		// Calculate grid coordinates
		gridX := adjustedX / blockSize
		gridY := adjustedY / blockSize

		// Ensure within bounds
		if gridX >= 0 && gridX < gridSize && gridY >= 0 && gridY < gridSize {
			g.jawbreaker.Move(g.connections)
			g.currentScore = g.jawbreaker.Score()

			if g.jawbreaker.GameOver() {
				g.lastScore = g.jawbreaker.LastScore()
				g.bestScore = g.jawbreaker.BestScore()
				g.resetGame()
			}
		}
	}
}

// resetGame resets the game state
func (g *Game) resetGame() {
	g.jawbreaker.Restart()
	g.currentScore = 0
	_ = util.SaveScores(g.jawbreaker.LastScore(), g.jawbreaker.BestScore())
}

// Update updates the game state
func (g *Game) Update() error {
	return g.handleInput()
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw dark background
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(screenWidth),
		float32(screenHeight),
		colorBlackEnd,
		true,
	)

	// Draw grid
	g.drawGrid(screen)

	// Draw score panel
	g.drawScorePanel(screen)
}

// drawGrid draws the game grid
func (g *Game) drawGrid(screen *ebiten.Image) {
	hover := make(map[jb.Vec2i]bool)

	for _, c := range g.connections {
		hover[c] = true
	}

	board := g.jawbreaker.Board()
	for y := range board {
		for x := range board[y] {
			posX := float64(x*blockSize + 2*borderSize)
			posY := float64(y*blockSize + 2*borderSize)

			hovered := hover[jb.Vec2i{X: x, Y: y}]
			g.drawBlock(screen, board[y][x], posX, posY, hovered)
		}
	}
}

// drawBlock draws a single block
func (g *Game) drawBlock(screen *ebiten.Image, piece jb.Tile, posX, posY float64, isHovered bool) {

	// Draw the block with modern gradient texture
	op := &ebiten.DrawImageOptions{}

	// Apply final position
	op.GeoM.Translate(posX, posY)

	// If hovered, draw a border of the same color as the piece
	if isHovered {
		g.drawHighlightedBorder(screen, piece, posX, posY)
	}

	// Draw the appropriate texture
	screen.DrawImage(g.images[piece], op)
}

// drawHighlightedBorder draws a border of the same color as the piece around highlighted pieces
func (g *Game) drawHighlightedBorder(screen *ebiten.Image, piece jb.Tile, posX, posY float64) {
	var borderColor color.RGBA
	switch piece {
	case jb.TileRed:
		borderColor = colorRedStart
	case jb.TileBlue:
		borderColor = colorBlueStart
	case jb.TileYellow:
		borderColor = colorYellowStart
	case jb.TileGreen:
		borderColor = colorGreenStart
	case jb.TilePurple:
		borderColor = colorPurpleStart
	default:
		return
	}

	width := float32(blockSize - borderSize*2)
	height := float32(blockSize - borderSize*2)
	bs := float32(borderSize)

	// Top border
	vector.DrawFilledRect(screen, float32(posX)-bs, float32(posY)-bs, width+2*bs, bs, borderColor, true)
	// Bottom border
	vector.DrawFilledRect(screen, float32(posX)-bs, float32(posY)+height, width+2*bs, bs, borderColor, true)
	// Left border
	vector.DrawFilledRect(screen, float32(posX)-bs, float32(posY), bs, height, borderColor, true)
	// Right border
	vector.DrawFilledRect(screen, float32(posX)+width, float32(posY), bs, height, borderColor, true)
}

// drawScorePanel draws the score panel at the bottom of the screen
func (g *Game) drawScorePanel(screen *ebiten.Image) {
	// Draw score panel background
	vector.DrawFilledRect(
		screen,
		0,
		float32(screenHeight-scoreHeight),
		float32(screenWidth),
		scoreHeight,
		color.RGBA{R: 15, G: 15, B: 25, A: 255}, // Slightly darker than main background
		true,
	)

	// Draw score panel dividers
	panelWidth := screenWidth / 3
	vector.StrokeLine(
		screen,
		float32(panelWidth),
		float32(screenHeight-scoreHeight),
		float32(panelWidth),
		float32(screenHeight),
		1,
		color.RGBA{R: 100, G: 100, B: 100, A: 255},
		true,
	)
	vector.StrokeLine(
		screen,
		float32(panelWidth*2),
		float32(screenHeight-scoreHeight),
		float32(panelWidth*2),
		float32(screenHeight),
		1,
		color.RGBA{R: 100, G: 100, B: 100, A: 255},
		true,
	)

	// Calculate positions for centered text
	currentX := panelWidth / 2
	lastX := panelWidth + panelWidth/2
	bestX := 2*panelWidth + panelWidth/2

	// Draw scores with new fonts
	// Current score
	textWidth := measureTextWidth(g.titleFont, "Current")
	titleOffset := 20
	scoreOffset := 55
	text.Draw(screen, "Current", g.titleFont, currentX-textWidth/2, screenHeight-scoreHeight+titleOffset, color.White)

	scoreText := strconv.Itoa(g.currentScore)
	textWidth = measureTextWidth(g.valueFont, scoreText)
	text.Draw(screen, scoreText, g.valueFont, currentX-textWidth/2, screenHeight-scoreHeight+scoreOffset, color.White)

	// Last score
	textWidth = measureTextWidth(g.titleFont, "Last")
	text.Draw(screen, "Last", g.titleFont, lastX-textWidth/2, screenHeight-scoreHeight+titleOffset, color.White)

	scoreText = strconv.Itoa(g.lastScore)
	textWidth = measureTextWidth(g.valueFont, scoreText)
	text.Draw(screen, scoreText, g.valueFont, lastX-textWidth/2, screenHeight-scoreHeight+scoreOffset, color.White)

	// Best score
	textWidth = measureTextWidth(g.titleFont, "Best")
	text.Draw(screen, "Best", g.titleFont, bestX-textWidth/2, screenHeight-scoreHeight+titleOffset, color.White)

	scoreText = strconv.Itoa(g.bestScore)
	textWidth = measureTextWidth(g.valueFont, scoreText)
	text.Draw(screen, scoreText, g.valueFont, bestX-textWidth/2, screenHeight-scoreHeight+scoreOffset, color.White)

	// Draw instructions
	instructionsText := "Press 'R' to reset"
	textWidth = measureTextWidth(g.normalFont, instructionsText)
	text.Draw(screen, instructionsText, g.normalFont, screenWidth/2-textWidth/2, screenHeight-10, color.White)
}

// measureTextWidth calculates the estimated rendering width of a string
// using the provided font face. It sums the advance widths of all runes
// and includes kerning adjustments between adjacent runes.
//
// Parameters:
//   - face: The font.Face to use for measurement.
//   - text: The string to measure.
//
// Returns:
//   - The estimated width as a int value.
//     Returns 0 if the text is empty.
//
// Note: If a rune in the text does not have a corresponding glyph in the
// font face, its advance width is treated as 0 for this estimation.
func measureTextWidth(face font.Face, text string) int {
	if text == "" {
		return fixed.Int26_6(0).Ceil()
	}

	var totalWidth fixed.Int26_6
	var prevRune rune = -1 // Initialize with a value that won't match any rune

	for _, currentRune := range text {
		// Get advance width for the current rune
		advance, ok := face.GlyphAdvance(currentRune)
		if !ok {
			// Handle glyph not found - for estimation, we can add 0 width.
			// Alternatively, could add width of a replacement char like ' ' or '?'.
			advance = fixed.Int26_6(0)
		}
		totalWidth += advance

		// Apply kerning if it's not the first rune
		if prevRune >= 0 {
			totalWidth += face.Kern(prevRune, currentRune)
		}

		// Update previous rune
		prevRune = currentRune
	}

	return totalWidth.Ceil()
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Calculate screen dimensions with double border size
	screenWidth = gridSize*blockSize + 4*borderSize                // Double the border size (2*borderSize*2)
	screenHeight = gridSize*blockSize + 4*borderSize + scoreHeight // Double the border size (2*borderSize*2)

	// Set up the game
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Grid Puzzle")

	// Prevent screen from being cleared to white every frame
	// This helps prevent the white flash when the game starts
	ebiten.SetScreenClearedEveryFrame(false)

	// Create and run the game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
