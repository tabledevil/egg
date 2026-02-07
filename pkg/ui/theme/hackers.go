package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"strings"
)

type HackersTheme struct{}

func NewHackersTheme() Theme                                { return &HackersTheme{} }
func (t *HackersTheme) Init() tea.Cmd                       { return nil }
func (t *HackersTheme) Update(msg tea.Msg) (Theme, tea.Cmd) { return t, nil }
func (t *HackersTheme) Name() string                        { return "Zero Cool" }
func (t *HackersTheme) Description() string                 { return "Messy, chaotic hacker style" }

func (t *HackersTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile <= termenv.ANSI
}

func (t *HackersTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	bg := lipgloss.Color("#000000")
	fg := lipgloss.Color("#FF3333") // Reddish

	banner := lipgloss.NewStyle().
		Background(fg).
		Foreground(bg).
		Bold(true).
		Width(width).
		Align(lipgloss.Center).
		Render("!!! GIBSON ACCESSED !!!")

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true).
		Render(q.Text)

	input := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Render("PAYLOAD > " + inputView)

	var s strings.Builder
	s.WriteString(banner)
	s.WriteString("\n\n")
	s.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(content))
	s.WriteString("\n\n")
	s.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(input))

	if hint != "" {
		s.WriteString("\n\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Align(lipgloss.Center).Width(width).Render("HINT: " + hint))
	}

	// Add some random garbage
	garbage := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render("\n\nheap dump: 0x443 0x552 0x111 ... SEGMENTATION FAULT PREVENTED")
	s.WriteString(garbage)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s.String())
}

func init() {
	Register(NewHackersTheme)
}
