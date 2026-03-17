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

// --- 1. Matrix Theme ---

type MatrixTheme struct {
	BaseTheme
	columns     []matrixColumn
	initialized bool
}

type matrixColumn struct {
	x, y   float64
	speed  float64
	length int
	chars  []rune
}

func NewMatrixTheme() Theme {
	return &MatrixTheme{}
}

func (t *MatrixTheme) Name() string        { return "Matrix" }
func (t *MatrixTheme) Description() string { return "The Digital Rain" }

func (t *MatrixTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		// Update columns
		for i := range t.columns {
			t.columns[i].y += t.columns[i].speed
			// Reset if fell off screen
			if t.columns[i].y-float64(t.columns[i].length) > 50 { // approximate max height check
				t.columns[i].y = float64(-rand.Intn(20))
				t.columns[i].speed = 0.5 + rand.Float64()
			}
			// Random character flip
			if rand.Float64() < 0.05 {
				idx := rand.Intn(len(t.columns[i].chars))
				t.columns[i].chars[idx] = randomKatakana()
			}
		}
		return t, Tick()
	}
	return t, nil
}

func (t *MatrixTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Initialize columns if resize or first run
	if !t.initialized || len(t.columns) != width {
		t.columns = make([]matrixColumn, width)
		for x := 0; x < width; x++ {
			length := 10 + rand.Intn(20)
			chars := make([]rune, length)
			for i := range chars {
				chars[i] = randomKatakana()
			}
			t.columns[x] = matrixColumn{
				x:      float64(x),
				y:      float64(rand.Intn(height) - height), // Start above
				speed:  0.2 + rand.Float64()*0.8,
				length: length,
				chars:  chars,
			}
		}
		t.initialized = true
	}

	// Draw Background (Rain)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#003300"))
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	for x, col := range t.columns {
		headY := int(col.y)
		for i := 0; i < col.length; i++ {
			y := headY - i
			if y >= 0 && y < height {
				char := col.chars[i]
				var style lipgloss.Style
				if i == 0 {
					style = whiteStyle
				} else if i < 5 {
					style = greenStyle
				} else {
					style = dimStyle
				}
				c.SetChar(x, y, char, style)
			}
		}
	}

	// Draw UI Box
	boxWidth := min(60, width-4)
	boxHeight := min(15, height-4)
	boxX := (width - boxWidth) / 2
	boxY := (height - boxHeight) / 2

	// Clear area for box
	for y := boxY; y < boxY+boxHeight; y++ {
		for x := boxX; x < boxX+boxWidth; x++ {
			c.SetChar(x, y, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#001100")))
		}
	}

	// Draw Text
	textWidth := boxWidth - 4
	var lines []string
	lines = append(lines, "WAKE UP NEO...", "")
	lines = append(lines, wrapText(q.Text, textWidth)...)
	lines = append(lines, "")
	lines = append(lines, wrapLabeled("> ", inputView, textWidth)...)
	if hint != "" {
		lines = append(lines, "")
		lines = append(lines, wrapLabeled("HINT: ", hint, textWidth)...)
	}
	lines = clampLines(lines, boxHeight-4, textWidth)

	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#001100")).Bold(true)

	currentY := boxY + 2
	for _, line := range lines {
		if currentY >= boxY+boxHeight-1 {
			break
		}
		c.SetString(boxX+2, currentY, line, textStyle)
		currentY++
	}

	// Draw Border
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#001100"))
	drawBox(c, boxX, boxY, boxWidth, boxHeight, borderStyle)

	return c.Render()
}

func randomKatakana() rune {
	// Half-width katakana
	return rune(0xFF61 + rand.Intn(0xFF9F-0xFF61))
}

// --- 2. Cyberpunk Theme ---

type CyberpunkTheme struct {
	BaseTheme
	glitchIntensity float64
	frameCount      int
}

func NewCyberpunkTheme() Theme {
	return &CyberpunkTheme{}
}

func (t *CyberpunkTheme) Name() string        { return "Cyberpunk 2077" }
func (t *CyberpunkTheme) Description() string { return "Night City Glitch" }

func (t *CyberpunkTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.frameCount++
		if rand.Float64() < 0.1 {
			t.glitchIntensity = rand.Float64() * 0.5
		} else {
			t.glitchIntensity *= 0.8
		}
		return t, Tick()
	}
	return t, nil
}

func (t *CyberpunkTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Styles
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#FCEE0A")).Bold(true) // Cyberpunk Yellow
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	pink := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	bg := lipgloss.NewStyle().Background(lipgloss.Color("#050505"))

	// Background grid logic
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%20 == 0 || rand.Float64() < 0.005 {
				c.SetChar(x, y, '+', lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")))
			}
		}
	}

	// Main content with glitch offset
	offsetX := 0
	if t.glitchIntensity > 0.2 {
		offsetX = rand.Intn(3) - 1
	}

	panelW := boundedSpan(width, 4, 24, 62)
	panelX := 5
	if panelX+panelW >= width {
		panelX = centeredStart(width, panelW)
	}
	row := height / 3
	if row < 2 {
		row = 2
	}

	title := "NET_BUNKER // ACCESS_REQUIRED"
	c.SetString(panelX+offsetX, row, truncateToWidth(title, panelW), yellow.Inherit(bg))
	row += 2

	questionLines := wrapAndClamp("", q.Text, panelW, 3)
	for _, line := range questionLines {
		if row >= height {
			break
		}
		c.SetString(panelX, row, line, yellow.Inherit(bg))
		row++
	}

	if row < height {
		row++
	}

	inputLines := wrapAndClamp("> ", inputView, panelW, 1)
	if len(inputLines) > 0 && row < height {
		line := inputLines[0]
		// Chromatic Aberration for Input
		c.SetString(panelX+2, row, line, cyan) // Shift right
		c.SetString(panelX-2, row, line, pink) // Shift left
		c.SetString(panelX, row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true))
		row++
	}

	if hint != "" && row < height {
		row++
		hintLines := wrapAndClamp("WARNING: ", hint, panelW, 2)
		for _, line := range hintLines {
			if row >= height {
				break
			}
			c.SetString(panelX, row, line, pink)
			row++
		}
	}

	// Decor
	conn := "CONN: SECURE"
	connX := width - runeLen(conn) - 2
	if connX < 2 {
		connX = 2
	}
	c.SetString(connX, 2, conn, cyan)

	ping := fmt.Sprintf("PING: %dms", rand.Intn(20)+10)
	c.SetString(connX, 3, truncateToWidth(ping, boundedSpan(width-connX, 0, 6, width-connX)), cyan)

	return c.Render()
}

// --- 3. Tron Theme ---

type TronTheme struct {
	BaseTheme
	gridOffset float64
}

func NewTronTheme() Theme                { return &TronTheme{} }
func (t *TronTheme) Name() string        { return "The Grid" }
func (t *TronTheme) Description() string { return "Digital Frontier" }

func (t *TronTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.gridOffset += 0.5
		if t.gridOffset >= 10 {
			t.gridOffset = 0
		}
		return t, Tick()
	}
	return t, nil
}

func (t *TronTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	glowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")) // Cyan
	darkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#004444"))

	// Horizon line
	horizonY := height / 3
	c.SetString(0, horizonY, strings.Repeat("━", width), glowStyle)

	// Perspective Grid (Fake)
	centerX := width / 2

	// Vertical lines converging
	for x := -width; x < width*2; x += 10 {
		// Draw line from bottom to horizon
		// Simple slope: (x - centerX) * scale
		for y := height - 1; y > horizonY; y-- {
			distFromHorizon := float64(y - horizonY)
			factor := distFromHorizon / float64(height-horizonY)

			// Spread x based on factor (closer = wider spread)
			// Wait, perspective: lines converge at horizon (factor 0)
			// So at bottom (factor 1), x is original.
			// At top (factor 0), x is centerX.

			screenX := centerX + int(float64(x-centerX)*factor)
			if screenX >= 0 && screenX < width {
				c.SetChar(screenX, y, '│', darkStyle)
			}
		}
	}

	// Horizontal lines moving
	// Logarithmic spacing
	for i := 0; i < 20; i++ {
		offset := (float64(i) * 2) + (t.gridOffset / 5) // Moving
		y := horizonY + int(math.Pow(1.5, offset))
		if y < height {
			c.SetString(0, y, strings.Repeat("─", width), darkStyle)
		}
	}

	// Text Floating
	textW := boundedSpan(width, 6, 24, 60)
	textX := centeredStart(width, textW)
	textY := horizonY - 5
	if textY < 1 {
		textY = 1
	}

	row := textY
	questionLines := wrapAndClamp("", q.Text, textW, 3)
	for _, line := range questionLines {
		if row >= height {
			break
		}
		c.SetString(textX+centeredStart(textW, runeLen(line)), row, line, glowStyle.Bold(true))
		row++
	}

	if row < height {
		row++
	}
	inputLines := wrapAndClamp("> ", inputView, textW, 2)
	for _, line := range inputLines {
		if row >= height {
			break
		}
		c.SetString(textX, row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")))
		row++
	}

	if hint != "" && row < height {
		if row < height {
			row++
		}
		hintLines := wrapAndClamp("", hint, textW, 2)
		for _, line := range hintLines {
			if row >= height {
				break
			}
			c.SetString(textX+centeredStart(textW, runeLen(line)), row, line, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9900")))
			row++
		}
	}

	return c.Render()
}

// --- 4. Alien (Nostromo) Theme ---

type AlienTheme struct {
	BaseTheme
	cursorBlink bool
}

func NewAlienTheme() Theme                { return &AlienTheme{} }
func (t *AlienTheme) Name() string        { return "Nostromo" }
func (t *AlienTheme) Description() string { return "MU-TH-UR 6000" }

func (t *AlienTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.cursorBlink = !t.cursorBlink
		return t, SlowTick()
	}
	return t, nil
}

func (t *AlienTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	amber := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB000")) // Classic Amber
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#553300"))

	// Header
	c.SetString(2, 1, "INTERFACE 2037", dim)
	c.SetString(width-20, 1, "WEYLAND-YUTANI", dim)

	// Frame
	drawBox(c, 2, 3, width-4, height-5, amber)

	// Corners
	c.SetChar(2, 3, '▛', amber)
	c.SetChar(width-3, 3, '▜', amber)
	c.SetChar(2, height-3, '▙', amber)
	c.SetChar(width-3, height-3, '▟', amber)

	// Content
	innerW := width - 10
	if innerW < 12 {
		innerW = 12
	}

	c.SetString(5, 6, "PRIORITY ONE:", amber)

	row := 8
	queryLines := clampLines(wrapLabeled("QUERY: ", strings.ToUpper(q.Text), innerW), 2, innerW)
	for _, line := range queryLines {
		if row >= height-3 {
			break
		}
		c.SetString(5, row, line, amber)
		row++
	}

	cursor := " "
	if t.cursorBlink {
		cursor = "█"
	}

	if row < height-3 {
		row++
	}
	inputLines := clampLines(wrapLabeled("INPUT: ", inputView+cursor, innerW), 2, innerW)
	for _, line := range inputLines {
		if row >= height-3 {
			break
		}
		c.SetString(5, row, line, amber)
		row++
	}

	if hint != "" {
		if row < height-3 {
			row++
		}
		hintLines := clampLines(wrapLabeled("ANALYSIS: ", strings.ToUpper(hint), innerW), 2, innerW)
		for _, line := range hintLines {
			if row >= height-3 {
				break
			}
			c.SetString(5, row, line, dim)
			row++
		}
	}

	return c.Render()
}

// --- 5. System Shock (SHODAN) Theme ---

type SystemShockTheme struct {
	BaseTheme
	corruptionLevel float64
}

func NewSystemShockTheme() Theme                { return &SystemShockTheme{} }
func (t *SystemShockTheme) Name() string        { return "Citadel" }
func (t *SystemShockTheme) Description() string { return "Look at you, Hacker" }

func (t *SystemShockTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.corruptionLevel = math.Sin(float64(time.Now().UnixNano())/1e9)*0.5 + 0.5
		return t, Tick()
	}
	return t, nil
}

func (t *SystemShockTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	c := canvas.New(width, height)

	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	grey := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))

	// Random background noise
	for i := 0; i < 50; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		c.SetChar(x, y, rune(rand.Intn(30)+33), grey)
	}

	// Big SHODAN face approximation (ASCII) - simplistic
	eyeX := width / 2
	eyeY := 5
	c.SetString(eyeX-5, eyeY, "[  O  ]", red.Bold(true))

	// Glitch blocks (draw before the message so we don't overwrite the question/input).
	for i := 0; i < 5; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		c.SetString(x, y, "ERR", lipgloss.NewStyle().Background(lipgloss.Color("#FF0000")).Foreground(lipgloss.Color("#000000")))
	}

	// Message
	msg := q.Text
	if rand.Float64() < 0.1 {
		msg = "INSECT"
	}

	contentW := boundedSpan(width, 4, 24, 62)
	contentX := centeredStart(width, contentW)
	row := eyeY + 4
	if row < 8 {
		row = 8
	}

	subjectLines := wrapAndClamp("SUBJECT: ", msg, contentW, 2)
	for _, line := range subjectLines {
		if row >= height {
			break
		}
		c.SetString(contentX, row, line, red)
		row++
	}

	if row < height {
		row++
	}

	responseLines := wrapAndClamp("RESPONSE REQUIRED: ", inputView, contentW, 2)
	for _, line := range responseLines {
		if row >= height {
			break
		}
		c.SetString(contentX, row, line, red)
		row++
	}

	if hint != "" && row < height {
		if row < height {
			row++
		}
		hintLines := wrapAndClamp("DATA FRAGMENT: ", hint, contentW, 2)
		for _, line := range hintLines {
			if row >= height {
				break
			}
			c.SetString(contentX, row, line, grey)
			row++
		}
	}

	return c.Render()
}

// Helper
func drawBox(c *canvas.Canvas, x, y, w, h int, style lipgloss.Style) {
	// Top/Bottom
	for i := 0; i < w; i++ {
		c.SetChar(x+i, y, '─', style)
		c.SetChar(x+i, y+h-1, '─', style)
	}
	// Sides
	for i := 0; i < h; i++ {
		c.SetChar(x, y+i, '│', style)
		c.SetChar(x+w-1, y+i, '│', style)
	}
	// Corners
	c.SetChar(x, y, '┌', style)
	c.SetChar(x+w-1, y, '┐', style)
	c.SetChar(x, y+h-1, '└', style)
	c.SetChar(x+w-1, y+h-1, '┘', style)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	Register(NewMatrixTheme)
	Register(NewCyberpunkTheme)
	Register(NewTronTheme)
	Register(NewAlienTheme)
	Register(NewSystemShockTheme)
}
