package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type MinimalTheme struct{}

func (t MinimalTheme) Name() string        { return "Root Shell" }
func (t MinimalTheme) Description() string { return "Clean, minimal root shell access" }

func (t MinimalTheme) Render(width, height int, q *game.Question, inputView string, hint string) string {
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
	Register(MinimalTheme{})
}
