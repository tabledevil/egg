package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// BaseTheme provides common functionality for all themes
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

// Tick generates the standard 30 FPS tick
func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}
