package boot

import (
	"ctf-tool/pkg/ui/caps"

	tea "github.com/charmbracelet/bubbletea"
)

// Intro represents a cinematic startup sequence shown before the first question.
type Intro interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Intro, tea.Cmd)
	View(width, height int) string
	Done() bool
	Name() string
	Description() string
}

// CapabilityAware is optional and lets intros opt out for unsupported terminals.
type CapabilityAware interface {
	IsCompatible(c caps.Capabilities) bool
}

type Constructor func() Intro

var Registry = []Constructor{}

func Register(c Constructor) {
	Registry = append(Registry, c)
}
