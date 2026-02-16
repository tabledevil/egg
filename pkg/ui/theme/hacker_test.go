package theme

import (
	"ctf-tool/pkg/game"
	"strings"
	"testing"
	"time"
)

func advanceCryptoFrames(theme *CryptoTheme, frames int) *CryptoTheme {
	current := theme
	for i := 0; i < frames; i++ {
		next, _ := current.Update(game.TickMsg(time.Now()))
		if cast, ok := next.(*CryptoTheme); ok {
			current = cast
		}
	}
	return current
}

func extractClueLine(view string) string {
	clean := stripANSI(view)
	for _, line := range strings.Split(clean, "\n") {
		idx := strings.Index(line, "CLUE:")
		if idx >= 0 {
			segment := line[idx:]
			if end := strings.IndexRune(segment, '│'); end >= 0 {
				segment = segment[:end]
			}
			return strings.TrimRight(segment, " ")
		}
	}
	return ""
}

func TestCryptoThemeHintStaticWhenFits(t *testing.T) {
	theme, ok := NewCryptoTheme().(*CryptoTheme)
	if !ok {
		t.Fatalf("expected NewCryptoTheme to return *CryptoTheme")
	}

	q := &game.Question{ID: 1, Text: "Hash challenge", Hint: "Try terminal command"}
	width, height := 120, 28

	first := extractClueLine(theme.View(width, height, q, "nonce123", q.Hint))
	if first == "" {
		t.Fatalf("expected clue line to be visible")
	}

	theme = advanceCryptoFrames(theme, 12)
	second := extractClueLine(theme.View(width, height, q, "nonce123", q.Hint))
	if second == "" {
		t.Fatalf("expected clue line to remain visible")
	}

	if first != second {
		t.Fatalf("expected fitting hint to remain static across frames")
	}
}

func TestCryptoThemeHintScrollsSlowlyWhenOverflowing(t *testing.T) {
	theme, ok := NewCryptoTheme().(*CryptoTheme)
	if !ok {
		t.Fatalf("expected NewCryptoTheme to return *CryptoTheme")
	}

	q := &game.Question{ID: 1, Text: "Hash challenge", Hint: "This hint is intentionally very long so it cannot fit in the blockchain clue panel and should scroll slowly"}
	width, height := 80, 24

	start := extractClueLine(theme.View(width, height, q, "nonce123", q.Hint))
	if start == "" {
		t.Fatalf("expected clue line to be visible")
	}

	theme = advanceCryptoFrames(theme, 3)
	soon := extractClueLine(theme.View(width, height, q, "nonce123", q.Hint))
	if soon != start {
		t.Fatalf("expected overflowing hint to scroll slower than every frame")
	}

	theme = advanceCryptoFrames(theme, 24)
	later := extractClueLine(theme.View(width, height, q, "nonce123", q.Hint))
	if later == start {
		t.Fatalf("expected overflowing hint to eventually scroll")
	}
}
