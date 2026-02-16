package transition

import (
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
)

// Transition represents a transition animation or mini-game between levels
type Transition interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Transition, tea.Cmd)
	View(width, height int) string
	Done() bool
	SetContent(oldView, newView string)
}

// CapabilityAware is an optional interface that transitions can implement to
// declare whether they should be used in the current terminal environment.
//
// Transitions that don't implement this are treated as always compatible.
type CapabilityAware interface {
	IsCompatible(c caps.Capabilities) bool
}

// Constructor for creating new transition instances
type Constructor func() Transition

var Registry = []Constructor{}

func Register(c Constructor) {
	Registry = append(Registry, c)
}
