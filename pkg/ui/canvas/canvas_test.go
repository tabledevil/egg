package canvas

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestSetStringHandlesNewlines(t *testing.T) {
	c := New(6, 4)
	c.SetString(1, 0, "AB\nCD", lipgloss.NewStyle())

	if got := c.Grid[0][1].Rune; got != 'A' {
		t.Fatalf("expected A at row 0 col 1, got %q", got)
	}
	if got := c.Grid[0][2].Rune; got != 'B' {
		t.Fatalf("expected B at row 0 col 2, got %q", got)
	}
	if got := c.Grid[1][1].Rune; got != 'C' {
		t.Fatalf("expected C at row 1 col 1, got %q", got)
	}
	if got := c.Grid[1][2].Rune; got != 'D' {
		t.Fatalf("expected D at row 1 col 2, got %q", got)
	}
}

func TestSetStringPreservesEmptyLinesAndTrailingNewline(t *testing.T) {
	c := New(5, 4)
	c.SetString(0, 0, "A\n\nB\n", lipgloss.NewStyle())

	if got := c.Grid[0][0].Rune; got != 'A' {
		t.Fatalf("expected A at row 0 col 0, got %q", got)
	}
	if got := c.Grid[1][0].Rune; got != ' ' {
		t.Fatalf("expected blank row for empty line, got %q", got)
	}
	if got := c.Grid[2][0].Rune; got != 'B' {
		t.Fatalf("expected B at row 2 col 0, got %q", got)
	}
	if got := c.Grid[3][0].Rune; got != ' ' {
		t.Fatalf("expected trailing newline to leave final row unchanged, got %q", got)
	}
}

func TestSetStringClipsAcrossNewline(t *testing.T) {
	c := New(3, 3)
	c.SetString(-1, 0, "ABCD\nEFG", lipgloss.NewStyle())

	if got := c.Grid[0][0].Rune; got != 'B' {
		t.Fatalf("expected clipped B at row 0 col 0, got %q", got)
	}
	if got := c.Grid[0][1].Rune; got != 'C' {
		t.Fatalf("expected clipped C at row 0 col 1, got %q", got)
	}
	if got := c.Grid[0][2].Rune; got != 'D' {
		t.Fatalf("expected clipped D at row 0 col 2, got %q", got)
	}
	if got := c.Grid[1][0].Rune; got != 'F' {
		t.Fatalf("expected clipped F at row 1 col 0, got %q", got)
	}
	if got := c.Grid[1][1].Rune; got != 'G' {
		t.Fatalf("expected clipped G at row 1 col 1, got %q", got)
	}
}
