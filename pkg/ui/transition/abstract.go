package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 21. Neural Network ---

type NeuralNetwork struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	done     bool
}

func NewNeuralNetwork() Transition { return &NeuralNetwork{} }

func (t *NeuralNetwork) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *NeuralNetwork) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.05
		if t.progress >= 1.5 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *NeuralNetwork) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	// Draw nodes and connections
	// Simplified: Random points connected by lines

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00AAFF"))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	// If progress < 0.5: draw growing network over Old
	// If progress > 0.5: draw network fading out over New

	lines := t.oldLines
	if t.progress > 0.75 { lines = t.newLines }

	// Background content
	for y := 0; y < height; y++ {
		line := ""
		if y < len(lines) { line = lines[y] }
		c.SetString(0, y, line, lipgloss.NewStyle().Faint(true))
	}

	// Network overlay
	for i := 0; i < 20; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)

		// Draw line (simple Bresenham or just dots)
		// For TUI, just start and end nodes
		c.SetChar(x1, y1, 'O', bright)
		c.SetChar(x2, y2, 'O', bright)

		// "Signal" traveling
		sigX := x1 + int(float64(x2-x1)*t.progress)
		sigY := y1 + int(float64(y2-y1)*t.progress)
		c.SetChar(sigX, sigY, '*', style)
	}

	return c.Render()
}

func (t *NeuralNetwork) Done() bool { return t.done }
func (t *NeuralNetwork) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 22. File Transfer ---

type FileTransfer struct {
	BaseTransition
	newLines []string
	progress float64
	done     bool
}

func NewFileTransfer() Transition { return &FileTransfer{} }
func (t *FileTransfer) SetContent(o, n string) {
	t.newLines = getLines(n)
}

func (t *FileTransfer) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.02
		if t.progress >= 1.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *FileTransfer) View(width, height int) string {
	c := canvas.New(width, height)

	// Title
	c.SetString(width/2-10, height/2-4, "DOWNLOADING CONTENT...", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")))

	// Bar
	barW := width - 20
	filled := int(float64(barW) * t.progress)
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat(" ", barW-filled) + "]"
	c.SetString(10, height/2-2, bar, lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))

	// Stats
	c.SetString(10, height/2, fmt.Sprintf("%d%% Complete", int(t.progress*100)), lipgloss.NewStyle())

	// Preview new content appearing
	previewLines := int(float64(len(t.newLines)) * t.progress)
	for i := 0; i < previewLines; i++ {
		if i < len(t.newLines) {
			c.SetString(0, i, t.newLines[i], lipgloss.NewStyle().Faint(true))
		}
	}

	return c.Render()
}

func (t *FileTransfer) Done() bool { return t.done }


// --- 23. Redaction ---

type Redaction struct {
	BaseTransition
	oldLines []string
	newLines []string
	phase    int // 0=redact old, 1=show new redacted, 2=reveal
	progress float64
	done     bool
}

func NewRedaction() Transition { return &Redaction{} }

func (t *Redaction) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *Redaction) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.05
		if t.progress >= 1.0 {
			t.progress = 0
			t.phase++
			if t.phase > 2 { t.done = true }
		}
		return t, Tick()
	}
	return t, nil
}

func (t *Redaction) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	black := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#000000"))

	var lines []string
	if t.phase == 0 { lines = t.oldLines } else { lines = t.newLines }

	for y, line := range lines {
		if y >= height { break }

		// If phase 0: add redaction bars randomly
		// If phase 1: full redaction bars
		// If phase 2: remove redaction bars

		for x, char := range line {
			if x >= width { break }

			drawRedacted := false
			if t.phase == 0 && rand.Float64() < t.progress { drawRedacted = true }
			if t.phase == 1 { drawRedacted = true }
			if t.phase == 2 && rand.Float64() > t.progress { drawRedacted = true }

			if drawRedacted {
				c.SetChar(x, y, '█', black)
			} else {
				c.SetChar(x, y, char, lipgloss.NewStyle())
			}
		}
	}

	if t.phase == 1 {
		c.SetString(width/2-5, height/2, "TOP SECRET", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Background(lipgloss.Color("#FFFFFF")))
	}

	return c.Render()
}

func (t *Redaction) Done() bool { return t.done }
func (t *Redaction) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 24. Slot Machine ---

type SlotMachine struct {
	BaseTransition
	newLines []string
	offsets  []float64
	speeds   []float64
	done     bool
}

func NewSlotMachine() Transition { return &SlotMachine{} }
func (t *SlotMachine) SetContent(o, n string) {
	t.newLines = getLines(n)
}

func (t *SlotMachine) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if len(t.offsets) == 0 {
			// Init
			t.offsets = make([]float64, 100) // Max lines assumed
			t.speeds = make([]float64, 100)
			for i := range t.speeds {
				t.speeds[i] = 2.0 + rand.Float64()*5.0
			}
		}

		allStopped := true
		for i := range t.speeds {
			if t.speeds[i] > 0 {
				t.offsets[i] += t.speeds[i]
				t.speeds[i] -= 0.1 // Friction
				if t.speeds[i] < 0 { t.speeds[i] = 0 }
				allStopped = false
			} else {
				// Snap to integer?
				// Just keep it 0
			}
		}

		if allStopped { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *SlotMachine) View(width, height int) string {
	c := canvas.New(width, height)
	if len(t.newLines) == 0 { return "" }

	for y := 0; y < height; y++ {
		if y >= len(t.offsets) { break }

		offset := int(t.offsets[y])
		// Draw line, but offset by random chars?
		// Or assume the lines themselves are spinning vertically?
		// Let's spin the lines vertically.

		targetLineIdx := y
		currentLineIdx := (targetLineIdx + offset) % len(t.newLines)

		line := t.newLines[currentLineIdx]
		c.SetString(0, y, line, lipgloss.NewStyle())
	}

	return c.Render()
}

func (t *SlotMachine) Done() bool { return t.done }


// --- 25. Firewall ---

type Firewall struct {
	BaseTransition
	newLines []string
	packets  []int // y positions
	done     bool
}

func NewFirewall() Transition { return &Firewall{} }
func (t *Firewall) SetContent(o, n string) {
	t.newLines = getLines(n)
}

func (t *Firewall) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if rand.Float64() < 0.3 {
			t.packets = append(t.packets, rand.Intn(30))
		}
		if len(t.packets) > 50 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *Firewall) View(width, height int) string {
	c := canvas.New(width, height)

	// Wall
	wallX := width / 2
	c.DrawBox(wallX, 0, 2, height, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")))

	// Packets hitting wall
	for _, y := range t.packets {
		c.SetString(wallX-2, y, "->", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")))
		c.SetString(wallX+2, y, "ALLOW", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))

		// Reveal line
		if y < len(t.newLines) {
			c.SetString(wallX+10, y, t.newLines[y], lipgloss.NewStyle())
		}
	}

	return c.Render()
}

func (t *Firewall) Done() bool { return t.done }

func init() {
	Register(NewNeuralNetwork)
	Register(NewFileTransfer)
	Register(NewRedaction)
	Register(NewSlotMachine)
	Register(NewFirewall)
}
