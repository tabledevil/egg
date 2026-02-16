package canvas

import (
	"testing"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func TestCanvas_Render(t *testing.T) {
	c := New(3, 2)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	c.SetChar(0, 0, 'A', style)
	c.SetChar(1, 0, 'B', style)
	c.SetChar(2, 0, 'C', style)

	c.SetChar(0, 1, '1', style)
	c.SetChar(1, 1, '2', style)
	c.SetChar(2, 1, '3', style)

	output := c.Render()

	// Just check if output contains the characters and has 2 lines (plus possibly newline)
	if !strings.Contains(output, "A") || !strings.Contains(output, "1") {
		t.Errorf("Render output missing content: %q", output)
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Errorf("Render output expected at least 2 lines, got %d", len(lines))
	}
}
