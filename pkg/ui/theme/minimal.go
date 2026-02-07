package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type MinimalTheme struct{}

func NewMinimalTheme() Theme {
	return &MinimalTheme{}
}

func (t *MinimalTheme) Init() tea.Cmd {
	return nil
}

func (t *MinimalTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	return t, nil
}

func (t *MinimalTheme) Name() string        { return "Root Shell" }
func (t *MinimalTheme) Description() string { return "Clean, minimal root shell access" }

func (t *MinimalTheme) IsCompatible(c caps.Capabilities) bool {
	return true
}

func (t *MinimalTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	gray := lipgloss.Color("#888888")
	white := lipgloss.Color("#FFFFFF")

	baseStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(2)

	header := lipgloss.NewStyle().
		Foreground(gray).
		Render(fmt.Sprintf("root@system:~/challenges/0x%02x# ./view_question", q.ID))

	qText := lipgloss.NewStyle().
		Foreground(white).
		Bold(true).
		MarginLeft(2).
		Width(width - 10).
		Render(q.Text)

	hintText := ""
	if hint != "" {
		hintText = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Italic(true).Render("// " + hint)
	}

	var content strings.Builder
	content.WriteString(header + "\n")
	content.WriteString(strings.Repeat("-", width-10) + "\n\n")
	content.WriteString(qText + "\n\n")
	if hintText != "" {
		content.WriteString(hintText + "\n\n")
	}
	content.WriteString("root@system:~/input# " + inputView)

	return baseStyle.Render(content.String())
}

func init() {
	Register(NewMinimalTheme)
}
