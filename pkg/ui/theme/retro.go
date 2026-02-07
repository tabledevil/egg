package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"strings"
)

type RetroTheme struct{}

func NewRetroTheme() Theme {
	return &RetroTheme{}
}

func (t *RetroTheme) Init() tea.Cmd {
	return nil
}

func (t *RetroTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	return t, nil
}

func (t *RetroTheme) Name() string        { return "Retro CRT" }
func (t *RetroTheme) Description() string { return "Classic green-on-black terminal" }

func (t *RetroTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile <= termenv.ANSI
}

func (t *RetroTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	// Styles
	green := lipgloss.Color("#00FF00")
	black := lipgloss.Color("#000000")

	baseStyle := lipgloss.NewStyle().
		Foreground(green).
		Background(black).
		Width(width).
		Height(height)

	contentWidth := width - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(green).
		Padding(1).
		Width(contentWidth)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		MarginBottom(1)

	textStyle := lipgloss.NewStyle().
		MarginBottom(2).
		Width(contentWidth - 6)

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
	Register(NewRetroTheme)
}
