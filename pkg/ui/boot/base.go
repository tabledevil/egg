package boot

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"

	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type BaseIntro struct{}

// Safe default for intro compatibility.
func (b BaseIntro) IsCompatible(c caps.Capabilities) bool {
	return true
}

func (b *BaseIntro) Init() tea.Cmd {
	return Tick()
}

func (b *BaseIntro) Update(msg tea.Msg) (Intro, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		return nil, Tick()
	}
	return nil, nil
}

func (b *BaseIntro) Done() bool {
	return false
}

func Tick() tea.Cmd {
	return tea.Tick(caps.BootFrameInterval(), func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}
