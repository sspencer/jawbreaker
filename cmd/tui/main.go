package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	jb "github.com/sspencer/jawbreaker"
)

// Game constants
const (
	gridSize    = 10 // Size of the game grid (gridSize x gridSize)
	blockWidth  = 4  // Width of each cell in characters
	blockHeight = 2  // Height of each cell in lines
)

// Define key mappings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Reset  key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Select, k.Reset, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Reset, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset game"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

// Define color styles
var (
	purpleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#8a2be2")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	blueStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0070ff")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	greenStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#00cc66")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	redStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ff0055")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	yellowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ff8800")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)

	emptyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#333333")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true).
			Padding(0, 1)

	scoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00")).
			Bold(true).
			Padding(0, 1)
)

// Model represents the game state
type Model struct {
	jawbreaker   *jb.Game
	cursorX      int
	cursorY      int
	currentScore int
	lastScore    int
	bestScore    int
	help         help.Model
	keys         keyMap
	width        int
	height       int
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// NewModel creates a new model
func NewModel() Model {
	game := jb.NewGame(gridSize, gridSize) // Create a grid as specified
	return Model{
		jawbreaker:   game,
		cursorX:      0,
		cursorY:      0,
		currentScore: 0,
		lastScore:    0,
		bestScore:    0,
		help:         help.New(),
		keys:         keys,
		width:        0,
		height:       0,
	}
}

// Update handles user input and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("UI update triggered")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursorY > 0 {
				m.cursorY--
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursorY < gridSize-1 {
				m.cursorY++
			}

		case key.Matches(msg, m.keys.Left):
			if m.cursorX > 0 {
				m.cursorX--
			}

		case key.Matches(msg, m.keys.Right):
			if m.cursorX < gridSize-1 {
				m.cursorX++
			}

		case key.Matches(msg, m.keys.Select):
			index := m.cursorY*gridSize + m.cursorX
			status := m.jawbreaker.Move(index)

			// Update scores
			m.currentScore = status.Score
			if status.GameOver {
				if m.currentScore > m.bestScore {
					m.bestScore = m.currentScore
				}
				m.lastScore = m.currentScore
				m.resetGame()
			}

		case key.Matches(msg, m.keys.Reset):
			m.resetGame()
		}

	case tea.MouseMsg:
		// Calculate grid position from mouse coordinates
		gridX := msg.X / (blockWidth + 2)
		gridY := (msg.Y - 2) / blockHeight

		// Check if position is within grid bounds
		if gridX >= 0 && gridX < gridSize && gridY >= 0 && gridY < gridSize {
			// Update cursor position for both movement and clicks
			m.cursorX = gridX
			m.cursorY = gridY

			// Handle mouse clicks
			if msg.Action == tea.MouseActionPress {
				// Make a move at the clicked position
				index := m.cursorY*gridSize + m.cursorX
				log.Printf("Mouse click at index %d: X=%d, Y=%d (%d, %d)\n", index, gridX, gridY, msg.X, msg.Y)

				status := m.jawbreaker.Move(index)

				// Update scores
				m.currentScore = status.Score
				if status.GameOver {
					if m.currentScore > m.bestScore {
						m.bestScore = m.currentScore
					}
					m.lastScore = m.currentScore
					m.resetGame()
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
	}

	return m, nil
}

// resetGame resets the game state
func (m *Model) resetGame() {
	m.jawbreaker = jb.NewGame(gridSize, gridSize)
	m.currentScore = 0
	m.cursorX = 0
	m.cursorY = 0
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Render the grid
	var s string

	// Calculate the width of the grid in characters
	gridWidthInChars := gridSize * blockWidth

	// Center the title over the game board
	title := titleStyle.Render("Jawbreaker TUI")
	titlePadding := (gridWidthInChars - lipgloss.Width(title)) / 2
	if titlePadding < 0 {
		titlePadding = 0
	}

	s += strings.Repeat(" ", titlePadding) + title + "\n\n"

	board := m.jawbreaker.Board()
	for y := 0; y < gridSize; y++ {
		for i := 0; i < blockHeight; i++ {
			for x := 0; x < gridSize; x++ {
				index := y*gridSize + x
				piece := board[index]

				// Choose style based on piece color
				var cellStyle lipgloss.Style
				switch piece {
				case jb.Purple:
					cellStyle = purpleStyle
				case jb.Blue:
					cellStyle = blueStyle
				case jb.Green:
					cellStyle = greenStyle
				case jb.Red:
					cellStyle = redStyle
				case jb.Yellow:
					cellStyle = yellowStyle
				default:
					cellStyle = emptyStyle
				}

				// Add pointer-like highlight if this is the cursor position
				block := " "
				if x == m.cursorX && y == m.cursorY {
					block = "-"
					cellStyle = cellStyle.Foreground(lipgloss.Color("#ffffff")).Bold(true)
				}

				// Render a cell with blockWidth characters
				s += cellStyle.Render(strings.Repeat(block, blockWidth))
			}
			s += "\n"
		}
	}

	// Render scores
	s += "\n"
	scorePanel := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Current"),
			scoreStyle.Render(fmt.Sprintf("%d", m.currentScore)),
		),
		"   ",
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Last"),
			scoreStyle.Render(fmt.Sprintf("%d", m.lastScore)),
		),
		"   ",
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Best"),
			scoreStyle.Render(fmt.Sprintf("%d", m.bestScore)),
		),
	)

	s += scorePanel + "\n\n"
	s += m.help.View(m.keys)

	return s
}

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("jawbreaker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
