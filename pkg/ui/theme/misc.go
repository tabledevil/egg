package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"fmt"
	"math"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 21. BBS Theme ---

type BBSTheme struct {
	BaseTheme
	frame int
}

func NewBBSTheme() Theme                { return &BBSTheme{} }
func (t *BBSTheme) Name() string        { return "BBS Era" }
func (t *BBSTheme) Description() string { return "14.4k Modem" }

func (t *BBSTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, Tick()
	}
	return t, nil
}

func (t *BBSTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// ANSI Colors
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	magenta := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// ASCII Art Banner
	banner := []string{
		" ▄▄▄       ███▄    █   ██████  ██▓",
		"▒████▄     ██ ▀█   █ ▒██    ▒ ▓██▒",
		"▒██  ▀█▄  ▓██  ▀█ ██▒░ ▓██▄   ▒██▒",
		"░██▄▄▄▄██ ▓██▒  ▐▌██▒  ▒   ██▒░██░",
		" ▓█   ▓██▒▒██░   ▓██░▒██████▒▒░██░",
		" ▒▒   ▓▒█░░ ▒░   ▒ ▒ ▒ ▒▓▒ ▒ ░░▓  ",
	}

	y := 2
	for _, line := range banner {
		c.SetString((width-len(line))/2, y, line, cyan)
		y++
	}

	// Box
	boxY := y + 2
	boxH := height - boxY - 4
	c.DrawBox(5, boxY, width-10, boxH, magenta)
	innerW := width - 14
	if innerW < 10 {
		innerW = 10
	}

	c.SetString(7, boxY+2, "Message from SysOp:", white)

	row := boxY + 4
	questionLines := clampLines(wrapText(q.Text, innerW), 3, innerW)
	for _, line := range questionLines {
		if row >= boxY+boxH-1 {
			break
		}
		c.SetString(7, row, line, white)
		row++
	}

	if row < boxY+boxH-1 {
		row++
	}

	responseLines := clampLines(wrapLabeled("Response: ", inputView, innerW), 2, innerW)
	for _, line := range responseLines {
		if row >= boxY+boxH-1 {
			break
		}
		c.SetString(7, row, line, cyan)
		row++
	}

	// Hint with marquee scrolling (modem-style)
	if hint != "" {
		hintPrefix := "HINT: "
		hintWidth := innerW - runeLen(hintPrefix)
		if hintWidth > 0 {
			displayHint := hint
			if runeLen(hint) > hintWidth {
				maxOffset := runeLen(hint) - hintWidth
				offset := (t.frame / 3) % (maxOffset + 1)
				displayHint = sliceRunes(hint, offset, hintWidth)
			} else {
				displayHint = truncateToWidth(hint, hintWidth)
			}

			hintY := row + 1
			if hintY >= boxY+boxH-1 {
				hintY = row
			}
			if hintY < boxY+boxH-1 {
				c.SetString(7, hintY, hintPrefix+displayHint, white)
			}
		}
	}

	c.SetString(width-20, height-2, "NO CARRIER", lipgloss.NewStyle().Foreground(lipgloss.Color("1")))

	return c.Render()
}

// --- 22. Stranger Things Theme ---

type StrangerThingsTheme struct {
	BaseTheme
	litChar rune
	frame   int
}

func NewStrangerThingsTheme() Theme                { return &StrangerThingsTheme{} }
func (t *StrangerThingsTheme) Name() string        { return "Upside Down" }
func (t *StrangerThingsTheme) Description() string { return "R-U-N" }

func (t *StrangerThingsTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		// Randomly light up a character A-Z
		t.litChar = rune('A' + rand.Intn(26))
		t.frame++
		return t, Tick()
	}
	return t, nil
}

func (t *StrangerThingsTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Alphabet wall
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	colors := []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00"}

	gridW := ((8 - 1) * 5) + 1
	startX := centeredStart(width, gridW)
	startY := 4
	if startY < 2 {
		startY = 2
	}

	for i, char := range letters {
		x := startX + (i%8)*5
		y := startY + (i/8)*3

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i%len(colors)]))
		if char == t.litChar {
			style = style.Bold(true).Background(lipgloss.Color("#FFFFFF"))
		} else {
			style = style.Faint(true)
		}

		c.SetChar(x, y, char, style)
		c.SetChar(x, y-1, '●', style) // Bulb
	}

	panelW := boundedSpan(width, 4, 20, 56)
	panelX := centeredStart(width, panelW)
	row := height - 7
	if row < startY+10 {
		row = startY + 10
	}

	questionLines := wrapAndClamp("", q.Text, panelW, 2)
	for _, line := range questionLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
		row++
	}

	if row < height {
		row++
	}

	inputLines := wrapAndClamp("> ", inputView, panelW, 1)
	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
		row++
	}

	// Hint with blink effect - hint pulses with alphabet
	if hint != "" && t.frame%10 < 5 && row < height {
		hintLines := wrapAndClamp("HINT: ", hint, panelW, 2)
		for _, line := range hintLines {
			if row >= height {
				break
			}
			c.SetString(panelX, row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true))
			row++
		}
	}

	return c.Render()
}

// --- 23. Blade Runner Theme ---

type BladeRunnerTheme struct {
	BaseTheme
	frame int
}

func NewBladeRunnerTheme() Theme                { return &BladeRunnerTheme{} }
func (t *BladeRunnerTheme) Name() string        { return "Voight-Kampff" }
func (t *BladeRunnerTheme) Description() string { return "Enhance 224 to 176" }

func (t *BladeRunnerTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, Tick()
	}
	return t, nil
}

func (t *BladeRunnerTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))

	// Grid overlay
	for y := 0; y < height; y += 4 {
		c.SetString(0, y, strings.Repeat("-", width), lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")))
	}
	for x := 0; x < width; x += 10 {
		for y := 0; y < height; y++ {
			c.SetChar(x, y, '|', lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")))
		}
	}

	// Eye image placeholder (ASCII)
	eye := []string{
		"      /\\      ",
		"    /    \\    ",
		"  /   O    \\  ",
		"  \\        /  ",
		"    \\    /    ",
		"      \\/      ",
	}

	for i, line := range eye {
		c.SetString(5, 5+i, line, orange)
	}

	contentX := 30
	contentY := 8
	if width < 70 {
		contentX = 2
		contentY = 13
	}
	contentW := width - contentX - 2
	if contentW < 16 {
		contentW = boundedSpan(width, 2, 16, width-4)
		contentX = centeredStart(width, contentW)
	}

	row := contentY
	subjectLines := wrapAndClamp("SUBJECT: ", q.Text, contentW, 2)
	for _, line := range subjectLines {
		if row >= height {
			break
		}
		c.SetString(contentX, row, line, orange)
		row++
	}

	if row < height {
		row++
	}
	responseLines := wrapAndClamp("EMOTIONAL RESPONSE: ", inputView, contentW, 2)
	for _, line := range responseLines {
		if row >= height {
			break
		}
		c.SetString(contentX, row, line, orange)
		row++
	}

	// Hint with blink effect - orange pulsing hint
	if hint != "" && t.frame%12 < 8 && row < height {
		if row < height {
			row++
		}
		hintLines := wrapAndClamp("ANALYSIS: ", hint, contentW, 2)
		for _, line := range hintLines {
			if row >= height {
				break
			}
			c.SetString(contentX, row, line, orange)
			row++
		}
	}

	footer := "VOIGHT-KAMPFF TEST"
	c.SetString(width-runeLen(footer)-2, height-2, truncateToWidth(footer, boundedSpan(width, 2, 8, width-4)), orange)

	return c.Render()
}

// --- 24. Boot Theme ---

type BootTheme struct {
	BaseTheme
	memCount int
}

func NewBootTheme() Theme                { return &BootTheme{} }
func (t *BootTheme) Name() string        { return "POST Screen" }
func (t *BootTheme) Description() string { return "BIOS Check" }

func (t *BootTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.memCount += 1024
		if t.memCount > 640000 {
			t.memCount = 0
		}
		return t, Tick()
	}
	return t, nil
}

func (t *BootTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	lineW := boundedSpan(width, 2, 20, 78)
	lineX := 2
	if lineX+lineW >= width {
		lineX = centeredStart(width, lineW)
	}

	row := 2
	c.SetString(lineX, row, truncateToWidth("CTF-BIOS (C) 2024 The Authors", lineW), white)
	row += 2
	c.SetString(lineX, row, truncateToWidth(fmt.Sprintf("Memory Test: %d KB OK", t.memCount), lineW), white)

	row += 2
	c.SetString(lineX, row, truncateToWidth("Detecting Primary Master ... Question Found", lineW), white)
	row++
	c.SetString(lineX, row, truncateToWidth("Detecting Primary Slave  ... Hint Drive", lineW), white)

	row += 2
	c.SetString(lineX, row, truncateToWidth("Booting from Question Sector...", lineW), white)

	row += 2
	kernelLines := wrapAndClamp("KERNEL MSG: ", q.Text, lineW, 3)
	for _, line := range kernelLines {
		if row >= height {
			break
		}
		c.SetString(lineX, row, line, white)
		row++
	}

	if row < height {
		row++
	}
	loginLines := wrapAndClamp("login: ", inputView, lineW, 2)
	for _, line := range loginLines {
		if row >= height {
			break
		}
		c.SetString(lineX, row, line, white)
		row++
	}

	// Hint with marquee scrolling (BIOS ticker style)
	if hint != "" && row < height {
		frame := t.memCount / 1024
		hintPrefix := "HINT: "
		hintWidth := lineW - runeLen(hintPrefix)
		if hintWidth > 0 {
			displayHint := hint
			if runeLen(hint) > hintWidth {
				maxOffset := runeLen(hint) - hintWidth
				offset := 0
				if maxOffset > 0 {
					offset = frame % (maxOffset + 1)
				}
				displayHint = sliceRunes(hint, offset, hintWidth)
			} else {
				displayHint = truncateToWidth(hint, hintWidth)
			}
			hintY := row + 1
			if hintY >= height {
				hintY = row
			}
			c.SetString(lineX, hintY, hintPrefix+displayHint, white)
		}
	}

	return c.Render()
}

// --- 25. Ghost in the Shell Theme ---

type GhostInShellTheme struct {
	BaseTheme
	frame int
}

func NewGhostInShellTheme() Theme                { return &GhostInShellTheme{} }
func (t *GhostInShellTheme) Name() string        { return "Section 9" }
func (t *GhostInShellTheme) Description() string { return "Stand Alone Complex" }

func (t *GhostInShellTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frame++
		return t, Tick()
	}
	return t, nil
}

func (t *GhostInShellTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	c := canvas.New(width, height)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	// Circular layout text
	centerX := width / 2
	centerY := height / 2
	radius := 10.0

	for i := 0; i < 12; i++ {
		angle := float64(i) * (3.14159 / 6.0)
		x := centerX + int(radius*2*float64(width)/float64(height)*0.5*math.Cos(angle)) // Aspect corrected
		y := centerY + int(radius*math.Sin(angle))

		c.SetString(x, y, "DATA", green)
	}

	panelW := boundedSpan(width, 6, 24, 60)
	panelX := centeredStart(width, panelW)
	row := centerY

	questionLines := wrapAndClamp("", q.Text, panelW, 2)
	for _, line := range questionLines {
		if row >= height {
			break
		}
		c.SetString(panelX+centeredStart(panelW, runeLen(line)), row, line, green.Bold(true))
		row++
	}

	if row < height {
		row++
	}
	inputLines := wrapAndClamp("> ", inputView, panelW, 2)
	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, green)
		row++
	}

	// Hint with typewriter effect
	if hint != "" && row < height {
		hintRunes := runeLen(hint)
		revealLen := (t.frame / 3) % (hintRunes + 3)
		if revealLen > hintRunes {
			revealLen = hintRunes
		}
		if revealLen > 0 {
			hintLines := wrapAndClamp("HINT: ", sliceRunes(hint, 0, revealLen), panelW, 2)
			for _, line := range hintLines {
				if row >= height {
					break
				}
				c.SetString(panelX, row, line, green)
				row++
			}
		}
	}

	return c.Render()
}

func init() {
	Register(NewBBSTheme)
	Register(NewStrangerThingsTheme)
	Register(NewBladeRunnerTheme)
	Register(NewBootTheme)
	Register(NewGhostInShellTheme)
}
