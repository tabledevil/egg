package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 6. C64 Theme ---

type C64Theme struct {
	BaseTheme
	cursorFlash bool
}

func NewC64Theme() Theme { return &C64Theme{} }
func (t *C64Theme) Name() string { return "Commodore 64" }
func (t *C64Theme) Description() string { return "64K RAM SYSTEM 38911 BASIC BYTES FREE" }

func (t *C64Theme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if time.Now().UnixNano()/int64(time.Millisecond*500)%2 == 0 {
			t.cursorFlash = true
		} else {
			t.cursorFlash = false
		}
		return t, Tick()
	}
	return t, nil
}

func (t *C64Theme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	blue := lipgloss.NewStyle().Foreground(lipgloss.Color("#70A4B2")).Background(lipgloss.Color("#352879"))

	// Fill background (Dark Blue)
	c.Fill(0, 0, width, height, ' ', blue)

	// Draw Border Area (Light Blue)

	// Inner viewport
	viewportW := width - 4
	viewportH := height - 4
	viewportX := 2
	viewportY := 2

	c.Fill(viewportX, viewportY, viewportW, viewportH, ' ', blue) // Clear center

	// Header
	header := "**** COMMODORE 64 BASIC V2 ****"
	c.SetString((width-len(header))/2, viewportY+1, header, blue)
	c.SetString((width-len("64K RAM SYSTEM"))/2, viewportY+2, "64K RAM SYSTEM  38911 BASIC BYTES FREE", blue)

	// Content
	c.SetString(viewportX+1, viewportY+5, "READY.", blue)
	c.SetString(viewportX+1, viewportY+6, "LOAD \""+strings.ToUpper(q.Text)+"\",8,1", blue)
	c.SetString(viewportX+1, viewportY+7, "SEARCHING FOR "+strings.ToUpper(q.Text), blue)
	c.SetString(viewportX+1, viewportY+8, "LOADING", blue)
	c.SetString(viewportX+1, viewportY+9, "READY.", blue)

	// Input
	cursor := " "
	if t.cursorFlash {
		cursor = "█"
	}
	c.SetString(viewportX+1, viewportY+11, "> "+inputView+cursor, blue)

	if hint != "" {
		c.SetString(viewportX+1, viewportY+13, "SYNTAX ERROR: "+strings.ToUpper(hint), blue)
	}

	return c.Render()
}

// --- 7. DOS Theme ---

type DOSTheme struct {
	BaseTheme
}

func NewDOSTheme() Theme { return &DOSTheme{} }
func (t *DOSTheme) Name() string { return "MS-DOS" }
func (t *DOSTheme) Description() string { return "C:\\>" }

func (t *DOSTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	bg := lipgloss.NewStyle().Background(lipgloss.Color("#0000AA")).Foreground(lipgloss.Color("#AAAAAA")) // Blue background
	hl := lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")).Foreground(lipgloss.Color("#0000AA")) // Selected

	c.Fill(0, 0, width, height, '░', bg)

	// Main Window
	winW := min(60, width-2)
	winH := min(18, height-2)
	winX := (width - winW) / 2
	winY := (height - winH) / 2

	c.Fill(winX, winY, winW, winH, ' ', bg)
	c.DrawBox(winX, winY, winW, winH, bg)

	// Shadow
	c.Fill(winX+1, winY+winH, winW, 1, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
	c.Fill(winX+winW, winY+1, 1, winH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))

	// Title
	title := " BIOS SETUP UTILITY - AWARD SOFTWARE "
	c.SetString(winX+(winW-len(title))/2, winY, title, hl)

	// Content
	c.SetString(winX+2, winY+2, "Question Item:", bg)
	c.SetString(winX+2, winY+3, q.Text, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF55")).Background(lipgloss.Color("#0000AA")))

	// Input Field
	c.SetString(winX+2, winY+6, "Enter Value:", bg)
	c.SetString(winX+15, winY+6, "["+inputView+"_]", hl)

	// Bottom Bar
	bottom := "F1: Help  F10: Save & Exit  Esc: Exit"
	c.SetString(0, height-1, bottom+strings.Repeat(" ", width-len(bottom)), hl)

	if hint != "" {
		c.SetString(winX+2, winY+10, "Hint: "+hint, bg)
	}

	return c.Render()
}

// --- 8. Amiga Theme ---

type AmigaTheme struct {
	BaseTheme
}

func NewAmigaTheme() Theme { return &AmigaTheme{} }
func (t *AmigaTheme) Name() string { return "Amiga Workbench" }
func (t *AmigaTheme) Description() string { return "Guru Meditation" }

func (t *AmigaTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	wbBlue := lipgloss.NewStyle().Background(lipgloss.Color("#0055AA")).Foreground(lipgloss.Color("#FFFFFF"))
	winGrey := lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")).Foreground(lipgloss.Color("#000000"))
	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8800"))

	c.Fill(0, 0, width, height, ' ', wbBlue)

	// Top Bar
	c.Fill(0, 0, width, 1, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF")).Foreground(lipgloss.Color("#0055AA")))
	c.SetString(2, 0, "Workbench 1.3  2563456 graphics mem  0 other mem", lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF")).Foreground(lipgloss.Color("#0055AA")))

	// Window
	winW := min(50, width-10)
	winH := min(12, height-6)
	winX := 5
	winY := 4

	// Shadow
	c.Fill(winX+1, winY+1, winW, winH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))

	// Window Body
	c.Fill(winX, winY, winW, winH, ' ', winGrey)

	// Window Title Bar
	c.Fill(winX, winY, winW, 1, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF")).Foreground(lipgloss.Color("#0055AA")))
	c.SetString(winX+1, winY, "Shell", lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF")).Foreground(lipgloss.Color("#0055AA")))
	c.SetChar(winX, winY, '◻', wbBlue) // Close gadget

	// Content
	c.SetString(winX+1, winY+2, "1.System:> "+q.Text, winGrey)
	c.SetString(winX+1, winY+4, "Answer:> "+inputView+"█", winGrey)

	if hint != "" {
		c.SetString(winX+1, winY+6, "Hint: "+hint, winGrey.Inherit(orange))
	}

	// Gadgets
	c.SetChar(winX+winW-1, winY+winH-1, '◢', winGrey)

	return c.Render()
}

// --- 9. VHS Theme ---

type VHSTheme struct {
	BaseTheme
	frameCount int
}

func NewVHSTheme() Theme { return &VHSTheme{} }
func (t *VHSTheme) Name() string { return "VHS / Analog Horror" }
func (t *VHSTheme) Description() string { return "Tracking Error" }

func (t *VHSTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frameCount++
		return t, Tick()
	}
	return t, nil
}

func (t *VHSTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Render plain text first
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	// Timestamp
	ts := time.Now().Format("PM 03:04:05")
	c.SetString(2, height-2, "PLAY  SP  "+ts, textStyle)

	content := fmt.Sprintf("%s\n\n> %s", q.Text, inputView)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		c.SetString(4, height/3+i, line, textStyle)
	}

	if hint != "" {
		c.SetString(4, height/3+len(lines)+2, hint, lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")))
	}

	// Apply VHS effects
	// 1. Tracking noise (bottom/top bar)
	noiseY := int(math.Sin(float64(t.frameCount)/10.0)*float64(height)/2 + float64(height)/2)
	for x := 0; x < width; x++ {
		c.SetChar(x, noiseY, '░', textStyle)
		if rand.Float64() < 0.3 {
			c.SetChar(x, noiseY+1, '▒', textStyle)
		}
	}

	// 2. Random Static
	for i := 0; i < 20; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		c.SetChar(x, y, '.', lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")))
	}

	return c.Render()
}

// --- 10. Soviet Theme ---

type SovietTheme struct {
	BaseTheme
}

func NewSovietTheme() Theme { return &SovietTheme{} }
func (t *SovietTheme) Name() string { return "Soviet Terminal" }
func (t *SovietTheme) Description() string { return "Top Secret / GRU" }

func (t *SovietTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	bg := lipgloss.NewStyle().Background(lipgloss.Color("#220000"))

	c.Fill(0, 0, width, height, ' ', bg)
	c.DrawBox(0, 0, width, height, red.Inherit(bg))

	// Header
	header := " СОВЕРШЕННО СЕКРЕТНО " // TOP SECRET
	c.SetString((width-len(header))/2, 0, header, red.Inherit(bg).Bold(true))

	// Eagle/Star (Simple ASCII)
	c.SetString(2, 2, "★ СССР ★", red)

	// Content
	c.SetString(4, 5, "OBJECTIVE: "+q.Text, red)
	c.SetString(4, 7, "INPUT DATA: "+inputView, red)

	if hint != "" {
		// Redacted hint
		c.SetString(4, 9, "INTELLIGENCE: "+hint, lipgloss.NewStyle().Foreground(lipgloss.Color("#550000")))
	}

	// Footer
	c.SetString(2, height-2, "AUTHORIZED PERSONNEL ONLY", red)

	return c.Render()
}

func init() {
	Register(NewC64Theme)
	Register(NewDOSTheme)
	Register(NewAmigaTheme)
	Register(NewVHSTheme)
	Register(NewSovietTheme)
}
