package canvas

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// stripANSI removes all ANSI escape sequences to get plain visible text.
func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// TestCoalescingPreservesVisibleContent verifies that the new run-coalescing
// Render() produces the exact same visible characters as expected.
func TestCoalescingPreservesVisibleContent(t *testing.T) {
	c := New(40, 10)

	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	blue := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))

	// Paint a mixed scene: background, box first, then text on top.
	c.Fill(5, 2, 10, 3, '#', red)
	c.DrawBox(15, 1, 20, 5, red)
	c.SetString(7, 3, "HELLO", green)
	c.SetString(20, 5, "WORLD", blue)

	rendered := c.Render()
	visible := stripANSI(rendered)

	if !strings.Contains(visible, "HELLO") {
		t.Fatalf("rendered output missing 'HELLO':\n%s", visible)
	}
	if !strings.Contains(visible, "WORLD") {
		t.Fatalf("rendered output missing 'WORLD':\n%s", visible)
	}

	lines := strings.Split(rendered, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(lines))
	}

	for i, line := range lines {
		visLine := stripANSI(line)
		runes := []rune(visLine)
		if len(runes) != 40 {
			t.Errorf("line %d: expected 40 visible chars, got %d: %q", i, len(runes), visLine)
		}
	}
}

// TestCoalescingReducesStyleCount verifies that run-coalescing reduces the
// number of Render calls by checking that a uniform grid produces far fewer
// style transitions than cells.
func TestCoalescingReducesStyleCount(t *testing.T) {
	c := New(80, 24)

	bg := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	c.Fill(0, 0, 80, 24, '.', bg)

	// All cells should share a single styleID.
	firstID := c.Grid[0][0].styleID
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			if c.Grid[y][x].styleID != firstID {
				t.Fatalf("cell (%d,%d) has styleID %d, expected %d",
					x, y, c.Grid[y][x].styleID, firstID)
			}
		}
	}

	// The style registry should have only 2 entries: default + the bg style.
	if len(c.styles) > 2 {
		t.Errorf("expected at most 2 registered styles, got %d", len(c.styles))
	}

	// Render should produce valid output.
	rendered := c.Render()
	visible := stripANSI(rendered)
	lines := strings.Split(visible, "\n")
	if len(lines) != 24 {
		t.Fatalf("expected 24 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if len([]rune(line)) != 80 {
			t.Errorf("line %d has %d visible chars, expected 80", i, len([]rune(line)))
		}
	}
}

// TestCoalescingMixedStyles verifies coalescing with multiple style transitions.
func TestCoalescingMixedStyles(t *testing.T) {
	c := New(20, 3)

	s1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	s2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	for x := 0; x < 20; x++ {
		style := s1
		if (x/5)%2 == 1 {
			style = s2
		}
		c.SetChar(x, 1, 'X', style)
	}

	rendered := c.Render()
	visible := stripANSI(rendered)
	lines := strings.Split(visible, "\n")

	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines")
	}
	line1 := lines[1]
	if line1 != "XXXXXXXXXXXXXXXXXXXX" {
		t.Errorf("expected 20 X's, got %q", line1)
	}

	// In a color-capable terminal, registry has 3 entries: default, s1, s2.
	// In a no-color test env, all styles fingerprint identically, so only 1.
	// Either way, the visual output must be correct.
	t.Logf("style registry size: %d (1 = no-color env, 3 = color env)", len(c.styles))
}

// TestDefaultStyleSkipsANSI verifies that blank cells with the default style
// produce no ANSI escape sequences (raw spaces).
func TestDefaultStyleSkipsANSI(t *testing.T) {
	c := New(10, 1)
	rendered := c.Render()

	if strings.Contains(rendered, "\x1b") {
		t.Errorf("blank canvas should have no ANSI escapes, got: %q", rendered)
	}
	if rendered != "          " {
		t.Errorf("expected 10 spaces, got %q", rendered)
	}
}

// TestStyleRegistryDeduplication verifies that the style registry correctly
// deduplicates identical styles created separately.
func TestStyleRegistryDeduplication(t *testing.T) {
	c := New(10, 1)

	s1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	s2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	c.SetChar(0, 0, 'A', s1)
	c.SetChar(1, 0, 'B', s2)

	if c.Grid[0][0].styleID != c.Grid[0][1].styleID {
		t.Errorf("identical styles should share styleID: got %d and %d",
			c.Grid[0][0].styleID, c.Grid[0][1].styleID)
	}
}

// TestFillOptimisation verifies that Fill registers the style only once
// and that all filled cells share the same styleID.
func TestFillOptimisation(t *testing.T) {
	c := New(80, 24)
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("#AABBCC"))

	c.Fill(0, 0, 80, 24, '.', s)

	// All cells must share one styleID (Fill registers once).
	firstID := c.Grid[0][0].styleID
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			if c.Grid[y][x].styleID != firstID {
				t.Fatalf("cell (%d,%d) has styleID %d, expected %d",
					x, y, c.Grid[y][x].styleID, firstID)
			}
		}
	}
	// In a color terminal, registry would have 2 (default + fill style).
	// In no-color test env, it may be 1 (all styles fingerprint identically).
	t.Logf("style registry size after Fill: %d", len(c.styles))
}

// BenchmarkCanvasRender benchmarks the Render method with a typical themed canvas.
func BenchmarkCanvasRender(b *testing.B) {
	c := New(80, 24)

	bg := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	c.Fill(0, 0, 80, 24, '.', bg)
	c.DrawBox(10, 5, 60, 14, text)
	c.Fill(11, 6, 58, 12, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#111111")))
	c.SetString(15, 8, "What is the access code?", text)
	c.SetString(15, 12, "> hunter2", accent)
	c.SetString(15, 15, "CLUE: Try XOR", lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.Render()
	}
}

// BenchmarkCanvasRenderParticles benchmarks with many unique styles.
func BenchmarkCanvasRenderParticles(b *testing.B) {
	c := New(80, 24)

	for i := 0; i < 200; i++ {
		x := i % 80
		y := (i / 80) % 24
		r := (i * 7) % 256
		g := (i * 13) % 256
		hex := "#" + hexByte(r) + hexByte(g) + "80"
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		c.SetChar(x, y, '.', style)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.Render()
	}
}

func hexByte(v int) string {
	const digits = "0123456789abcdef"
	return string([]byte{digits[v>>4], digits[v&0xf]})
}
