package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 11. Game Boy Theme ---

type GameboyTheme struct {
	BaseTheme
}

func NewGameboyTheme() Theme { return &GameboyTheme{} }
func (t *GameboyTheme) Name() string { return "Dot Matrix Game" }
func (t *GameboyTheme) Description() string { return "160x144 pixels of fun" }

func (t *GameboyTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// GB Palette
	darkest := lipgloss.NewStyle().Foreground(lipgloss.Color("#0F380F")).Background(lipgloss.Color("#8BAC0F"))
	dark := lipgloss.NewStyle().Foreground(lipgloss.Color("#306230")).Background(lipgloss.Color("#8BAC0F"))
	light := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BAC0F")).Background(lipgloss.Color("#9BBC0F"))
	lightest := lipgloss.NewStyle().Foreground(lipgloss.Color("#9BBC0F")).Background(lipgloss.Color("#9BBC0F"))

	// Fill background (Lightest Green)
	c.Fill(0, 0, width, height, ' ', lightest)

	// Viewport (160x144 ratio approx)
	vpW := 40
	vpH := 20
	vpX := (width - vpW) / 2
	vpY := (height - vpH) / 2

	// Draw Bezel
	c.DrawBox(vpX-1, vpY-1, vpW+2, vpH+2, dark)
	c.SetString(vpX, vpY-2, "DOT MATRIX WITH STEREO SOUND", dark.Italic(true))

	// Screen content background
	c.Fill(vpX, vpY, vpW, vpH, ' ', light)

	// Pokemon-style text box
	boxH := 6
	c.DrawBox(vpX, vpY+vpH-boxH, vpW, boxH, darkest)

	// Content
	c.SetString(vpX+2, vpY+2, "Trainer wants to battle!", darkest)
	c.SetString(vpX+2, vpY+4, "Q: "+q.Text, darkest)

	c.SetString(vpX+2, vpY+vpH-boxH+2, "> "+inputView, darkest)

	// Sprite placeholder
	c.SetString(vpX+vpW-5, vpY+vpH-boxH-1, "🐭", darkest) // Mouse/Pikachu

	if hint != "" {
		c.SetString(vpX+2, vpY+vpH-boxH+4, hint, dark)
	}

	return c.Render()
}

// --- 12. NES RPG Theme ---

type NESTheme struct {
	BaseTheme
}

func NewNESTheme() Theme { return &NESTheme{} }
func (t *NESTheme) Name() string { return "8-Bit RPG" }
func (t *NESTheme) Description() string { return "It's dangerous to go alone" }

func (t *NESTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#000000"))
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#000000"))

	// Background map pattern (grass)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%4 == 0 {
				c.SetChar(x, y, '"', lipgloss.NewStyle().Foreground(lipgloss.Color("#00AA00")))
			} else {
				c.SetChar(x, y, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#008800")))
			}
		}
	}

	// Dialog Box
	boxW := min(50, width-4)
	boxH := 10
	boxX := (width - boxW) / 2
	boxY := height - boxH - 2

	// Black background for box
	c.Fill(boxX, boxY, boxW, boxH, ' ', white)

	// Fancy double border
	c.DrawBox(boxX, boxY, boxW, boxH, border)

	// Text
	c.SetString(boxX+2, boxY+2, "OLD MAN:", white)
	c.SetString(boxX+2, boxY+4, q.Text, white)
	c.SetString(boxX+4, boxY+6, "▶ "+inputView, white)

	if hint != "" {
		c.SetString(boxX+2, boxY+8, "HINT: "+hint, white)
	}

	return c.Render()
}

// --- 13. SNES Mode 7 Theme ---

type SNESTheme struct {
	BaseTheme
	rotation float64
}

func NewSNESTheme() Theme { return &SNESTheme{} }
func (t *SNESTheme) Name() string { return "Super Mode 7" }
func (t *SNESTheme) Description() string { return "F-Zero Style" }

func (t *SNESTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.rotation += 0.05
		return t, Tick()
	}
	return t, nil
}

func (t *SNESTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Horizon
	horizonY := height / 3
	skyStyle := lipgloss.NewStyle().Background(lipgloss.Color("#87CEEB")) // Sky Blue
	c.Fill(0, 0, width, horizonY, ' ', skyStyle)

	// Mode 7 Floor
	for y := horizonY; y < height; y++ {
		// Perspective calculation
		z := float64(y - horizonY)
		scale := 10.0 / (z + 0.1)

		for x := 0; x < width; x++ {
			// Rotate coordinates
			worldX := float64(x - width/2) * scale
			worldY := 100.0 / scale // Move forward

			rotX := worldX*math.Cos(t.rotation) - worldY*math.Sin(t.rotation)
			rotY := worldX*math.Sin(t.rotation) + worldY*math.Cos(t.rotation)

			// Checkerboard
			if (int(rotX/5)+int(rotY/5))%2 == 0 {
				c.SetChar(x, y, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#555555")))
			} else {
				c.SetChar(x, y, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#333333")))
			}
		}
	}

	// Floating Text Box
	boxW := 40
	boxX := (width - boxW) / 2
	boxY := height / 2

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true)
	c.SetString(boxX, boxY, q.Text, style)
	c.SetString(boxX, boxY+2, "> "+inputView, style)

	return c.Render()
}

// --- 14. Fallout Pip-Boy Theme ---

type FalloutTheme struct {
	BaseTheme
}

func NewFalloutTheme() Theme { return &FalloutTheme{} }
func (t *FalloutTheme) Name() string { return "Pip-Boy 3000" }
func (t *FalloutTheme) Description() string { return "Vault-Tec Approved" }

func (t *FalloutTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#001100"))

	// Scan lines background
	for y := 0; y < height; y += 2 {
		c.Fill(0, y, width, 1, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#002200")))
	}

	// UI Frame
	c.DrawBox(2, 2, width-4, height-4, green)

	// Header
	c.SetString(4, 2, " STATS   ITEMS   DATA ", green.Background(lipgloss.Color("#003300")))
	c.SetString(width-10, 2, " HP 100/100 ", green)

	// Vault Boy (ASCII Art Placeholder)
	vbX := 5
	vbY := 5
	c.SetString(vbX, vbY, " (^_^)", green)
	c.SetString(vbX, vbY+1, "/|  |\\", green)
	c.SetString(vbX, vbY+2, " |__| ", green)
	c.SetString(vbX, vbY+3, "  LL  ", green)

	// Content
	contentX := 20
	c.SetString(contentX, 6, "QUEST: Answer the Question", green)
	c.SetString(contentX, 8, q.Text, green)

	c.SetString(contentX, 12, "> "+inputView+"_", green)

	// Footer
	c.SetString(4, height-3, " [ENTER] SELECT   [TAB] BACK ", green)

	return c.Render()
}

// --- 15. Deus Ex Theme ---

type DeusExTheme struct {
	BaseTheme
}

func NewDeusExTheme() Theme { return &DeusExTheme{} }
func (t *DeusExTheme) Name() string { return "Augmented Reality" }
func (t *DeusExTheme) Description() string { return "I never asked for this" }

func (t *DeusExTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	gold := lipgloss.NewStyle().Foreground(lipgloss.Color("#D4AF37")) // Gold
	black := lipgloss.NewStyle().Background(lipgloss.Color("#000000"))

	// Hexagonal pattern hint (dots)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x%4 == 0 && y%2 == 0 {
				c.SetChar(x, y, '·', gold.Inherit(black))
			}
		}
	}

	// UI Elements (Triangles/Angles)
	c.SetString(2, 2, "► SYSTEM ACCESS", gold)
	c.SetString(width-15, 2, "STATUS: OK ◄", gold)

	// Main Bar
	c.Fill(0, height/2-1, width, 1, '─', gold)

	// Content
	c.SetString(10, height/2-3, q.Text, gold.Bold(true))
	c.SetString(10, height/2+1, "INPUT: "+inputView, gold)

	// Decoration
	c.SetString(width-10, height-5, "AUGMENTATION", gold)
	c.SetString(width-10, height-4, "ACTIVE", gold)

	return c.Render()
}

func init() {
	Register(NewGameboyTheme)
	Register(NewNESTheme)
	Register(NewSNESTheme)
	Register(NewFalloutTheme)
	Register(NewDeusExTheme)
}
