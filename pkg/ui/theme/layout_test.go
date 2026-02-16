package theme

import "strings"
import "testing"

func TestWrapTextWrapsLongToken(t *testing.T) {
	lines := wrapText("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 8)
	if len(lines) < 3 {
		t.Fatalf("expected wrapped long token into multiple lines")
	}
	for _, line := range lines {
		if runeLen(line) > 8 {
			t.Fatalf("wrapped line exceeded width: %q", line)
		}
	}
}

func TestWrapLabeledKeepsPrefixOnFirstLine(t *testing.T) {
	lines := wrapLabeled("HINT: ", "one two three four five", 12)
	if len(lines) < 2 {
		t.Fatalf("expected wrapped labeled lines")
	}
	if !strings.HasPrefix(lines[0], "HINT: ") {
		t.Fatalf("expected first line to include prefix, got %q", lines[0])
	}
	if runeLen(lines[1]) > 12 {
		t.Fatalf("continuation line exceeded width: %q", lines[1])
	}
}

func TestClampLinesAddsEllipsisWhenTruncated(t *testing.T) {
	lines := []string{"line one", "line two", "line three"}
	clamped := clampLines(lines, 2, 10)
	if len(clamped) != 2 {
		t.Fatalf("expected 2 lines after clamp, got %d", len(clamped))
	}
	if clamped[1] != "line two…" {
		t.Fatalf("expected ellipsis on final clamped line, got %q", clamped[1])
	}
}

func TestSliceRunes(t *testing.T) {
	got := sliceRunes("abcdef", 2, 3)
	if got != "cde" {
		t.Fatalf("unexpected rune slice: %q", got)
	}
}
