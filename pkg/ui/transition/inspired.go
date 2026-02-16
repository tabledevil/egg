package transition

import (
	"math"

	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"ctf-tool/pkg/ui/caps"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// DecryptLockTransition is inspired by the egg2 decrypt-lock style, but remains
// a fully automatic non-interactive transition.
type DecryptLockTransition struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	frame    int
	done     bool
}

func NewDecryptLockTransition() Transition { return &DecryptLockTransition{} }

func (t *DecryptLockTransition) IsCompatible(c caps.Capabilities) bool {
	return profileAtLeastInspired(c.ColorProfile, termenv.ANSI)
}

func (t *DecryptLockTransition) SetContent(oldView, newView string) {
	t.oldLines = getLines(oldView)
	t.newLines = getLines(newView)
}

func (t *DecryptLockTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.04
		t.frame++
		if t.progress >= 1.12 {
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *DecryptLockTransition) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	faint := lipgloss.NewStyle().Foreground(lipgloss.Color("#30585B"))
	lock := lipgloss.NewStyle().Foreground(lipgloss.Color("#65FFDD")).Bold(true)
	final := lipgloss.NewStyle().Foreground(lipgloss.Color("#E8FFF9"))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			threshold := float64((x*17+y*29)%100) / 100.0
			switch {
			case t.progress > threshold+0.10:
				c.SetChar(x, y, charAt(t.newLines, x, y), final)
			case t.progress > threshold-0.08:
				c.SetChar(x, y, noiseGlyph(t.frame+x+y), lock)
			default:
				c.SetChar(x, y, charAt(t.oldLines, x, y), faint)
			}
		}
	}

	if height > 0 {
		barW := width - 2
		if barW < 3 {
			barW = width
		}
		if barW > 0 {
			pct := clampFloatInspired(t.progress/1.0, 0, 1)
			fill := int(float64(barW) * pct)
			if fill < 0 {
				fill = 0
			}
			if fill > barW {
				fill = barW
			}
			bar := "[" + repeatRune('=', fill) + repeatRune('-', barW-fill) + "]"
			c.SetString(0, height-1, bar, lock)
		}
	}

	if width > 18 {
		c.SetString(2, 0, "DECRYPT LOCK", lock)
	}

	return c.Render()
}

func (t *DecryptLockTransition) Done() bool { return t.done }

func (t *DecryptLockTransition) ensureLines(h int) {
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

// PacketFirewallTransition is inspired by egg2 firewall/packet visuals.
type PacketFirewallTransition struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	frame    int
	done     bool
}

func NewPacketFirewallTransition() Transition { return &PacketFirewallTransition{} }

func (t *PacketFirewallTransition) IsCompatible(c caps.Capabilities) bool {
	return profileAtLeastInspired(c.ColorProfile, termenv.ANSI)
}

func (t *PacketFirewallTransition) SetContent(oldView, newView string) {
	t.oldLines = getLines(oldView)
	t.newLines = getLines(newView)
}

func (t *PacketFirewallTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.035
		t.frame++
		if t.progress >= 1.18 {
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *PacketFirewallTransition) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	faint := lipgloss.NewStyle().Foreground(lipgloss.Color("#365157"))
	clean := lipgloss.NewStyle().Foreground(lipgloss.Color("#D7FFF1"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5C5C")).Bold(true)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#56FF8E")).Bold(true)

	revealRows := int(float64(height) * clampFloatInspired(t.progress, 0, 1))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if y < revealRows {
				c.SetChar(x, y, charAt(t.newLines, x, y), clean)
			} else {
				c.SetChar(x, y, charAt(t.oldLines, x, y), faint)
			}
		}
	}

	wallX := width / 2
	if width > 8 {
		wallX += int(math.Sin(float64(t.frame)/5.0) * 2)
		if wallX < 3 {
			wallX = 3
		}
		if wallX > width-3 {
			wallX = width - 3
		}
	}

	for y := 0; y < height; y++ {
		c.SetChar(wallX, y, '|', red)
	}

	packetCount := 5
	if height < packetCount {
		packetCount = height
	}
	for p := 0; p < packetCount; p++ {
		y := (p*3 + t.frame) % maxIntInspired(height, 1)
		x := (t.frame*3 + p*7) % maxIntInspired(wallX-1, 1)
		if x >= wallX-1 {
			c.SetString(maxIntInspired(wallX-2, 0), y, "=>", red)
			if wallX+2 < width {
				if (p+t.frame)%2 == 0 {
					c.SetString(wallX+2, y, "PASS", green)
				} else {
					c.SetString(wallX+2, y, "DROP", red)
				}
			}
		} else {
			c.SetString(x, y, ">>>", red)
		}
	}

	if width > 24 {
		c.SetString(2, 0, "FIREWALL PACKET SWEEP", green)
	}

	return c.Render()
}

func (t *PacketFirewallTransition) Done() bool { return t.done }

func (t *PacketFirewallTransition) ensureLines(h int) {
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

func charAt(lines []string, x, y int) rune {
	if y < 0 || y >= len(lines) || x < 0 {
		return ' '
	}
	r := []rune(lines[y])
	if x >= len(r) {
		return ' '
	}
	return r[x]
}

func noiseGlyph(seed int) rune {
	glyphs := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&@")
	if len(glyphs) == 0 {
		return '#'
	}
	if seed < 0 {
		seed = -seed
	}
	return glyphs[seed%len(glyphs)]
}

func repeatRune(r rune, count int) string {
	if count <= 0 {
		return ""
	}
	buf := make([]rune, count)
	for i := range buf {
		buf[i] = r
	}
	return string(buf)
}

func clampFloatInspired(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func maxIntInspired(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func profileAtLeastInspired(profile termenv.Profile, required termenv.Profile) bool {
	return profileRankInspired(profile) >= profileRankInspired(required)
}

func profileRankInspired(profile termenv.Profile) int {
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

func init() {
	Register(NewDecryptLockTransition)
	Register(NewPacketFirewallTransition)
}
