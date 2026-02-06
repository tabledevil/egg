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

// --- 11. Fake Decryption ---

type FakeDecryption struct {
	BaseTransition
	oldLines []string
	newLines []string
	progress float64
	done     bool
}

func NewFakeDecryption() Transition { return &FakeDecryption{} }
func (t *FakeDecryption) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *FakeDecryption) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.02
		if t.progress >= 1.2 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *FakeDecryption) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Deterministic "lock" time based on position
			lockTime := float64(x+y*width) / float64(width*height)

			// Show old until progress starts
			// Show random while progress < lockTime
			// Show new when progress > lockTime

			char := ' '
			style := lipgloss.NewStyle()

			if t.progress < 0.1 {
				// Initial delay
				if y < len(t.oldLines) && x < len(t.oldLines[y]) {
					char = rune(t.oldLines[y][x])
				}
			} else if t.progress > lockTime + 0.2 {
				// Locked New
				if y < len(t.newLines) && x < len(t.newLines[y]) {
					char = rune(t.newLines[y][x])
				}
				style = green
			} else {
				// Cycling
				char = rune(33 + rand.Intn(94))
				style = green.Bold(true)
			}

			c.SetChar(x, y, char, style)
		}
	}
	return c.Render()
}

func (t *FakeDecryption) Done() bool { return t.done }
func (t *FakeDecryption) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 12. Terminal Scroll ---

type TerminalScroll struct {
	BaseTransition
	oldLines []string
	newLines []string
	offset   int
	done     bool
}

func NewTerminalScroll() Transition { return &TerminalScroll{} }
func (t *TerminalScroll) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *TerminalScroll) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.offset += 2
		if t.offset > 50 { t.done = true } // Heuristic
		return t, Tick()
	}
	return t, nil
}

func (t *TerminalScroll) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	// Draw Old at -offset
	// Draw New coming up from bottom?
	// Or standard terminal behavior: Old scrolls UP, New appears at bottom line by line.
	// But we have a full screen of "New".
	// Let's assume New scrolls up FROM bottom to replace Old.

	// This is getting complicated to simulate exact terminal history.
	// Simpler: Slide Old UP, slide New UP.

	// Render Old
	for y := 0; y < height; y++ {
		drawY := y - t.offset
		if drawY >= 0 && drawY < height {
			c.SetString(0, drawY, t.oldLines[y], lipgloss.NewStyle())
		}
	}

	// Render New (starting below old)
	startNewY := height - t.offset
	for y := 0; y < height; y++ {
		drawY := startNewY + y
		if drawY >= 0 && drawY < height {
			line := ""
			if y < len(t.newLines) { line = t.newLines[y] }
			c.SetString(0, drawY, line, lipgloss.NewStyle())
		}
	}

	// If new screen has fully scrolled in (offset >= height)
	if t.offset >= height { t.done = true }

	return c.Render()
}

func (t *TerminalScroll) Done() bool { return t.done }
func (t *TerminalScroll) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
}


// --- 13. Typewriter Overtype ---

type TypewriterOvertype struct {
	BaseTransition
	oldLines []string
	newLines []string
	cursorX, cursorY int
	done     bool
}

func NewTypewriterOvertype() Transition { return &TypewriterOvertype{} }
func (t *TypewriterOvertype) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *TypewriterOvertype) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		// Speed typing
		for i := 0; i < 5; i++ {
			t.cursorX++
			if t.cursorX > 80 { // Line width limit
				t.cursorX = 0
				t.cursorY++
			}
		}
		if t.cursorY > 30 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *TypewriterOvertype) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF"))

	for y := 0; y < height; y++ {
		line := ""

		if y < t.cursorY {
			// Fully New
			if y < len(t.newLines) { line = t.newLines[y] }
		} else if y == t.cursorY {
			// Mixed
			newLine := ""
			if y < len(t.newLines) { newLine = t.newLines[y] }
			oldLine := ""
			if y < len(t.oldLines) { oldLine = t.oldLines[y] }

			// Construct mixed line
			res := ""
			maxLen := max(len(newLine), len(oldLine))
			for x := 0; x < maxLen; x++ {
				if x < t.cursorX {
					if x < len(newLine) { res += string(newLine[x]) } else { res += " " }
				} else if x == t.cursorX {
					// Cursor pos
					// handled by setChar overlay usually, but here string building
					res += " "
				} else {
					if x < len(oldLine) { res += string(oldLine[x]) } else { res += " " }
				}
			}
			line = res
		} else {
			// Fully Old
			if y < len(t.oldLines) { line = t.oldLines[y] }
		}

		c.SetString(0, y, line, lipgloss.NewStyle())
	}

	c.SetChar(t.cursorX, t.cursorY, '█', cursorStyle)

	return c.Render()
}

func (t *TypewriterOvertype) Done() bool { return t.done }
func (t *TypewriterOvertype) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}
func max(a, b int) int { if a > b { return a }; return b }


// --- 14. Memory Rewrite ---

type MemoryRewrite struct {
	BaseTransition
	newLines []string
	addr     int
	done     bool
}

func NewMemoryRewrite() Transition { return &MemoryRewrite{} }
func (t *MemoryRewrite) SetContent(o, n string) {
	t.newLines = getLines(n)
}

func (t *MemoryRewrite) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.addr += 4
		if t.addr > 200 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *MemoryRewrite) View(width, height int) string {
	c := canvas.New(width, height)

	hexStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	// Hex Dump View
	for y := 0; y < height; y++ {
		// Fake address
		c.SetString(0, y, fmt.Sprintf("%08X ", y*16+t.addr*16), hexStyle)

		// Content
		if y < t.addr/2 && y < len(t.newLines) {
			// Show actual new content (decoded)
			c.SetString(12, y, t.newLines[y], lipgloss.NewStyle())
		} else {
			// Show hex soup
			soup := ""
			for i := 0; i < 8; i++ { soup += fmt.Sprintf("%02X ", rand.Intn(256)) }
			c.SetString(12, y, soup, hexStyle)
		}
	}
	return c.Render()
}

func (t *MemoryRewrite) Done() bool { return t.done }


// --- 15. Compile Transition ---

type CompileTransition struct {
	BaseTransition
	logs     []string
	progress float64
	done     bool
}

func NewCompileTransition() Transition { return &CompileTransition{} }
func (t *CompileTransition) SetContent(o, n string) {}

func (t *CompileTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.05
		if rand.Float64() < 0.3 {
			t.logs = append(t.logs, fmt.Sprintf("Compiling object_%d.o ...", rand.Intn(999)))
		}
		if t.progress >= 1.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *CompileTransition) View(width, height int) string {
	c := canvas.New(width, height)

	// Draw logs scrolling
	start := 0
	if len(t.logs) > height-2 { start = len(t.logs) - (height - 2) }

	for i, log := range t.logs[start:] {
		c.SetString(0, i, log, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")))
	}

	// Progress Bar
	barW := int(t.progress * float64(width))
	c.SetString(0, height-1, strings.Repeat("=", barW)+">", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))

	return c.Render()
}

func (t *CompileTransition) Done() bool { return t.done }

func init() {
	Register(NewFakeDecryption)
	Register(NewTerminalScroll)
	Register(NewTypewriterOvertype)
	Register(NewMemoryRewrite)
	Register(NewCompileTransition)
}
