package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"math"
	"math/rand"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
)

// --- Helper for parsing view strings ---
var leftoverSGRTextPattern = regexp.MustCompile(`\[[0-9;]*m`)

func getLines(s string) []string {
	clean := xansi.Strip(s)
	clean = leftoverSGRTextPattern.ReplaceAllString(clean, "")
	clean = strings.ReplaceAll(clean, "\r\n", "\n")
	clean = strings.ReplaceAll(clean, "\r", "\n")
	clean = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, clean)
	return strings.Split(clean, "\n")
}

// --- 1. Matrix Wipe ---

type MatrixWipe struct {
	BaseTransition
	oldLines   []string
	newLines   []string
	columns    []wipeColumn
	progress   int
	done       bool
	viewWidth  int
	viewHeight int
}

type wipeColumn struct {
	yOffset float64
	speed   float64
}

func (t *MatrixWipe) initColumns(width int) {
	if width < 1 {
		width = 1
	}

	t.columns = make([]wipeColumn, width)
	for i := range t.columns {
		t.columns[i] = wipeColumn{
			yOffset: -float64(rand.Intn(20)),
			speed:   0.5 + rand.Float64(),
		}
	}
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
			width := t.viewWidth
			if width < 1 {
				width = 80
			}
			t.initColumns(width)
		}

		targetHeight := t.viewHeight
		if targetHeight < 1 {
			targetHeight = 50
		}

		allDone := true
		for i := range t.columns {
			t.columns[i].yOffset += t.columns[i].speed
			if int(t.columns[i].yOffset) < targetHeight+10 {
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
	t.viewWidth = width
	t.viewHeight = height

	if len(t.columns) != width {
		t.initColumns(width)
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	faint := lipgloss.NewStyle().Foreground(lipgloss.Color("#003300"))

	for x, col := range t.columns {
		if x >= width {
			continue
		}
		wipeY := int(col.yOffset)

		for y := 0; y < height; y++ {
			if y < wipeY {
				char := charAt(t.newLines, x, y)
				c.SetChar(x, y, char, lipgloss.NewStyle())
			} else if y < wipeY+10 {
				char := rune(0xFF61 + rand.Intn(20))
				if y == wipeY {
					c.SetChar(x, y, char, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
				} else {
					c.SetChar(x, y, char, green)
				}
			} else {
				char := charAt(t.oldLines, x, y)
				c.SetChar(x, y, char, faint)
			}
		}
	}
	return c.Render()
}

func (t *MatrixWipe) Done() bool { return t.done }
func (t *MatrixWipe) ensureLines(h int) {
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

// --- 2. Pixel Dissolve ---

type PixelDissolve struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	done     bool
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
			noise := float64((x*57+y*13)%100) / 100.0
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
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

// --- 3. Scan Line ---

type ScanLineTransition struct {
	BaseTransition
	oldLines   []string
	newLines   []string
	scanY      int
	done       bool
	viewHeight int
}

func NewScanLineTransition() Transition { return &ScanLineTransition{} }

func (t *ScanLineTransition) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *ScanLineTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.scanY += 2
		targetHeight := t.viewHeight
		if targetHeight < 1 {
			targetHeight = 100
		}
		if t.scanY > targetHeight+1 {
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *ScanLineTransition) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)
	t.viewHeight = height
	bright := lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF"))

	for y := 0; y < height; y++ {
		if y < t.scanY {
			line := ""
			if y < len(t.newLines) {
				line = t.newLines[y]
			}
			c.SetString(0, y, line, lipgloss.NewStyle())
		} else if y == t.scanY {
			c.Fill(0, y, width, 1, ' ', bright)
		} else {
			line := ""
			if y < len(t.oldLines) {
				line = t.oldLines[y]
			}
			c.SetString(0, y, line, lipgloss.NewStyle().Faint(true))
		}
	}
	return c.Render()
}

func (t *ScanLineTransition) Done() bool { return t.done }
func (t *ScanLineTransition) ensureLines(h int) {
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
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
		if t.phase >= 3.14159 {
			t.done = true
		}
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
		effectiveProgress := progress*2 - stagger
		if effectiveProgress < 0 {
			effectiveProgress = 0
		}
		if effectiveProgress > 1 {
			effectiveProgress = 1
		}

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
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

// --- 5. Ink Spread ---

type InkSpread struct {
	BaseTransition
	oldLines  []string
	newLines  []string
	radius    float64
	done      bool
	maxRadius float64
}

func NewInkSpread() Transition { return &InkSpread{} }

func (t *InkSpread) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *InkSpread) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		step := 2.0
		if t.maxRadius > 0 {
			step = t.maxRadius / 45.0
			if step < 1.5 {
				step = 1.5
			}
		}
		t.radius += step

		targetRadius := t.maxRadius
		if targetRadius <= 0 {
			targetRadius = 150
		}
		if t.radius > targetRadius {
			t.done = true
		}
		return t, Tick()
	}
	return t, nil
}

func (t *InkSpread) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)
	centerX := width / 2
	centerY := height / 2

	maxDX := math.Max(float64(centerX), float64(width-centerX))
	maxDY := math.Max(float64(centerY), float64(height-centerY)) * 2.0
	t.maxRadius = math.Sqrt(maxDX*maxDX+maxDY*maxDY) + 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x - centerX)
			dy := float64(y-centerY) * 2.0
			dist := (dx*dx + dy*dy)
			radSq := t.radius * t.radius

			if dist < radSq {
				char := charAt(t.newLines, x, y)
				c.SetChar(x, y, char, lipgloss.NewStyle())
			} else {
				char := charAt(t.oldLines, x, y)
				c.SetChar(x, y, char, lipgloss.NewStyle())
			}
		}
	}
	return c.Render()
}

func (t *InkSpread) Done() bool { return t.done }
func (t *InkSpread) ensureLines(h int) {
	for len(t.oldLines) < h {
		t.oldLines = append(t.oldLines, "")
	}
	for len(t.newLines) < h {
		t.newLines = append(t.newLines, "")
	}
}

func init() {
	Register(NewMatrixWipe)
	Register(NewPixelDissolve)
	Register(NewScanLineTransition)
	Register(NewBlindsTransition)
	Register(NewInkSpread)
}
