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

func NewSneakersTheme() Theme                { return &SneakersTheme{} }
func (t *SneakersTheme) Name() string        { return "Sneakers" }
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

	// Box - expanded to fit hint
	boxW := min(60, width-4)
	boxH := 10
	boxX := (width - boxW) / 2
	boxY := (height - boxH) / 2
	if hint != "" {
		boxH = 13
		boxY = (height - boxH) / 2
	}

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

	// Hint with slot machine animation
	if hint != "" {
		displayHint := hint
		frame := t.tick % 30
		if frame < 20 {
			// Slot machine effect - random characters cycling
			runes := []rune(hint)
			displayHint = ""
			for i := range runes {
				if (frame/3+i)%3 == 0 {
					displayHint += string(rune('A' + rand.Intn(26)))
				} else {
					displayHint += string(runes[i])
				}
			}
			if frame < 10 {
				// Still cycling - show more random
				for i := len(hint); i < 15; i++ {
					displayHint += string(rune('A' + rand.Intn(10)))
				}
			}
		}
		c.SetString(boxX+2, boxY+9, "CLUE: "+displayHint, green)
	}

	return c.Render()
}

// --- 17. Hackers (Acid Burn) Theme ---

type HackersTheme struct {
	BaseTheme
	rotation float64
}

func NewHackersTheme() Theme                { return &HackersTheme{} }
func (t *HackersTheme) Name() string        { return "Hackers (1995)" }
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

	// Hint with glitch effect
	if hint != "" {
		frame := int(t.rotation * 10)
		displayHint := hint
		// Glitch: randomly distort some characters
		if frame%8 < 3 {
			runes := []rune(hint)
			for i := range runes {
				if rand.Float64() < 0.3 {
					runes[i] = rune("#%@"[rand.Intn(3)])
				}
			}
			displayHint = string(runes)
		}
		// Offset hint position slightly for glitch effect
		offsetX := 0
		if frame%4 < 1 {
			offsetX = rand.Intn(3) - 1
		}
		c.SetString(centerX-10+offsetX, centerY+4, "HINT: "+displayHint, bgStyle)
	}

	return c.Render()
}

// --- 18. Mr. Robot Theme ---

type MrRobotTheme struct {
	BaseTheme
	tick int
}

func NewMrRobotTheme() Theme                { return &MrRobotTheme{} }
func (t *MrRobotTheme) Name() string        { return "fsociety" }
func (t *MrRobotTheme) Description() string { return "Hello friend." }

func (t *MrRobotTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.tick++
		return t, Tick()
	}
	return t, nil
}

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

	// Hint with typewriter effect
	if hint != "" {
		// Reveal more characters as tick increases
		revealLen := (t.tick / 3) % (len(hint) + 5)
		if revealLen > len(hint) {
			revealLen = len(hint)
		}
		if revealLen > 0 {
			c.SetString(0, 8, "> Hint: "+hint[:revealLen], bg)
		}
	}

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

func NewWargamesTheme() Theme                { return &WargamesTheme{} }
func (t *WargamesTheme) Name() string        { return "WOPR" }
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
	if t.blink {
		cursor = "█"
	}
	c.SetString(22, 14, inputView+cursor, cyan)

	// Hint with blink effect
	if hint != "" {
		if t.blink {
			c.SetString(10, 16, "HINT: "+hint, cyan)
		}
	}

	return c.Render()
}

// --- 20. Crypto Theme ---

type CryptoTheme struct {
	BaseTheme
	hashRate int
	frame    int
}

func NewCryptoTheme() Theme                { return &CryptoTheme{} }
func (t *CryptoTheme) Name() string        { return "Blockchain" }
func (t *CryptoTheme) Description() string { return "HODL" }

func (t *CryptoTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.hashRate = rand.Intn(1000) + 9000
		t.frame++
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

	// Overlay Box - expanded for hint
	boxW := min(70, width-4)
	boxH := 12
	boxX := (width - boxW) / 2
	boxY := (height - boxH) / 2
	if hint != "" {
		boxH = 14
		boxY = (height - boxH) / 2
	}

	c.Fill(boxX, boxY, boxW, boxH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
	c.DrawBox(boxX, boxY, boxW, boxH, gold)

	c.SetString(boxX+2, boxY+1, fmt.Sprintf("MINING BLOCK #%d", time.Now().Unix()/600), gold)
	c.SetString(boxX+2, boxY+2, fmt.Sprintf("HASHRATE: %d MH/s", t.hashRate), gold)

	c.SetString(boxX+2, boxY+5, "CHALLENGE: "+q.Text, green)
	c.SetString(boxX+2, boxY+7, "NONCE: "+inputView, green)

	// Hint: static when it fits, slow marquee only on overflow
	if hint != "" {
		hintPrefix := "CLUE: "
		innerWidth := boxW - 4
		availableHintRunes := innerWidth - len([]rune(hintPrefix))
		if availableHintRunes > 0 {
			displayHint := hint
			hintRunes := len([]rune(hint))

			if hintRunes > availableHintRunes {
				maxOffset := hintRunes - availableHintRunes
				const dwellFrames = 10
				const stepEveryFrames = 4

				sweepFrames := (maxOffset + 1) * stepEveryFrames
				cycleFrames := dwellFrames + sweepFrames + dwellFrames
				phase := 0
				if cycleFrames > 0 {
					phase = t.frame % cycleFrames
				}

				offset := 0
				switch {
				case phase < dwellFrames:
					offset = 0
				case phase < dwellFrames+sweepFrames:
					offset = (phase - dwellFrames) / stepEveryFrames
				default:
					offset = maxOffset
				}

				r := []rune(hint)
				displayHint = string(r[offset : offset+availableHintRunes])
			}

			c.SetString(boxX+2, boxY+9, hintPrefix+displayHint, green)
		}
	}

	return c.Render()
}

func init() {
	Register(NewSneakersTheme)
	Register(NewHackersTheme)
	Register(NewMrRobotTheme)
	Register(NewWargamesTheme)
	Register(NewCryptoTheme)
}
