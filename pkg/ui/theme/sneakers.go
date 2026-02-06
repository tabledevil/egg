package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"time"
)

type SneakersTheme struct {
	tickCount int
}

func NewSneakersTheme() Theme { return &SneakersTheme{} }

func (t *SneakersTheme) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

func (t *SneakersTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.tickCount++
		return t, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return game.TickMsg(t)
		})
	}
	return t, nil
}

func (t *SneakersTheme) Name() string { return "Setec Astronomy" }
func (t *SneakersTheme) Description() string { return "Too Many Secrets" }

func (t *SneakersTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	boxWidth := 60
	if boxWidth > width-4 { boxWidth = width-4 }

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF")).
		Padding(1).
		Width(boxWidth).
		Align(lipgloss.Center)

	header := "NO MORE SECRETS"

	// Animate Text Scramble
	// For the first 20 ticks (1 sec), show scrambled text
	text := q.Text
	if t.tickCount < 20 {
		runes := []rune(text)
		for i := range runes {
			if runes[i] != ' ' && rand.Float32() < 0.5 {
				runes[i] = rune('A' + rand.Intn(26))
			}
		}
		text = string(runes)
	}

	content := fmt.Sprintf("%s\n\n%s\n\nCODE: %s", header, text, inputView)
	if hint != "" {
		content += "\n\n[" + hint + "]"
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

func init() {
	Register(NewSneakersTheme)
}
