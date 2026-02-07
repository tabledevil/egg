package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

type LoadingTransition struct {
	progress float64
}

func NewLoadingTransition() Transition {
	return &LoadingTransition{progress: 0}
}

func (t *LoadingTransition) Init() tea.Cmd {
	return tick()
}

func (t *LoadingTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	switch msg.(type) {
	case game.TickMsg:
		t.progress += 0.005
		return t, tick()
	case tea.KeyMsg:
		t.progress += 0.05 // Mash keys to speed up
		return t, nil
	}
	return t, nil
}

func (t *LoadingTransition) View(width, height int) string {
	barWidth := 40
	if barWidth > width-4 {
		barWidth = width - 4
	}
	if barWidth < 0 {
		barWidth = 0
	}

	filled := int(t.progress * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	style := lipgloss.NewStyle().
		Width(width).Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#FF00FF"))

	return style.Render(fmt.Sprintf("ESTABLISHING CONNECTION...\n\n[%s]", bar))
}

func (t *LoadingTransition) Done() bool {
	return t.progress >= 1.0
}

func (t *LoadingTransition) IsCompatible(c caps.Capabilities) bool {
	// Unicode for blocks
	return c.HasUnicode
}

func init() {
	Register(NewLoadingTransition)
}
