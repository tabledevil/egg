package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Helper for parsing view strings ---
func getLines(s string) []string {
	return strings.Split(s, "\n")
}

// --- 1. Matrix Wipe ---

type MatrixWipe struct {
	BaseTransition
	oldLines   []string
	newLines   []string
	columns    []wipeColumn
	progress   int
	done       bool
}

type wipeColumn struct {
	yOffset float64
	speed   float64
}

func NewMatrixWipe() Transition { return &MatrixWipe{} }

func (t *MatrixWipe) Init() tea.Cmd {
	return Tick()
}

func (t *MatrixWipe) SetContent(oldView, newView string) {
	t.oldLines = getLines(oldView)
	t.newLines = getLines(newView)
}

func (t *MatrixWipe) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress++

		if len(t.columns) == 0 {
			width := 100
			if len(t.oldLines) > 0 { width = len(t.oldLines[0]) }
			// Ensure reasonable width
			if width < 1 { width = 80 }

			t.columns = make([]wipeColumn, width)
			for i := range t.columns {
				t.columns[i] = wipeColumn{
					yOffset: -float64(rand.Intn(20)),
					speed:   0.5 + rand.Float64(),
				}
			}
		}

		allDone := true
		for i := range t.columns {
			t.columns[i].yOffset += t.columns[i].speed
			if int(t.columns[i].yOffset) < 50 {
				allDone = false
			}
		}

		if allDone {
			t.done = true
		}

		return t, Tick()
	}
	return t, nil
}

func (t *MatrixWipe) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	faint := lipgloss.NewStyle().Foreground(lipgloss.Color("#003300"))

	for x, col := range t.columns {
		if x >= width { continue }
		wipeY := int(col.yOffset)

		for y := 0; y < height; y++ {
			if y < wipeY {
				char := ' '
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
				c.SetChar(x, y, char, lipgloss.NewStyle())
			} else if y < wipeY + 10 {
				char := rune(0xFF61 + rand.Intn(20))
				if y == wipeY {
					c.SetChar(x, y, char, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
				} else {
					c.SetChar(x, y, char, green)
				}
			} else {
				char := ' '
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
				c.SetChar(x, y, char, faint)
			}
		}
	}
	return c.Render()
}

func (t *MatrixWipe) Done() bool { return t.done }
func (t *MatrixWipe) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

// --- 2. Pixel Dissolve ---

type PixelDissolve struct {
	BaseTransition
	oldLines   []string
	newLines   []string
	progress   float64
	done       bool
}

func NewPixelDissolve() Transition { return &PixelDissolve{} }

func (t *PixelDissolve) SetContent(oldView, newView string) {
	t.oldLines = getLines(oldView)
	t.newLines = getLines(newView)
}

func (t *PixelDissolve) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.02
		if t.progress >= 1.2 {
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *PixelDissolve) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			noise := float64((x * 57 + y * 13) % 100) / 100.0
			char := ' '
			if noise < t.progress {
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
			} else {
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
			}
			c.SetChar(x, y, char, lipgloss.NewStyle())
		}
	}
	return c.Render()
}

func (t *PixelDissolve) Done() bool { return t.done }
func (t *PixelDissolve) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

// --- 3. Scan Line ---

type ScanLineTransition struct {
	BaseTransition
	oldLines []string
	newLines []string
	scanY    int
	done     bool
}

func NewScanLineTransition() Transition { return &ScanLineTransition{} }

func (t *ScanLineTransition) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *ScanLineTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.scanY += 2
		if t.scanY > 100 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *ScanLineTransition) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)
	bright := lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF"))

	for y := 0; y < height; y++ {
		if y < t.scanY {
			line := ""
			if y < len(t.newLines) { line = t.newLines[y] }
			c.SetString(0, y, line, lipgloss.NewStyle())
		} else if y == t.scanY {
			c.Fill(0, y, width, 1, ' ', bright)
		} else {
			line := ""
			if y < len(t.oldLines) { line = t.oldLines[y] }
			c.SetString(0, y, line, lipgloss.NewStyle().Faint(true))
		}
	}
	return c.Render()
}

func (t *ScanLineTransition) Done() bool { return t.done }
func (t *ScanLineTransition) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

// --- 4. Blinds Transition ---

type BlindsTransition struct {
	BaseTransition
	oldLines []string
	newLines []string
	phase    float64
	done     bool
}

func NewBlindsTransition() Transition { return &BlindsTransition{} }

func (t *BlindsTransition) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *BlindsTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.phase += 0.1
		if t.phase >= 3.14159 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *BlindsTransition) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)
	blindH := 4

	for y := 0; y < height; y++ {
		blindIndex := y / blindH
		progress := t.phase / 3.14159
		stagger := float64(blindIndex) * 0.1
		effectiveProgress := progress * 2 - stagger
		if effectiveProgress < 0 { effectiveProgress = 0 }
		if effectiveProgress > 1 { effectiveProgress = 1 }

		if effectiveProgress > 0.5 {
			if y < len(t.newLines) {
				c.SetString(0, y, t.newLines[y], lipgloss.NewStyle())
			}
		} else {
			if y < len(t.oldLines) {
				c.SetString(0, y, t.oldLines[y], lipgloss.NewStyle())
			}
		}
	}
	return c.Render()
}

func (t *BlindsTransition) Done() bool { return t.done }
func (t *BlindsTransition) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

// --- 5. Ink Spread ---

type InkSpread struct {
	BaseTransition
	oldLines []string
	newLines []string
	radius   float64
	done     bool
}

func NewInkSpread() Transition { return &InkSpread{} }

func (t *InkSpread) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *InkSpread) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.radius += 2.0
		if t.radius > 150 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *InkSpread) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)
	centerX := width / 2
	centerY := height / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY) * 2.0
			dist := (dx*dx + dy*dy)
			radSq := t.radius * t.radius

			if dist < radSq {
				char := ' '
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
				c.SetChar(x, y, char, lipgloss.NewStyle())
			} else {
				char := ' '
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
				c.SetChar(x, y, char, lipgloss.NewStyle())
			}
		}
	}
	return c.Render()
}

func (t *InkSpread) Done() bool { return t.done }
func (t *InkSpread) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

func init() {
	Register(NewMatrixWipe)
	Register(NewPixelDissolve)
	Register(NewScanLineTransition)
	Register(NewBlindsTransition)
	Register(NewInkSpread)
}
