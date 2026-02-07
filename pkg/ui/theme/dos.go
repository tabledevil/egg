package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"time"
)

type DOSTheme struct {
	blink bool
}

func NewDOSTheme() Theme {
	return &DOSTheme{}
}

func (t *DOSTheme) Init() tea.Cmd {
	return nil
}

func (t *DOSTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if tick, ok := msg.(game.TickMsg); ok {
		// Blink every 500ms based on the global clock
		t.blink = (time.Time(tick).UnixMilli()/500)%2 == 0
	}
	return t, nil
}

func (t *DOSTheme) Name() string        { return "MS-DOS" }
func (t *DOSTheme) Description() string { return "Blue background, gray text" }

func (t *DOSTheme) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile <= termenv.ANSI
}

func (t *DOSTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	bg := lipgloss.Color("#0000AA")
	fg := lipgloss.Color("#AAAAAA")
	white := lipgloss.Color("#FFFFFF")

	baseStyle := lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Width(width).
		Height(height)

	titleBar := lipgloss.NewStyle().
		Background(white).
		Foreground(bg).
		Width(width).
		Align(lipgloss.Center).
		Render(" C:\\WINDOWS\\SYSTEM32\\HACK.EXE ")

	menuBar := lipgloss.NewStyle().
		Background(lipgloss.Color("#AAAAAA")).
		Foreground(lipgloss.Color("#000000")).
		Width(width).
		Render(" File  Edit  Search  Run  Compile  Debug  Options  Help ")

	content := lipgloss.NewStyle().
		Padding(2).
		Width(width).
		Render(fmt.Sprintf("\nQuestion ID: %04X\n\n%s", q.ID, q.Text))

	cursor := "_"
	if !t.blink {
		cursor = " "
	}

	inputArea := lipgloss.NewStyle().
		Padding(0, 2).
		Render(fmt.Sprintf("C:\\> %s%s", inputView, cursor))

	footer := lipgloss.NewStyle().
		Background(lipgloss.Color("#AAAAAA")).
		Foreground(lipgloss.Color("#000000")).
		Width(width).
		Align(lipgloss.Center).
		Render(" F1=Help  Alt+X=Exit ")

	// Layout
	mainAreaHeight := height - 3 // Title + Menu + Footer
	if mainAreaHeight < 0 {
		mainAreaHeight = 0
	}

	mainArea := lipgloss.NewStyle().Height(mainAreaHeight).Render(content + "\n\n" + inputArea)

	if hint != "" {
		mainArea = lipgloss.NewStyle().Height(mainAreaHeight).Render(content + "\n\n" + inputArea + "\n\nBad command or filename.\nHint: " + hint)
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			titleBar,
			menuBar,
			mainArea,
			footer,
		),
	)
}

func init() {
	Register(NewDOSTheme)
}
