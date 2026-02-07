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
	IsCompatible(c caps.Capabilities) bool
}

type Constructor func() Theme

var Registry = []Constructor{}

func Register(c Constructor) {
	Registry = append(Registry, c)
}
