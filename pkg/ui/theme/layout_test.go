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

func TestBoundedSpan(t *testing.T) {
	if got := boundedSpan(80, 4, 20, 60); got != 60 {
		t.Fatalf("expected span clamped to max 60, got %d", got)
	}
	if got := boundedSpan(20, 6, 14, 0); got != 14 {
		t.Fatalf("expected span to honor min when margins are tight, got %d", got)
	}
	if got := boundedSpan(5, 8, 10, 0); got != 5 {
		t.Fatalf("expected span to never exceed total width, got %d", got)
	}
}

func TestCenteredBoxAndInset(t *testing.T) {
	box := centeredBox(80, 24, 40, 10)
	if box.x != 20 || box.y != 7 {
		t.Fatalf("unexpected centered box position: %#v", box)
	}

	inner := box.inset(2, 1)
	if inner.x != 22 || inner.y != 8 || inner.w != 36 || inner.h != 8 {
		t.Fatalf("unexpected inset box: %#v", inner)
	}
}

func TestWrapAndClamp(t *testing.T) {
	lines := wrapAndClamp("", "alpha beta gamma delta", 10, 2)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasSuffix(lines[1], "…") {
		t.Fatalf("expected clamped final line to include ellipsis, got %q", lines[1])
	}

	labeled := wrapAndClamp("Q: ", "one two three", 8, 3)
	if len(labeled) == 0 || !strings.HasPrefix(labeled[0], "Q: ") {
		t.Fatalf("expected labeled wrapping to keep prefix, got %#v", labeled)
	}
}

func TestRemainingRows(t *testing.T) {
	if got := remainingRows(10, 3, 1); got != 7 {
		t.Fatalf("expected 7 remaining rows, got %d", got)
	}
	if got := remainingRows(5, 10, 2); got != 2 {
		t.Fatalf("expected minimum rows when over budget, got %d", got)
	}
	if got := remainingRows(0, 1, 1); got != 0 {
		t.Fatalf("expected 0 rows when no space available, got %d", got)
	}
}
