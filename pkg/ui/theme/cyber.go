package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type CyberTheme struct{}

func NewCyberTheme() Theme {
	return &CyberTheme{}
}

func (t *CyberTheme) Init() tea.Cmd {
	return nil
}

func (t *CyberTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	return t, nil
}

func (t *CyberTheme) Name() string        { return "Neon City" }
func (t *CyberTheme) Description() string { return "High contrast neon colors" }

func (t *CyberTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	pink := lipgloss.Color("#FF00FF")
	cyan := lipgloss.Color("#00FFFF")
	darkBlue := lipgloss.Color("#000033")

	baseStyle := lipgloss.NewStyle().
		Background(darkBlue).
		Width(width).
		Height(height)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(pink).
		BorderTop(true).
		BorderBottom(true).
		Padding(1, 2).
		Width(width).
		Align(lipgloss.Center)

	questionStyle := lipgloss.NewStyle().
		Foreground(cyan).
		Bold(true).
		Background(lipgloss.Color("#111111")).
		Padding(1).
		Width(width - 10).
		Align(lipgloss.Center)

	inputStyle := lipgloss.NewStyle().
		Foreground(pink).
		MarginTop(2)

	// Content
	var content strings.Builder

	header := lipgloss.NewStyle().
		Foreground(darkBlue).
		Background(pink).
		Padding(0, 1).
		Bold(true).
		Render(fmt.Sprintf(" MISSION #%d ", q.ID))

	content.WriteString(header + "\n\n")
	content.WriteString(questionStyle.Render(q.Text))
	content.WriteString("\n")
	content.WriteString(inputStyle.Render(inputView))

	if hint != "" {
		content.WriteString("\n\n")
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("⚠ " + hint))
	}

	return baseStyle.Render(
		lipgloss.PlaceVertical(height, lipgloss.Center,
			boxStyle.Render(content.String()),
		),
	)
}

func init() {
	Register(NewCyberTheme)
}
