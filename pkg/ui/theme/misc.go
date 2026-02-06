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
}

func NewBBSTheme() Theme { return &BBSTheme{} }
func (t *BBSTheme) Name() string { return "BBS Era" }
func (t *BBSTheme) Description() string { return "14.4k Modem" }

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
	c.DrawBox(5, boxY, width-10, height-boxY-2, magenta)

	c.SetString(7, boxY+2, "Message from SysOp:", white)
	c.SetString(7, boxY+4, q.Text, white)

	c.SetString(7, boxY+7, "Response: "+inputView, cyan)

	c.SetString(width-20, height-2, "NO CARRIER", lipgloss.NewStyle().Foreground(lipgloss.Color("1")))

	return c.Render()
}

// --- 22. Stranger Things Theme ---

type StrangerThingsTheme struct {
	BaseTheme
	litChar rune
}

func NewStrangerThingsTheme() Theme { return &StrangerThingsTheme{} }
func (t *StrangerThingsTheme) Name() string { return "Upside Down" }
func (t *StrangerThingsTheme) Description() string { return "R-U-N" }

func (t *StrangerThingsTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		// Randomly light up a character A-Z
		t.litChar = rune('A' + rand.Intn(26))
		return t, Tick()
	}
	return t, nil
}

func (t *StrangerThingsTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Alphabet wall
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	colors := []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00"}

	startX := (width - 13*4) / 2
	startY := 5

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

	c.SetString(width/2-10, height-5, q.Text, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
	c.SetString(width/2-10, height-3, "> "+inputView, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))

	return c.Render()
}

// --- 23. Blade Runner Theme ---

type BladeRunnerTheme struct {
	BaseTheme
}

func NewBladeRunnerTheme() Theme { return &BladeRunnerTheme{} }
func (t *BladeRunnerTheme) Name() string { return "Voight-Kampff" }
func (t *BladeRunnerTheme) Description() string { return "Enhance 224 to 176" }

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

	c.SetString(30, 8, "SUBJECT: "+q.Text, orange)
	c.SetString(30, 10, "EMOTIONAL RESPONSE: "+inputView, orange)

	c.SetString(width-20, height-2, "VOIGHT-KAMPFF TEST", orange)

	return c.Render()
}

// --- 24. Boot Theme ---

type BootTheme struct {
	BaseTheme
	memCount int
}

func NewBootTheme() Theme { return &BootTheme{} }
func (t *BootTheme) Name() string { return "POST Screen" }
func (t *BootTheme) Description() string { return "BIOS Check" }

func (t *BootTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.memCount += 1024
		if t.memCount > 640000 { t.memCount = 0 }
		return t, Tick()
	}
	return t, nil
}

func (t *BootTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	c.SetString(2, 2, "CTF-BIOS (C) 2024 The Authors", white)
	c.SetString(2, 4, fmt.Sprintf("Memory Test: %d KB OK", t.memCount), white)

	c.SetString(2, 6, "Detecting Primary Master ... Question Found", white)
	c.SetString(2, 7, "Detecting Primary Slave  ... Hint Drive", white)

	c.SetString(2, 10, "Booting from Question Sector...", white)

	c.SetString(2, 12, "KERNEL MSG: "+q.Text, white)
	c.SetString(2, 14, "login: "+inputView, white)

	return c.Render()
}

// --- 25. Ghost in the Shell Theme ---

type GhostInShellTheme struct {
	BaseTheme
}

func NewGhostInShellTheme() Theme { return &GhostInShellTheme{} }
func (t *GhostInShellTheme) Name() string { return "Section 9" }
func (t *GhostInShellTheme) Description() string { return "Stand Alone Complex" }

func (t *GhostInShellTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
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

	c.SetString(centerX-len(q.Text)/2, centerY, q.Text, green.Bold(true))
	c.SetString(centerX-10, centerY+2, "> "+inputView, green)

	return c.Render()
}

func init() {
	Register(NewBBSTheme)
	Register(NewStrangerThingsTheme)
	Register(NewBladeRunnerTheme)
	Register(NewBootTheme)
	Register(NewGhostInShellTheme)
}
