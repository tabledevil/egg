package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
)

type Theme interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Theme, tea.Cmd)
	View(width, height int, q *game.Question, inputView string, hint string) string
	Name() string
	Description() string
}

// CapabilityAware is an optional interface that themes can implement to declare
// whether they should be used in the current terminal environment.
//
// Themes that don't implement this are treated as always compatible.
type CapabilityAware interface {
	IsCompatible(c caps.Capabilities) bool
}

type Constructor func() Theme

var Registry = []Constructor{}

func Register(c Constructor) {
	Registry = append(Registry, c)
}
