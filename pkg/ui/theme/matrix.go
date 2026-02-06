package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

type MatrixTheme struct {
	tickCount int
}

func NewMatrixTheme() Theme {
	return &MatrixTheme{}
}

func (t *MatrixTheme) Init() tea.Cmd {
	return tick()
}

func (t *MatrixTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.tickCount++
		return t, tick()
	}
	return t, nil
}

func (t *MatrixTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	boxWidth := 50
	if boxWidth > width-4 { boxWidth = width-4 }

	// Animate Border Color
	color := "#00FF00"
	if t.tickCount % 6 < 3 {
		color = "#005500"
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(color)).
		Background(lipgloss.Color("#000000")).
		Padding(1).
		Width(boxWidth).
		Align(lipgloss.Center)

	content := fmt.Sprintf("WAKE UP NEO...\n\n%s\n\n> %s", q.Text, inputView)
	if hint != "" {
		content += "\n\nHINT: " + hint
	}

	renderedBox := boxStyle.Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, renderedBox)
}

func (t *MatrixTheme) Name() string { return "The Matrix" }
func (t *MatrixTheme) Description() string { return "Digital Rain" }

func init() {
	Register(NewMatrixTheme)
}
