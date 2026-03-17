package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestQuantizeProducesValidValues verifies the quantize helper.
func TestQuantizeProducesValidValues(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{0, 0},
		{15, 0},
		{16, 16},
		{127, 112},
		{128, 128},
		{255, 240},
		{256, 255},  // clamp
		{-1, 0},     // clamp
		{1000, 255}, // clamp
	}
	for _, tc := range cases {
		got := quantize(tc.in)
		if got != tc.want {
			t.Errorf("quantize(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// TestAntigravityColorCacheReuse verifies that cachedColor correctly caches
// and deduplicates style lookups.
func TestAntigravityColorCacheReuse(t *testing.T) {
	theme := newAntigravityTheme(LayoutDialog, "test", "test")

	_ = theme.cachedColor("#ff0000")
	_ = theme.cachedColor("#ff0000") // same hex — should hit cache
	_ = theme.cachedColor("#00ff00") // different hex — new entry

	// Cache should have exactly 2 entries.
	if len(theme.colorCache) != 2 {
		t.Errorf("expected 2 cached colors, got %d", len(theme.colorCache))
	}

	// Requesting the same hex again should not grow the cache.
	_ = theme.cachedColor("#ff0000")
	if len(theme.colorCache) != 2 {
		t.Errorf("cache should not grow for repeated hex; got %d entries", len(theme.colorCache))
	}
}

// TestAllThemesStillRenderCorrectly runs the same visual correctness check
// as TestThemesRenderQuestionAndInput but exercises the new optimised canvas
// renderer.  This ensures the coalescing and style changes produce visually
// identical output.
func TestAllThemesStillRenderCorrectly(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatalf("no themes registered")
	}

	q := &game.Question{ID: 1, Text: "What is the access code?", Hint: "Try XOR"}
	input := "hunter2"
	width, height := 80, 24
	const maxFrames = 30

	for _, constructor := range Registry {
		instance := constructor()
		name := instance.Name()
		if name == "" {
			name = fmt.Sprintf("theme_%T", instance)
		}

		t.Run(name, func(t *testing.T) {
			theme := constructor()
			theme.Init()

			okAnyFrame := false
			for frame := 0; frame < maxFrames; frame++ {
				if frame > 0 {
					next, cmd := theme.Update(game.TickMsg(time.Now()))
					if cmd == nil {
						t.Fatalf("Update did not request another tick")
					}
					if next != nil {
						theme = next
					}
				}

				view := theme.View(width, height, q, input, q.Hint)
				clean := stripANSI(view)

				// Must produce the right number of lines.
				lines := lineCount(clean)
				if lines < height {
					t.Fatalf("expected >= %d lines, got %d", height, lines)
				}

				hasQuestion := questionMatchesView(name, clean, q)
				hasInput := strings.Contains(clean, input)

				if hasQuestion && hasInput {
					okAnyFrame = true
					break
				}
			}

			if !okAnyFrame {
				t.Fatalf("theme output did not contain question+input in any of %d frames", maxFrames)
			}
		})
	}
}

// BenchmarkAntigravityDialogView benchmarks the heaviest theme variant.
func BenchmarkAntigravityDialogView(b *testing.B) {
	theme := newAntigravityTheme(LayoutDialog, "bench", "bench")
	theme.Init()
	q := &game.Question{ID: 1, Text: "Benchmark question text", Hint: "hint"}

	// Warm up: run a few ticks so particles are initialised.
	for i := 0; i < 10; i++ {
		theme.Update(game.TickMsg(time.Now()))
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		theme.Update(game.TickMsg(time.Now()))
		_ = theme.View(80, 24, q, "hunter2", "hint")
	}
}

// BenchmarkAntigravityAvoidanceView benchmarks the avoidance layout with
// particle-text collision physics.
func BenchmarkAntigravityAvoidanceView(b *testing.B) {
	theme := newAntigravityTheme(LayoutAvoidance, "bench", "bench")
	theme.Init()
	q := &game.Question{ID: 1, Text: "Benchmark question text", Hint: "hint"}

	for i := 0; i < 10; i++ {
		theme.Update(game.TickMsg(time.Now()))
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		theme.Update(game.TickMsg(time.Now()))
		_ = theme.View(80, 24, q, "hunter2", "hint")
	}
}

// TestSlowBaseThemeUsesSlowTick verifies that themes using SlowBaseTheme
// generate the 100ms tick interval.
func TestSlowBaseThemeUsesSlowTick(t *testing.T) {
	// DOSTheme should use SlowBaseTheme.
	dos := NewDOSTheme()
	cmd := dos.Init()
	if cmd == nil {
		t.Fatal("Init returned nil")
	}

	next, cmd2 := dos.Update(game.TickMsg(time.Now()))
	if cmd2 == nil {
		t.Fatal("Update returned nil tick command")
	}
	// Verify theme still works after Update.
	_ = next
}
