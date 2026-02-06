package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 6. CRT Power Off ---

type CRTPowerOff struct {
	BaseTransition
	oldLines []string
	vScale   float64
	hScale   float64
	phase    int
	done     bool
}

func NewCRTPowerOff() Transition {
	return &CRTPowerOff{vScale: 1.0, hScale: 1.0}
}

func (t *CRTPowerOff) SetContent(o, n string) {
	t.oldLines = getLines(o)
}

func (t *CRTPowerOff) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		switch t.phase {
		case 0: // Vertical compress
			t.vScale -= 0.1
			if t.vScale <= 0.05 { t.phase = 1 }
		case 1: // Horizontal compress
			t.hScale -= 0.1
			if t.hScale <= 0.02 { t.phase = 2 }
		case 2: // Fade out dot
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *CRTPowerOff) View(width, height int) string {
	if t.phase == 2 { return "" } // Black screen

	c := canvas.New(width, height)
	t.ensureLines(height)

	visibleH := int(float64(height) * t.vScale)
	visibleW := int(float64(width) * t.hScale)

	if visibleH < 1 { visibleH = 1 }
	if visibleW < 1 { visibleW = 1 }

	startX := (width - visibleW) / 2
	startY := (height - visibleH) / 2

	bright := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	for y := 0; y < visibleH; y++ {
		// Sample from original
		srcY := int(float64(y) / t.vScale)
		if srcY >= height { srcY = height - 1 }

		line := ""
		if srcY < len(t.oldLines) { line = t.oldLines[srcY] }

		for x := 0; x < visibleW; x++ {
			srcX := int(float64(x) / t.hScale)
			char := ' '
			if srcX < len(line) { char = rune(line[srcX]) }

			c.SetChar(startX+x, startY+y, char, bright)
		}
	}
	return c.Render()
}

func (t *CRTPowerOff) Done() bool { return t.done }
func (t *CRTPowerOff) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
}


// --- 7. Data Wipe ---

type DataWipe struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	pass     int
	done     bool
}

func NewDataWipe() Transition { return &DataWipe{} }
func (t *DataWipe) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *DataWipe) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.1
		if t.progress >= 1.0 {
			t.progress = 0
			t.pass++
			if t.pass > 2 { t.done = true }
		}
		return t, Tick()
	}
	return t, nil
}

func (t *DataWipe) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	wipeY := int(float64(height) * t.progress)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	for y := 0; y < height; y++ {
		line := ""

		if t.pass == 0 {
			// First pass: Overwrite old with garbage
			if y < len(t.oldLines) { line = t.oldLines[y] }
			if y < wipeY {
				// Garbage
				chars := []rune(line)
				for i := range chars { chars[i] = rune(33 + rand.Intn(94)) }
				line = string(chars)
			}
		} else if t.pass == 1 {
			// Second pass: Show new
			if y < len(t.newLines) { line = t.newLines[y] }
			if y >= wipeY {
				// Still garbage from prev pass
				line = strings.Repeat("█", width)
			}
		} else {
			// Final
			if y < len(t.newLines) { line = t.newLines[y] }
		}

		c.SetString(0, y, line, green)
	}

	return c.Render()
}

func (t *DataWipe) Done() bool { return t.done }
func (t *DataWipe) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 8. Glitch Spread ---

type GlitchSpread struct {
	BaseTransition
	oldLines []string
	newLines []string
	mask     [][]int // 0=old, 1=glitch, 2=new
	done     bool
}

func NewGlitchSpread() Transition { return &GlitchSpread{} }

func (t *GlitchSpread) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *GlitchSpread) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		// Spread logic would go here, simplified: random spread
		if t.mask == nil { return t, Tick() }

		height := len(t.mask)
		width := len(t.mask[0])

		// Grow new (2) into glitch (1)
		// Grow glitch (1) into old (0)
		for i := 0; i < 50; i++ {
			x := rand.Intn(width)
			y := rand.Intn(height)
			if t.mask[y][x] == 1 {
				t.mask[y][x] = 2

			} else if t.mask[y][x] == 0 {
				if rand.Float64() < 0.1 {
					t.mask[y][x] = 1

				}
			}
		}

		// Aggressive finish
		if rand.Float64() < 0.05 {
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					if t.mask[y][x] < 2 { t.mask[y][x]++;  }
				}
			}
		}

		// Check done
		count := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if t.mask[y][x] == 2 { count++ }
			}
		}
		if count > width*height*9/10 { t.done = true }

		return t, Tick()
	}
	return t, nil
}

func (t *GlitchSpread) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	if t.mask == nil {
		t.mask = make([][]int, height)
		for y := 0; y < height; y++ {
			t.mask[y] = make([]int, width)
			// Seed center
			if y == height/2 { t.mask[y][width/2] = 1 }
		}
	}

	glitchStyle := lipgloss.NewStyle().Background(lipgloss.Color("#FF00FF")).Foreground(lipgloss.Color("#000000"))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			state := t.mask[y][x]
			char := ' '
			style := lipgloss.NewStyle()

			if state == 0 {
				// Old
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
			} else if state == 1 {
				// Glitch
				char = rune(33 + rand.Intn(94))
				style = glitchStyle
			} else {
				// New
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
			}
			c.SetChar(x, y, char, style)
		}
	}
	return c.Render()
}

func (t *GlitchSpread) Done() bool { return t.done }
func (t *GlitchSpread) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 9. Boot Corruption ---

type BootCorruption struct {
	BaseTransition
	newLines []string
	progress float64
	done     bool
}

func NewBootCorruption() Transition { return &BootCorruption{} }
func (t *BootCorruption) SetContent(o, n string) {
	t.newLines = getLines(n)
}

func (t *BootCorruption) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.02
		if t.progress >= 1.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *BootCorruption) View(width, height int) string {
	c := canvas.New(width, height)

	bsod := lipgloss.NewStyle().Background(lipgloss.Color("#0000AA")).Foreground(lipgloss.Color("#FFFFFF"))

	if t.progress < 0.5 {
		// BSOD
		c.Fill(0, 0, width, height, ' ', bsod)
		c.SetString(width/2-10, height/2, "FATAL ERROR", bsod.Bold(true))
		c.SetString(width/2-15, height/2+2, "RECOVERING...", bsod)
	} else {
		// Loading new
		for y := 0; y < height; y++ {
			line := ""
			if y < len(t.newLines) { line = t.newLines[y] }
			c.SetString(0, y, line, lipgloss.NewStyle())
		}
		// Progress bar overlay
		barW := int(t.progress * float64(width))
		c.SetString(0, height-1, strings.Repeat("█", barW), lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))
	}
	return c.Render()
}

func (t *BootCorruption) Done() bool { return t.done }


// --- 10. Hologram Flicker ---

type HologramFlicker struct {
	BaseTransition
	oldLines []string
	newLines []string
	flicker  float64
	done     bool
}

func NewHologramFlicker() Transition { return &HologramFlicker{} }

func (t *HologramFlicker) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *HologramFlicker) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.flicker += 0.05
		if t.flicker > 2.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *HologramFlicker) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	magenta := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))

	// If flicker < 1.0: Mostly old, flickering to new
	// If flicker > 1.0: Mostly new, stabilizing

	threshold := t.flicker / 2.0

	for y := 0; y < height; y++ {
		// Scanline effect
		if y % 2 == 0 { continue }

		for x := 0; x < width; x++ {
			showNew := rand.Float64() < threshold

			char := ' '
			style := lipgloss.NewStyle()

			if showNew {
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
				// RGB Shift
				if rand.Float64() < 0.1 {
					style = cyan
					xOffset := x + 1
					if xOffset < width { c.SetChar(xOffset, y, char, magenta) }
				}
			} else {
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
			}

			c.SetChar(x, y, char, style)
		}
	}
	return c.Render()
}

func (t *HologramFlicker) Done() bool { return t.done }
func (t *HologramFlicker) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

func init() {
	Register(NewCRTPowerOff)
	Register(NewDataWipe)
	Register(NewGlitchSpread)
	Register(NewBootCorruption)
	Register(NewHologramFlicker)
}
