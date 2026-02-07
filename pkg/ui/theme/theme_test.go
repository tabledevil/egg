package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestThemes(t *testing.T) {
	// Debug constants

	q := &game.Question{
		ID:     1,
		Text:   "What is the capital of cyberspace?",
		Answer: "None",
		Hint:   "It's everywhere",
	}

	for _, constructor := range Registry {
		th := constructor()
		name := th.Name()

		t.Run(name, func(t *testing.T) {
			fullCaps := caps.Capabilities{
				ColorProfile:  termenv.TrueColor,
				HasUnicode:    true,
				IsInteractive: true,
			}

			if !th.IsCompatible(fullCaps) {
				t.Errorf("Theme %s should be compatible with full capabilities (Profile: %d >= ANSI: %d)", name, fullCaps.ColorProfile, termenv.ANSI)
			}

			// Test View rendering
			view := th.View(80, 24, q, "InputText", "HintText")
			if len(view) == 0 {
				t.Errorf("Theme %s View returned empty string", name)
			}

			// Basic content checks
			if name != "Setec Astronomy" {
				if !strings.Contains(view, "What is the capital") {
					// Some themes rely on async ticks (Matrix, Sneakers), might not show text immediately if logic prevents it?
					// Matrix shows "WAKE UP NEO" and then q.Text.
					// Let's relax this check or verify specific themes logic.
					// Matrix Theme: Always shows q.Text
					// Cyber Theme: Always shows q.Text
					// Minimal: Always shows q.Text
					// Retro: Always shows q.Text
					// Hackers: Always shows q.Text
					// DOS: Always shows q.Text
					// Console: Always shows q.Text
					t.Errorf("Theme %s View does not contain question text", name)
				}
			}
		})
	}
}

func TestMinimalThemeCompatibility(t *testing.T) {
	min := NewMinimalTheme()
	asciiCaps := caps.Capabilities{
		ColorProfile: termenv.Ascii,
		HasUnicode:   false,
	}
	if !min.IsCompatible(asciiCaps) {
		t.Errorf("Minimal theme must be compatible with Ascii/NoUnicode")
	}
}

func TestMatrixThemeRequirements(t *testing.T) {
	mat := NewMatrixTheme()
	asciiCaps := caps.Capabilities{
		ColorProfile: termenv.Ascii,
	}
	if mat.IsCompatible(asciiCaps) {
		t.Errorf("Matrix theme should NOT be compatible with Ascii")
	}
}
