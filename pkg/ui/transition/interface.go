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
	IsCompatible(c caps.Capabilities) bool
}

// Constructor for creating new transition instances
type Constructor func() Transition

var Registry = []Constructor{}

func Register(c Constructor) {
	Registry = append(Registry, c)
}
