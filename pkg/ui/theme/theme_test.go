package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

var ansiPattern = regexp.MustCompile("\\x1b\\[[0-9;]*[A-Za-z]")

func stripANSI(in string) string {
	return ansiPattern.ReplaceAllString(in, "")
}

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func TestThemesRenderQuestionAndInput(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatalf("expected at least one registered theme")
	}

	q := &game.Question{ID: 1, Text: "What is the access code?", Hint: "Try XOR"}
	input := "hunter2"
	width, height := 80, 24

	for _, constructor := range Registry {
		instance := constructor()
		name := instance.Name()
		if name == "" {
			name = fmt.Sprintf("theme_%T", instance)
		}

		t.Run(name, func(t *testing.T) {
			theme := constructor()

			if cmd := theme.Init(); cmd == nil {
				t.Errorf("Init returned nil command")
			}

			if _, cmd := theme.Update(game.TickMsg(time.Now())); cmd == nil {
				t.Errorf("Update did not request another tick")
			}

			view := theme.View(width, height, q, input, q.Hint)
			clean := stripANSI(view)

			if got := lineCount(clean); got < height {
				t.Fatalf("expected at least %d lines, got %d", height, got)
			}

			upperView := strings.ToUpper(clean)
			question := strings.ToUpper(q.Text)
			if !strings.Contains(upperView, question) {
				if !(name == "Citadel" && strings.Contains(upperView, "INSECT")) {
					t.Errorf("theme output missing question text: %q", q.Text)
				}
			}

			if !strings.Contains(clean, input) {
				t.Errorf("theme output missing input view text %q", input)
			}
		})
	}
}
