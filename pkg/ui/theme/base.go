package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// BaseTheme provides common functionality for all themes.
// Themes that inherit BaseTheme without overriding Update get 30 FPS ticks.
type BaseTheme struct{}

// IsCompatible provides a safe default: themes work everywhere unless they opt
// into stricter checks by overriding this method.
func (b BaseTheme) IsCompatible(c caps.Capabilities) bool { return true }

func (b *BaseTheme) Init() tea.Cmd {
	return Tick()
}

func (b *BaseTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		return nil, Tick()
	}
	return nil, nil
}

// SlowBaseTheme is identical to BaseTheme but generates ticks at ~10 FPS.
// Embed this instead of BaseTheme in themes that have little to no animation
// (static layouts, cursor-blink only) to reduce CPU and ANSI output volume.
type SlowBaseTheme struct{}

func (b SlowBaseTheme) IsCompatible(c caps.Capabilities) bool { return true }

func (b *SlowBaseTheme) Init() tea.Cmd {
	return SlowTick()
}

func (b *SlowBaseTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		return nil, SlowTick()
	}
	return nil, nil
}

// Tick generates the standard 30 FPS tick
func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

// SlowTick generates a 10 FPS tick for themes with minimal animation.
func SlowTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}
