package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type RetroTheme struct{}

func (t RetroTheme) Name() string        { return "Retro CRT" }
func (t RetroTheme) Description() string { return "Classic green-on-black terminal" }

func (t RetroTheme) Render(width, height int, q *game.Question, inputView string, hint string) string {
	// Styles
	green := lipgloss.Color("#00FF00")
	black := lipgloss.Color("#000000")

	baseStyle := lipgloss.NewStyle().
		Foreground(green).
		Background(black).
		Width(width).
		Height(height)

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(green).
		Padding(1).
		Width(width - 4) // Adjust for border

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		MarginBottom(1)

	textStyle := lipgloss.NewStyle().
		MarginBottom(2).
		Width(width - 10)

	// Content Construction
	var content strings.Builder

	content.WriteString(titleStyle.Render(fmt.Sprintf("ACCESS NODE %02d", q.ID)))
	content.WriteString("\n")
	content.WriteString(textStyle.Render(q.Text))
	content.WriteString("\n\n")

	content.WriteString("> ")
	content.WriteString(inputView)

	if hint != "" {
		content.WriteString("\n\n")
		content.WriteString(lipgloss.NewStyle().Faint(true).Render("HINT: " + hint))
	}

	// Wrap in border and base
	return baseStyle.Render(
		borderStyle.Render(content.String()),
	)
}

func init() {
	Register(RetroTheme{})
}
