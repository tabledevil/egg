package theme

import (
	"ctf-tool/pkg/game"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type ConsoleConfig struct {
	Name        string
	Desc        string
	BgColor     string
	FgColor     string
	BorderColor string
	TitleColor  string
	BorderStyle lipgloss.Border
	FontBold    bool
}

type ConsoleTheme struct {
	Config ConsoleConfig
}

func NewConsoleTheme(cfg ConsoleConfig) Theme {
	return &ConsoleTheme{Config: cfg}
}

func (t *ConsoleTheme) Init() tea.Cmd { return nil }
func (t *ConsoleTheme) Update(msg tea.Msg) (Theme, tea.Cmd) { return t, nil }
func (t *ConsoleTheme) Name() string { return t.Config.Name }
func (t *ConsoleTheme) Description() string { return t.Config.Desc }

func (t *ConsoleTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	bg := lipgloss.Color(t.Config.BgColor)
	fg := lipgloss.Color(t.Config.FgColor)
	border := lipgloss.Color(t.Config.BorderColor)

	base := lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	boxWidth := width - 6
	if boxWidth < 20 { boxWidth = 20 }

	box := lipgloss.NewStyle().
		Border(t.Config.BorderStyle).
		BorderForeground(border).
		Padding(1).
		Width(boxWidth).
		Align(lipgloss.Center)

	if t.Config.FontBold {
		box = box.Bold(true)
	}

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Config.TitleColor)).
		Bold(true).
		Render(fmt.Sprintf("*** %s MODE - LEVEL %d ***", strings.ToUpper(t.Config.Name), q.ID))

	content := fmt.Sprintf("%s\n\n%s\n\nINPUT: %s", header, q.Text, inputView)

	if hint != "" {
		content += "\n\n(HINT: " + hint + ")"
	}

	return base.Render(box.Render(content))
}

func NewNESTheme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "NES", Desc: "8-bit Classic",
		BgColor: "#000000", FgColor: "#FFFFFF", BorderColor: "#FF0000", TitleColor: "#FF0000",
		BorderStyle: lipgloss.DoubleBorder(),
	})
}

func NewGameboyTheme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "Gameboy", Desc: "Dot Matrix Green",
		BgColor: "#8BAC0F", FgColor: "#0F380F", BorderColor: "#306230", TitleColor: "#0F380F",
		BorderStyle: lipgloss.RoundedBorder(),
	})
}

func NewC64Theme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "C64", Desc: "Commodore 64 Blue",
		BgColor: "#40318D", FgColor: "#7B6FBB", BorderColor: "#7B6FBB", TitleColor: "#FFFFFF",
		BorderStyle: lipgloss.ThickBorder(), FontBold: true,
	})
}

func NewAmigaTheme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "Amiga", Desc: "Workbench Style",
		BgColor: "#0055AA", FgColor: "#FFFFFF", BorderColor: "#FF8800", TitleColor: "#FFFFFF",
		BorderStyle: lipgloss.NormalBorder(),
	})
}

func init() {
	Register(NewNESTheme)
	Register(NewGameboyTheme)
	Register(NewC64Theme)
	Register(NewAmigaTheme)
}

func NewAtariTheme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "Atari", Desc: "2600 Style",
		BgColor: "#000000", FgColor: "#D4A017", BorderColor: "#8B4513", TitleColor: "#D4A017",
		BorderStyle: lipgloss.DoubleBorder(),
	})
}

func NewStarWarsTheme() Theme {
	return NewConsoleTheme(ConsoleConfig{
		Name: "Targeting Computer", Desc: "Stay on target",
		BgColor: "#000000", FgColor: "#FF0000", BorderColor: "#FF0000", TitleColor: "#FF0000",
		BorderStyle: lipgloss.NormalBorder(),
	})
}

// Re-registering allows init to run again? No, I need to append to the existing init block or add a new one.
// Go allows multiple init functions in one package, but inside the same file? Yes.
func init() {
	Register(NewAtariTheme)
	Register(NewStarWarsTheme)
}
