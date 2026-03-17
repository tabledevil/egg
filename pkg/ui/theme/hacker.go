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
	innerW := boxW - 4

	questionLines := wrapText(q.Text, innerW)
	inputLines := wrapLabeled("> ", inputView, innerW)

	baseH := 10
	if hint != "" {
		baseH = 13
	}
	extraH := 0
	if len(questionLines) > 1 {
		extraH += min(len(questionLines)-1, 3)
	}
	if len(inputLines) > 1 {
		extraH += min(len(inputLines)-1, 2)
	}
	boxH := baseH + extraH
	maxBoxH := height - 2
	if boxH > maxBoxH {
		boxH = maxBoxH
	}
	if boxH < 8 {
		boxH = 8
	}

	contentRows := boxH - 2
	reserved := 1 + 1 + 1 + len(inputLines)
	if hint != "" {
		reserved += 1 + 1
	}
	questionMax := contentRows - reserved
	if questionMax < 1 {
		questionMax = 1
	}
	questionLines = clampLines(questionLines, questionMax, innerW)

	remaining := contentRows - (1 + 1 + len(questionLines) + 1)
	if hint != "" {
		remaining--
	}
	if remaining < 1 {
		remaining = 1
	}
	inputLines = clampLines(inputLines, remaining, innerW)

	boxX := (width - boxW) / 2
	boxY := (height - boxH) / 2

	c.Fill(boxX, boxY, boxW, boxH, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
	c.DrawBox(boxX, boxY, boxW, boxH, green)

	// Content
	c.SetString(boxX+2, boxY+2, "DECRYPTING MESSAGE...", green)

	row := boxY + 4
	for _, line := range questionLines {
		if row >= boxY+boxH-1 {
			break
		}
		c.SetString(boxX+2, row, line, green)
		row++
	}

	if row < boxY+boxH-1 {
		row++
	}

	for _, line := range inputLines {
		if row >= boxY+boxH-1 {
			break
		}
		c.SetString(boxX+2, row, line, green)
		row++
	}

	// Hint with slot machine animation
	if hint != "" {
		hintPrefix := "CLUE: "
		hintWidth := innerW - runeLen(hintPrefix)
		if hintWidth <= 0 {
			return c.Render()
		}

		displayHint := hint
		frame := t.tick % 30
		if frame < 20 {
			// Slot machine effect - random characters cycling
			runes := []rune(hint)
			morphed := make([]rune, 0, len(runes))
			for i := range runes {
				if (frame/3+i)%3 == 0 {
					morphed = append(morphed, rune('A'+rand.Intn(26)))
				} else {
					morphed = append(morphed, runes[i])
				}
			}
			displayHint = string(morphed)
		}

		if runeLen(displayHint) > hintWidth {
			maxOffset := runeLen(displayHint) - hintWidth
			offset := 0
			if maxOffset > 0 {
				offset = (t.tick / 4) % (maxOffset + 1)
			}
			displayHint = sliceRunes(displayHint, offset, hintWidth)
		} else {
			displayHint = truncateToWidth(displayHint, hintWidth)
		}

		hintY := boxY + boxH - 2
		if row+1 < hintY {
			hintY = row + 1
		}
		if hintY < boxY+boxH-1 {
			c.SetString(boxX+2, hintY, hintPrefix+displayHint, green)
		}
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

	panelW := boundedSpan(width, 4, 26, 56)
	panelX := centeredStart(width, panelW)

	questionLines := wrapAndClamp("", q.Text, panelW, 3)
	inputLines := wrapAndClamp("> ", inputView, panelW, 2)

	panelRows := 1 + 1 + len(questionLines) + 1 + len(inputLines)
	if hint != "" {
		panelRows += 2
	}
	row := centeredStart(height, panelRows)

	title := "ACCESS GRANTED"
	c.SetString(panelX+centeredStart(panelW, runeLen(title)), row, title, bgStyle)
	row += 2

	for _, line := range questionLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, bgStyle)
		row++
	}

	if row < height {
		row++
	}

	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, bgStyle)
		row++
	}

	// Hint with glitch effect
	if hint != "" && row < height {
		row++
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

		hintPrefix := "HINT: "
		hintWidth := panelW - runeLen(hintPrefix)
		if hintWidth > 0 {
			if runeLen(displayHint) > hintWidth {
				maxOffset := runeLen(displayHint) - hintWidth
				offset := 0
				if maxOffset > 0 {
					offset = frame % (maxOffset + 1)
				}
				displayHint = sliceRunes(displayHint, offset, hintWidth)
			} else {
				displayHint = truncateToWidth(displayHint, hintWidth)
			}

			// Offset hint position slightly for glitch effect
			offsetX := 0
			if frame%4 < 1 {
				offsetX = rand.Intn(3) - 1
			}
			c.SetString(panelX+offsetX, row, hintPrefix+displayHint, bgStyle)
		}
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
	lineWidth := boundedSpan(width, 0, 12, width)
	c.SetString(0, 0, truncateToWidth(prompt+"./crack_question.sh", lineWidth), bg)

	c.SetString(0, 2, "Running exploit...", bg)

	targetLines := wrapAndClamp("[+] Target acquired: ", q.Text, lineWidth, 3)
	row := 3
	for _, line := range targetLines {
		if row >= height {
			break
		}
		c.SetString(0, row, line, bg)
		row++
	}

	if row < height {
		row++
	}
	if row < height {
		c.SetString(0, row, "Enter payload:", bg)
		row++
	}

	inputLines := wrapAndClamp("> ", inputView, lineWidth, 2)
	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(0, row, line, bg)
		row++
	}

	// Hint with typewriter effect
	if hint != "" && row < height {
		row++
		// Reveal more characters as tick increases
		hintRunes := runeLen(hint)
		revealLen := (t.tick / 3) % (hintRunes + 5)
		if revealLen > hintRunes {
			revealLen = hintRunes
		}
		if revealLen > 0 {
			hintLines := wrapAndClamp("> Hint: ", sliceRunes(hint, 0, revealLen), lineWidth, 2)
			for _, line := range hintLines {
				if row >= height {
					break
				}
				c.SetString(0, row, line, bg)
				row++
			}
		}
	}

	// Random "fsociety" hidden message
	if rand.Float64() < 0.01 && height > 0 {
		mark := "fsociety"
		c.SetString(width-runeLen(mark), height-1, mark, lipgloss.NewStyle().Foreground(lipgloss.Color("#330000")))
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
		return t, SlowTick()
	}
	return t, nil
}

func (t *WargamesTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Standard monochromatic phosphor
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	// Header
	header := "GREETINGS PROFESSOR FALKEN"
	c.SetString(centeredStart(width, runeLen(header)), 2, header, cyan)

	// Simple Menu
	menuW := boundedSpan(width, 2, 20, 64)
	menuX := centeredStart(width, menuW)

	row := 5
	c.SetString(menuX, row, "GAME SELECTION:", cyan)
	row += 2
	c.SetString(menuX, row, "1. FALKEN'S MAZE", cyan)
	row++
	c.SetString(menuX, row, "2. BLACK JACK", cyan)
	row++
	c.SetString(menuX, row, "3. GLOBAL THERMONUCLEAR WAR", cyan)
	row++

	gameLines := wrapAndClamp("4. ", strings.ToUpper(q.Text), menuW, 2)
	for _, line := range gameLines {
		if row >= height {
			break
		}
		c.SetString(menuX, row, line, cyan)
		row++
	}

	if row < height {
		row += 2
	}

	cursor := " "
	if t.blink {
		cursor = "█"
	}

	inputLines := wrapAndClamp("ENTER MOVE: ", inputView+cursor, menuW, 2)
	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(menuX, row, line, cyan)
		row++
	}

	// Hint with blink effect
	if hint != "" && row < height {
		if t.blink {
			hintLines := wrapAndClamp("HINT: ", hint, menuW, 2)
			for _, line := range hintLines {
				if row >= height {
					break
				}
				c.SetString(menuX, row, line, cyan)
				row++
			}
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
	innerW := boxW - 4
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

	challengeLines := clampLines(wrapLabeled("CHALLENGE: ", q.Text, innerW), 2, innerW)
	for i, line := range challengeLines {
		c.SetString(boxX+2, boxY+5+i, line, green)
	}

	nonceLines := clampLines(wrapLabeled("NONCE: ", inputView, innerW), 1, innerW)
	if len(nonceLines) > 0 {
		c.SetString(boxX+2, boxY+7, nonceLines[0], green)
	}

	// Hint: static when it fits, slow marquee only on overflow
	if hint != "" {
		hintPrefix := "CLUE: "
		innerWidth := innerW
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
