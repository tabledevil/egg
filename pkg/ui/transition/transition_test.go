package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func TestTransitions(t *testing.T) {
	fullCaps := caps.Capabilities{
		ColorProfile:  termenv.TrueColor,
		HasUnicode:    true,
		IsInteractive: true,
	}

	for _, constructor := range Registry {
		tr := constructor()

		// 1. Init
		cmd := tr.Init()
		if cmd == nil {
			// Some transitions might return nil Init, which is valid.
		}

		// 2. Compatibility
		if !tr.IsCompatible(fullCaps) {
			t.Errorf("Transition should be compatible with full capabilities")
		}

		// 3. View (Initial)
		view := tr.View(80, 24)
		if len(view) == 0 {
			t.Errorf("Transition View returned empty string initially")
		}

		// Use strings package to avoid unused import error
		if strings.Contains(view, "error") {
			t.Errorf("Transition view contains error")
		}

		// 4. Update Loop Simulation
		for i := 0; i < 100; i++ {
			msg := game.TickMsg(time.Now())
			var newCmd tea.Cmd
			tr, newCmd = tr.Update(msg)
			_ = newCmd

			if tr.Done() {
				break
			}
		}

		view = tr.View(80, 24)
		if len(view) == 0 {
			t.Errorf("Transition View returned empty string after update")
		}
	}
}

func TestLoadingTransitionUnicode(t *testing.T) {
	l := NewLoadingTransition()
	asciiCaps := caps.Capabilities{
		HasUnicode: false,
	}
	if l.IsCompatible(asciiCaps) {
		t.Errorf("Loading transition requires Unicode")
	}
}

func TestTicTacToeCompatibility(t *testing.T) {
	tt := NewTicTacToeTransition()
	asciiCaps := caps.Capabilities{
		ColorProfile: termenv.Ascii,
		HasUnicode:   false,
	}
	if !tt.IsCompatible(asciiCaps) {
		t.Errorf("TicTacToe should be compatible everywhere")
	}
}
