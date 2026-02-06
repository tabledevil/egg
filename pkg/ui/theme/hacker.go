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

// --- 16. Sneakers Theme ---

type SneakersTheme struct {
	BaseTheme
	codeStream []rune
	tick       int
}

func NewSneakersTheme() Theme { return &SneakersTheme{} }
func (t *SneakersTheme) Name() string { return "Sneakers" }
func (t *SneakersTheme) Description() string { return "SETEC ASTRONOMY" }

func (t *SneakersTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.tick++
		return t, Tick()
	}
	return t, nil
}

func (t *SneakersTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

	// Background: Random changing characters like the code breaking scene
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float64() < 0.05 {
				char := rune(rand.Intn(26) + 65)
				c.SetChar(x, y, char, lipgloss.NewStyle().Foreground(lipgloss.Color("#004400")))
			}
		}
	}

	// Box
	boxW := min(60, width-4)
	boxH := 10
	boxX := (width - boxW) / 2
	boxY := (height - boxH) / 2

	c.Fill(boxX, boxY, boxW, boxH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
	c.DrawBox(boxX, boxY, boxW, boxH, green)

	// Content
	c.SetString(boxX+2, boxY+2, "DECRYPTING MESSAGE...", green)

	// "Scrambled" question that reveals itself
	// In a real implementation we'd track per-character lock status.
	// Here we simulate it by locking more chars over time based on tick?
	// Actually, just show static text for now but surrounded by scramble
	c.SetString(boxX+2, boxY+4, q.Text, green)

	c.SetString(boxX+2, boxY+7, "> "+inputView, green)

	return c.Render()
}

// --- 17. Hackers (Acid Burn) Theme ---

type HackersTheme struct {
	BaseTheme
	rotation float64
}

func NewHackersTheme() Theme { return &HackersTheme{} }
func (t *HackersTheme) Name() string { return "Hackers (1995)" }
func (t *HackersTheme) Description() string { return "MESS WITH THE BEST, DIE LIKE THE REST" }

func (t *HackersTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.rotation += 0.1
		return t, Tick()
	}
	return t, nil
}

func (t *HackersTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Psychedelic colors
	colors := []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00", "#00FFFF", "#FF00FF"}

	// Tunnel effect
	centerX := width / 2
	centerY := height / 2

	for r := 0; r < min(width, height)/2; r += 2 {
		color := colors[(r/2+int(t.rotation*2))%len(colors)]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		// Draw circle(ish)
		for theta := 0.0; theta < 2*math.Pi; theta += 0.1 {
			x := centerX + int(float64(r)*math.Cos(theta))
			y := centerY + int(float64(r)*math.Sin(theta)/2) // Aspect ratio correction
			c.SetChar(x, y, '*', style)
		}
	}

	// Floating Text
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	c.SetString(centerX-10, centerY-2, "ACCESS GRANTED", bgStyle)
	c.SetString(centerX-len(q.Text)/2, centerY, q.Text, bgStyle)
	c.SetString(centerX-10, centerY+2, "> "+inputView, bgStyle)

	return c.Render()
}

// --- 18. Mr. Robot Theme ---

type MrRobotTheme struct {
	BaseTheme
}

func NewMrRobotTheme() Theme { return &MrRobotTheme{} }
func (t *MrRobotTheme) Name() string { return "fsociety" }
func (t *MrRobotTheme) Description() string { return "Hello friend." }

func (t *MrRobotTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Minimalist, authentic terminal
	bg := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#CCCCCC"))

	c.Fill(0, 0, width, height, ' ', bg)

	// Prompt
	prompt := "root@kali:~# "
	c.SetString(0, 0, prompt+"./crack_question.sh", bg)

	c.SetString(0, 2, "Running exploit...", bg)
	c.SetString(0, 3, "[+] Target acquired: "+q.Text, bg)
	c.SetString(0, 5, "Enter payload:", bg)
	c.SetString(0, 6, "> "+inputView, bg)

	// Random "fsociety" hidden message
	if rand.Float64() < 0.01 {
		c.SetString(width-10, height-1, "fsociety", lipgloss.NewStyle().Foreground(lipgloss.Color("#330000")))
	}

	return c.Render()
}

// --- 19. WarGames Theme ---

type WargamesTheme struct {
	BaseTheme
	blink bool
}

func NewWargamesTheme() Theme { return &WargamesTheme{} }
func (t *WargamesTheme) Name() string { return "WOPR" }
func (t *WargamesTheme) Description() string { return "Shall we play a game?" }

func (t *WargamesTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.blink = !t.blink
		return t, Tick()
	}
	return t, nil
}

func (t *WargamesTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Standard monochromatic phosphor
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	// Header
	c.SetString(width/2-10, 2, "GREETINGS PROFESSOR FALKEN", cyan)

	// Simple Menu
	c.SetString(10, 5, "GAME SELECTION:", cyan)
	c.SetString(10, 7, "1. FALKEN'S MAZE", cyan)
	c.SetString(10, 8, "2. BLACK JACK", cyan)
	c.SetString(10, 9, "3. GLOBAL THERMONUCLEAR WAR", cyan)
	c.SetString(10, 10, "4. "+strings.ToUpper(q.Text), cyan) // The question is the game

	c.SetString(10, 14, "ENTER MOVE:", cyan)
	cursor := " "
	if t.blink { cursor = "█" }
	c.SetString(22, 14, inputView+cursor, cyan)

	return c.Render()
}

// --- 20. Crypto Theme ---

type CryptoTheme struct {
	BaseTheme
	hashRate int
}

func NewCryptoTheme() Theme { return &CryptoTheme{} }
func (t *CryptoTheme) Name() string { return "Blockchain" }
func (t *CryptoTheme) Description() string { return "HODL" }

func (t *CryptoTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.hashRate = rand.Intn(1000) + 9000
		return t, Tick()
	}
	return t, nil
}

func (t *CryptoTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	gold := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))

	// Background hashes
	for y := 0; y < height; y++ {
		hash := fmt.Sprintf("%x", rand.Int63())
		c.SetString(0, y, hash+hash+hash, dim)
	}

	// Overlay Box
	boxW := min(70, width-4)
	boxH := 12
	boxX := (width - boxW) / 2
	boxY := (height - boxH) / 2

	c.Fill(boxX, boxY, boxW, boxH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
	c.DrawBox(boxX, boxY, boxW, boxH, gold)

	c.SetString(boxX+2, boxY+1, fmt.Sprintf("MINING BLOCK #%d", time.Now().Unix()/600), gold)
	c.SetString(boxX+2, boxY+2, fmt.Sprintf("HASHRATE: %d MH/s", t.hashRate), gold)

	c.SetString(boxX+2, boxY+5, "CHALLENGE: "+q.Text, green)
	c.SetString(boxX+2, boxY+7, "NONCE: "+inputView, green)

	return c.Render()
}

func init() {
	Register(NewSneakersTheme)
	Register(NewHackersTheme)
	Register(NewMrRobotTheme)
	Register(NewWargamesTheme)
	Register(NewCryptoTheme)
}
