package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	jb "github.com/sspencer/jawbreaker"
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

	selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))

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
	game := jb.NewGame(6, 6) // Create a 6x6 grid as specified
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
			if m.cursorY < 5 {
				m.cursorY++
			}

		case key.Matches(msg, m.keys.Left):
			if m.cursorX > 0 {
				m.cursorX--
			}

		case key.Matches(msg, m.keys.Right):
			if m.cursorX < 5 {
				m.cursorX++
			}

		case key.Matches(msg, m.keys.Select):
			index := m.cursorY*6 + m.cursorX
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
		// Handle mouse events
		if msg.Type == tea.MouseLeft {
			// Calculate grid position from mouse coordinates
			// Each cell is now 4 characters wide and 2 lines tall
			gridX := msg.X / 5
			gridY := (msg.Y - 3) / 2 // Adjust for title and spacing, and divide by 2 for double-height cells

			// Check if click is within grid bounds
			if gridX >= 0 && gridX < 6 && gridY >= 0 && gridY < 6 {
				// Update cursor position
				m.cursorX = gridX
				m.cursorY = gridY

				// Make a move at the clicked position
				index := m.cursorY*6 + m.cursorX
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
	m.jawbreaker = jb.NewGame(6, 6)
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
	s += titleStyle.Render("Jawbreaker TUI") + "\n\n"

	board := m.jawbreaker.Board()
	for y := 0; y < 6; y++ {
		for x := 0; x < 6; x++ {
			index := y*6 + x
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

			// Add selection border if this is the cursor position
			if x == m.cursorX && y == m.cursorY {
				cellStyle = selectedStyle.Copy().Inherit(cellStyle)
			}

			s += cellStyle.Render("  ") + cellStyle.Render("  ")
		}
		s += "\n"
		// Add another line with the same cells to make them more square-like
		for x := 0; x < 6; x++ {
			index := y*6 + x
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

			// Add selection border if this is the cursor position
			if x == m.cursorX && y == m.cursorY {
				cellStyle = selectedStyle.Copy().Inherit(cellStyle)
			}

			s += cellStyle.Render("  ") + cellStyle.Render("  ")
		}
		s += "\n"
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
