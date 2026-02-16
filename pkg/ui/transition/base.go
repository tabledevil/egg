package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// BaseTransition provides common functionality for all transitions
type BaseTransition struct{}

// IsCompatible provides a safe default: transitions work everywhere unless they
// opt into stricter checks by overriding this method.
func (b BaseTransition) IsCompatible(c caps.Capabilities) bool { return true }

func (b *BaseTransition) Init() tea.Cmd {
	return Tick()
}

// Tick generates the standard 30 FPS tick
func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

func (b *BaseTransition) SetContent(oldView, newView string) {
	// No-op by default
}

func (b *BaseTransition) Done() bool {
	return false
}
