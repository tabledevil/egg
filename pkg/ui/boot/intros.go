package boot

import (
	"fmt"
	"math"
	"strings"

	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type NeonCipherIntro struct {
	BaseIntro
	frame int
	done  bool
}

func NewNeonCipherIntro() Intro { return &NeonCipherIntro{} }

func (i *NeonCipherIntro) Name() string        { return "Neon Cipher" }
func (i *NeonCipherIntro) Description() string { return "Spectral handshake with decrypt pulse" }

func (i *NeonCipherIntro) IsCompatible(c caps.Capabilities) bool {
	return c.IsInteractive && c.HasUnicode && profileAtLeast(c.ColorProfile, termenv.ANSI256)
}

func (i *NeonCipherIntro) Update(msg tea.Msg) (Intro, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if i.done {
			return i, nil
		}
		i.frame++
		if i.frame >= 150 {
			i.done = true
		}
		return i, Tick()
	}
	return i, nil
}

func (i *NeonCipherIntro) Done() bool { return i.done }

func (i *NeonCipherIntro) View(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	panelW := panelWidth(width, 8, 42, 88)
	if panelW < 1 {
		return ""
	}

	palette := []lipgloss.Color{"#00FFD5", "#2DE2E6", "#6AFFD9", "#00B8FF"}
	accent := palette[(i.frame/8)%len(palette)]

	script := []string{
		"SYNTHETIC SIGNAL BUS: ONLINE",
		"ROTATING KEYMESH SHARDS... OK",
		"ALIGNING TERMINAL CHROMA... LOCKED",
		"PROMPT CHANNEL: READY",
	}

	logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A7FFF2"))
	cursorStyle := lipgloss.NewStyle().Foreground(accent).Bold(true)

	logs := make([]string, 0, len(script))
	for idx, line := range script {
		reveal := i.frame - idx*16
		if reveal <= 0 {
			continue
		}
		text := prefixRunes(line, reveal)
		if !i.done && idx == (i.frame/16) && i.frame%2 == 0 {
			text += cursorStyle.Render("█")
		}
		logs = append(logs, logStyle.Render("> "+text))
	}
	if len(logs) == 0 {
		logs = append(logs, logStyle.Render(">"))
	}

	gridW := panelW - 6
	if gridW < 10 {
		gridW = 10
	}
	gridStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#1D6760"))
	var gridLines []string
	for y := 0; y < 4; y++ {
		glyph := "."
		if (y+i.frame)%2 == 0 {
			glyph = ":"
		}
		line := strings.Repeat(glyph+" ", maxInt(1, gridW/2))
		gridLines = append(gridLines, gridStyle.Render(line))
	}

	progress := float64(i.frame) / 150.0
	bar := progressBar(maxInt(10, panelW-10), progress, '#', '-')
	bar = lipgloss.NewStyle().Foreground(accent).Render(bar)

	head := lipgloss.NewStyle().Foreground(accent).Bold(true).Render(":: NEON CIPHER BOOT ::")
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		head,
		"",
		strings.Join(gridLines, "\n"),
		"",
		strings.Join(logs, "\n"),
		"",
		bar,
	)

	card := lipgloss.NewStyle().
		Width(panelW).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(accent).
		Background(lipgloss.Color("#021317")).
		Foreground(lipgloss.Color("#D6FFF9")).
		Padding(1).
		Render(body)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type AmberGridIntro struct {
	BaseIntro
	frame int
	done  bool
}

func NewAmberGridIntro() Intro { return &AmberGridIntro{} }

func (i *AmberGridIntro) Name() string        { return "Amber Grid" }
func (i *AmberGridIntro) Description() string { return "CRT diagnostics with warm amber telemetry" }

func (i *AmberGridIntro) IsCompatible(c caps.Capabilities) bool {
	return c.IsInteractive && profileAtLeast(c.ColorProfile, termenv.ANSI)
}

func (i *AmberGridIntro) Update(msg tea.Msg) (Intro, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if i.done {
			return i, nil
		}
		i.frame++
		if i.frame >= 120 {
			i.done = true
		}
		return i, Tick()
	}
	return i, nil
}

func (i *AmberGridIntro) Done() bool { return i.done }

func (i *AmberGridIntro) View(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	panelW := panelWidth(width, 6, 40, 84)
	if panelW < 1 {
		return ""
	}

	amber := lipgloss.Color("#FFB547")
	text := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD89A"))
	accent := lipgloss.NewStyle().Foreground(amber).Bold(true)

	checks := []string{
		"MEMORY MAP............... OK",
		"I/O BUS INTEGRITY........ OK",
		"CHROMA CALIBRATION....... OK",
		"PROMPT RELAY............. READY",
	}

	visible := minInt(len(checks), i.frame/18+1)
	var out []string
	for idx := 0; idx < visible; idx++ {
		line := checks[idx]
		reveal := i.frame - idx*18
		line = prefixRunes(line, reveal)
		out = append(out, text.Render(line))
	}

	scanWidth := maxInt(10, panelW-12)
	scanPos := 0
	if scanWidth > 0 {
		scanPos = (i.frame * 2) % scanWidth
	}
	scan := strings.Repeat(".", scanPos) + "|" + strings.Repeat(".", maxInt(0, scanWidth-scanPos-1))
	scan = accent.Render(scan)

	progress := float64(i.frame) / 120.0
	bar := progressBar(maxInt(12, panelW-12), progress, '=', ' ')
	bar = accent.Render(bar)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		accent.Render("AMBER GRID DIAGNOSTICS"),
		"",
		strings.Join(out, "\n"),
		"",
		text.Render("SCAN: "+scan),
		text.Render("LOAD: "+bar),
	)

	card := lipgloss.NewStyle().
		Width(panelW).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(amber).
		Background(lipgloss.Color("#1A1205")).
		Foreground(lipgloss.Color("#FFE7BE")).
		Padding(1).
		Render(body)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

type PrismPulseIntro struct {
	BaseIntro
	frame int
	done  bool
}

func NewPrismPulseIntro() Intro { return &PrismPulseIntro{} }

func (i *PrismPulseIntro) Name() string        { return "Prism Pulse" }
func (i *PrismPulseIntro) Description() string { return "TrueColor burst with shifting spectral bars" }

func (i *PrismPulseIntro) IsCompatible(c caps.Capabilities) bool {
	return c.IsInteractive && profileAtLeast(c.ColorProfile, termenv.TrueColor)
}

func (i *PrismPulseIntro) Update(msg tea.Msg) (Intro, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if i.done {
			return i, nil
		}
		i.frame++
		if i.frame >= 132 {
			i.done = true
		}
		return i, Tick()
	}
	return i, nil
}

func (i *PrismPulseIntro) Done() bool { return i.done }

func (i *PrismPulseIntro) View(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	panelW := panelWidth(width, 6, 44, 88)
	if panelW < 1 {
		return ""
	}

	phase := float64(i.frame) / 8.0
	colors := []lipgloss.Color{"#FF4E7A", "#8A6BFF", "#2DD5FF", "#4CFFB5"}

	barW := maxInt(12, panelW-12)
	bands := make([]string, 0, 4)
	for band := 0; band < 4; band++ {
		center := (math.Sin(phase+float64(band)*0.9) + 1) * 0.5
		fill := int(center * float64(barW-1))
		line := strings.Repeat("·", fill) + "◆" + strings.Repeat("·", maxInt(0, barW-fill-1))
		bands = append(bands, lipgloss.NewStyle().Foreground(colors[band%len(colors)]).Render(line))
	}

	progress := float64(i.frame) / 132.0
	headColor := colors[(i.frame/9)%len(colors)]
	head := lipgloss.NewStyle().Foreground(headColor).Bold(true).Render("PRISM PULSE SYNCHRONIZER")
	foot := lipgloss.NewStyle().Foreground(lipgloss.Color("#DCE6FF")).Render(
		fmt.Sprintf("HARMONIC LOCK: %3d%%", int(clampFloat(progress, 0, 1)*100)),
	)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		head,
		"",
		strings.Join(bands, "\n"),
		"",
		foot,
	)

	card := lipgloss.NewStyle().
		Width(panelW).
		Border(lipgloss.ThickBorder()).
		BorderForeground(headColor).
		Background(lipgloss.Color("#0C1022")).
		Foreground(lipgloss.Color("#EEF3FF")).
		Padding(1).
		Render(body)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
}

func progressBar(width int, progress float64, fill, empty rune) string {
	if width <= 0 {
		return ""
	}
	progress = clampFloat(progress, 0, 1)
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return "[" + strings.Repeat(string(fill), filled) + strings.Repeat(string(empty), width-filled) + "]"
}

func prefixRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if n >= len(r) {
		return s
	}
	return string(r[:n])
}

func panelWidth(width, margin, minPreferred, maxPreferred int) int {
	w := width - margin
	if w > maxPreferred {
		w = maxPreferred
	}
	if w < minPreferred {
		w = minPreferred
	}
	if w > width {
		w = width
	}
	if w < 1 {
		w = 1
	}
	return w
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func profileAtLeast(profile termenv.Profile, required termenv.Profile) bool {
	return profileRank(profile) >= profileRank(required)
}

func profileRank(profile termenv.Profile) int {
	switch profile {
	case termenv.Ascii:
		return 0
	case termenv.ANSI:
		return 1
	case termenv.ANSI256:
		return 2
	case termenv.TrueColor:
		return 3
	default:
		return -1
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func init() {
	Register(NewNeonCipherIntro)
	Register(NewAmberGridIntro)
	Register(NewPrismPulseIntro)
}
