package boot

import (
	"testing"
	"time"

	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"

	"github.com/muesli/termenv"
)

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	count := 1
	for _, r := range s {
		if r == '\n' {
			count++
		}
	}
	return count
}

func TestBootProfilesRenderAndComplete(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatalf("expected boot profile registry to be populated")
	}

	const (
		width    = 72
		height   = 20
		maxTicks = 500
	)

	for _, constructor := range Registry {
		profile := constructor()
		name := profile.Name()
		if name == "" {
			name = "boot_profile"
		}

		t.Run(name, func(t *testing.T) {
			intro := constructor()
			if cmd := intro.Init(); cmd == nil {
				t.Fatalf("Init returned nil command")
			}

			if got := lineCount(intro.View(width, height)); got == 0 {
				t.Fatalf("expected rendered output for %s", name)
			}

			done := intro.Done()
			for tick := 0; tick < maxTicks && !done; tick++ {
				next, cmd := intro.Update(game.TickMsg(time.Now()))
				if cmd == nil {
					t.Fatalf("Update did not request next tick")
				}
				if next != nil {
					intro = next
				}
				_ = intro.View(width, height)
				done = intro.Done()
			}

			if !intro.Done() {
				t.Fatalf("boot profile %s did not finish within %d ticks", name, maxTicks)
			}
		})
	}
}

func TestBootProfileCompatibility(t *testing.T) {
	asciiCaps := caps.Capabilities{ColorProfile: termenv.Ascii, HasUnicode: false, IsInteractive: true}
	ansiCaps := caps.Capabilities{ColorProfile: termenv.ANSI, HasUnicode: true, IsInteractive: true}
	trueColorCaps := caps.Capabilities{ColorProfile: termenv.TrueColor, HasUnicode: true, IsInteractive: true}

	neonAware := NewNeonCipherIntro().(CapabilityAware)
	if neonAware.IsCompatible(asciiCaps) {
		t.Fatalf("neon profile should not support ascii-only terminals")
	}

	amberAware := NewAmberGridIntro().(CapabilityAware)
	if !amberAware.IsCompatible(ansiCaps) {
		t.Fatalf("amber profile should support ANSI terminals")
	}

	prismAware := NewPrismPulseIntro().(CapabilityAware)
	if prismAware.IsCompatible(ansiCaps) {
		t.Fatalf("prism profile should require truecolor terminal")
	}
	if !prismAware.IsCompatible(trueColorCaps) {
		t.Fatalf("prism profile should support truecolor terminal")
	}
}
