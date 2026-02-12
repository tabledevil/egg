package transition

import (
	"ctf-tool/pkg/game"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

var ansiRegexp = regexp.MustCompile("\\x1b\\[[0-9;]*[A-Za-z]")

func stripANSICodes(in string) string {
	return ansiRegexp.ReplaceAllString(in, "")
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func buildFrame(label string, width, height int) string {
	var b strings.Builder
	for y := 0; y < height; y++ {
		line := fmt.Sprintf("%s LINE %02d", label, y)
		if len(line) < width {
			line += strings.Repeat(" ", width-len(line))
		} else if len(line) > width {
			line = line[:width]
		}
		b.WriteString(line)
		if y < height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func TestTransitionsRenderAndComplete(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatalf("expected transition registry to be populated")
	}

	width, height := 60, 20
	oldView := buildFrame("OLD", width, height)
	newView := buildFrame("NEW", width, height)

	for _, constructor := range Registry {
		instance := constructor()
		name := fmt.Sprintf("transition_%T", instance)

		t.Run(name, func(t *testing.T) {
			tr := constructor()
			tr.SetContent(oldView, newView)

			if cmd := tr.Init(); cmd == nil {
				t.Errorf("Init returned nil command")
			}

			view := tr.View(width, height)
			clean := stripANSICodes(view)
			if got := countLines(clean); got != height {
				t.Fatalf("expected %d lines, got %d", height, got)
			}

			const maxTicks = 1000
			done := tr.Done()
			for tick := 0; tick < maxTicks && !done; tick++ {
				next, cmd := tr.Update(game.TickMsg(time.Now()))
				if cmd == nil {
					t.Fatalf("update for %s did not request another tick", name)
				}
				if next != nil {
					tr = next
				}
				done = tr.Done()
			}

			if !tr.Done() {
				t.Fatalf("transition %s never finished after %d ticks", name, maxTicks)
			}
		})
	}
}
