package theme

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"fmt"
	"math"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	numParticles = 800
)

const (
	ModeVector = iota
	ModeRings
	ModeDiamond
	ModeWaves
)

const (
	ColorThemeGoogle = iota
	ColorThemeFrostyBlue
	ColorThemeMatrix
	ColorThemeAnimatedRGB
	ColorThemeFire
	ColorThemeSparkle
)

const (
	LayoutDialog = iota
	LayoutDiamondSplat
	LayoutAvoidance
	LayoutFinale
)

type Particle struct {
	X, Y       float64
	VX, VY     float64 // Used for vector field mode
	TargetX    float64 // Used for shape modes
	TargetY    float64 // Used for shape modes
	Color      lipgloss.Style
	Char       rune
	ShapeGroup int // Used or ignored depending on shape
}

var googleColors = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#4285F4")), // Blue
	lipgloss.NewStyle().Foreground(lipgloss.Color("#EA4335")), // Red
	lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBC05")), // Yellow
	lipgloss.NewStyle().Foreground(lipgloss.Color("#34A853")), // Green
}

type AntigravityTheme struct {
	BaseTheme
	particles   []*Particle
	mode        int
	colorTheme  int
	layoutStyle int
	timeOffset  float64
	keyCount    int
	lastWidth   int
	lastHeight  int
	name        string
	description string

	// Used for Finale theme
	finaleTick int
}

func NewAntigravityThemeDialog() Theme {
	return newAntigravityTheme(LayoutDialog, "Antigravity (Dialog)", "Unified particles with classic dialog overlay")
}

func NewAntigravityThemeDiamond() Theme {
	return newAntigravityTheme(LayoutDiamondSplat, "Antigravity (Diamond Splat)", "Text exists in floating void bubbles inside the diamond")
}

func NewAntigravityThemeAvoidance() Theme {
	return newAntigravityTheme(LayoutAvoidance, "Antigravity (Avoidance)", "Particles organically avoid the text")
}

func NewAntigravityThemeFinale() Theme {
	return newAntigravityTheme(LayoutFinale, "Antigravity (Finale)", "Chaotic avoidance theme that auto-cycles")
}

func newAntigravityTheme(layout int, name, desc string) *AntigravityTheme {
	t := &AntigravityTheme{
		particles:   make([]*Particle, numParticles),
		layoutStyle: layout,
		colorTheme:  rand.Intn(4), // Spawn with random color theme
		name:        name,
		description: desc,
	}

	if layout == LayoutDiamondSplat {
		t.mode = ModeDiamond
	} else {
		t.mode = rand.Intn(4)
	}

	// Finale and new themes can use all colors. Splat uses new colors too.
	t.colorTheme = rand.Intn(6)

	for i := range t.particles {
		t.particles[i] = &Particle{}
		// Values will be initialized on first render when we know width/height
	}
	return t
}

func (t *AntigravityTheme) Name() string        { return t.name }
func (t *AntigravityTheme) Description() string { return t.description }

func (t *AntigravityTheme) Init() tea.Cmd {
	return Tick()
}

func (t *AntigravityTheme) resetParticle(p *Particle, width, height float64) {
	p.X = rand.Float64() * width
	p.Y = rand.Float64() * height
	p.VX = 0
	p.VY = 0
	p.Color = googleColors[rand.Intn(len(googleColors))]
	p.Char = '.'
}

func (t *AntigravityTheme) Update(msg tea.Msg) (Theme, tea.Cmd) {
	switch msg := msg.(type) {
	case game.TickMsg:
		t.timeOffset += 0.05
		if t.layoutStyle == LayoutFinale {
			t.finaleTick++
			// Auto cycle every ~2.5 seconds (approx 75 ticks at 30fps)
			if t.finaleTick > 75 {
				t.finaleTick = 0
				t.cycleMode()
				t.colorTheme = rand.Intn(6)
			}
		}
		return t, Tick()
	case tea.KeyMsg:
		// Don't count system keys for mode switching
		if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace || msg.Type == tea.KeyBackspace {
			t.keyCount++
			if t.keyCount%3 == 0 {
				t.cycleMode()
			}
		}
	}
	return t, nil
}

func (t *AntigravityTheme) cycleMode() {
	if t.layoutStyle == LayoutDiamondSplat {
		t.mode = ModeDiamond
		t.colorTheme = (t.colorTheme + 1) % 6
		t.recalculateTargets()
		return
	}

	t.mode = (t.mode + 1) % 4
	if t.mode == 0 {
		// Cycle colors when wrapping around modes to keep it fresh
		t.colorTheme = (t.colorTheme + 1) % 6
	}
	t.recalculateTargets()
}

func (t *AntigravityTheme) recalculateTargets() {
	if t.lastWidth <= 0 || t.lastHeight <= 0 {
		return
	}
	width, height := t.lastWidth, t.lastHeight

	switch t.mode {
	case ModeRings:
		t.generateRings(width, height)
	case ModeDiamond:
		t.generateDiamond(width, height)
	case ModeWaves:
		t.generateWaves(width, height)
	}
}

// flowField calculates the vector flow map
func flowField(x, y, timeOffset float64) (float64, float64) {
	scale := 0.08
	sx := x * scale
	sy := y * scale * 2.5
	angle := math.Sin(sx+timeOffset)*math.Cos(sy+timeOffset*0.5)*math.Pi + timeOffset*0.2
	magnitude := 0.3 + math.Abs(math.Sin(sx*0.5+timeOffset))*0.5
	return math.Cos(angle) * magnitude, math.Sin(angle) * magnitude
}

func (t *AntigravityTheme) generateRings(width, height int) {
	centerX, centerY := float64(width)/2.0, float64(height)/2.0
	centers := []struct{ x, y, r float64 }{
		{centerX, centerY - float64(height)*0.2, float64(height) * 0.3},
		{centerX - float64(width)*0.15, centerY + float64(height)*0.15, float64(height) * 0.3},
		{centerX + float64(width)*0.15, centerY + float64(height)*0.15, float64(height) * 0.3},
	}

	for i, p := range t.particles {
		ringIdx := i % 3
		angle := rand.Float64() * math.Pi * 2
		rOffset := rand.NormFloat64() * float64(height) * 0.04
		r := centers[ringIdx].r + rOffset

		p.TargetX = centers[ringIdx].x + math.Cos(angle)*r*2.0
		p.TargetY = centers[ringIdx].y + math.Sin(angle)*r
	}
}

func (t *AntigravityTheme) generateDiamond(width, height int) {
	centerX, centerY := float64(width)/2.0, float64(height)/2.0
	size := float64(height) * 0.35

	for _, p := range t.particles {
		var x, y float64
		side := rand.Intn(4)
		pos := rand.Float64()*size*2 - size

		switch side {
		case 0:
			x, y = pos, -size
		case 1:
			x, y = size, pos
		case 2:
			x, y = pos, size
		case 3:
			x, y = -size, pos
		}

		x += rand.NormFloat64() * float64(height) * 0.04
		y += rand.NormFloat64() * float64(height) * 0.04

		rotAngle := math.Pi / 4.0
		rotatedX := x*math.Cos(rotAngle) - y*math.Sin(rotAngle)
		rotatedY := x*math.Sin(rotAngle) + y*math.Cos(rotAngle)

		p.TargetX = centerX + rotatedX*2.5
		p.TargetY = centerY + rotatedY
	}
}

func (t *AntigravityTheme) generateWaves(width, height int) {
	bandY := []float64{float64(height) * 0.25, float64(height) * 0.5, float64(height) * 0.75}

	for i, p := range t.particles {
		bandIdx := i % 3
		p.ShapeGroup = bandIdx
		p.TargetX = rand.Float64()*float64(width)*0.8 + float64(width)*0.1
		p.TargetY = bandY[bandIdx]
	}
}

func getStreakChar(vx, vy float64) rune {
	angle := math.Atan2(vy, vx)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	deg := angle * 180 / math.Pi

	if deg >= 337.5 || deg < 22.5 || (deg >= 157.5 && deg < 202.5) {
		return '-'
	} else if (deg >= 22.5 && deg < 67.5) || (deg >= 202.5 && deg < 247.5) {
		return '\\'
	} else if (deg >= 67.5 && deg < 112.5) || (deg >= 247.5 && deg < 292.5) {
		return '|'
	}
	return '/'
}

type TextBox struct {
	X, Y, W, H int
	// Splat physics info
	CenterX, CenterY float64
	Radius           float64
}

func (t *AntigravityTheme) View(width, height int, q *game.Question, inputView string, hint string) string {
	c := canvas.New(width, height)

	// Handle resize and initialization
	if width != t.lastWidth || height != t.lastHeight {
		if t.lastWidth == 0 {
			// First time init
			for _, p := range t.particles {
				t.resetParticle(p, float64(width), float64(height))
			}
		}
		t.lastWidth = width
		t.lastHeight = height
		t.recalculateTargets()
	}

	// Prepare text components for all layouts
	boxW := min(60, width-4)
	innerW := boxW - 4

	questionLines := wrapText(q.Text, innerW)
	inputLines := wrapLabeled("> ", inputView, innerW)

	baseH := 8
	if hint != "" {
		baseH = 11
	}
	extraH := 0
	if len(questionLines) > 1 {
		extraH += min(len(questionLines)-1, 3)
	}
	if len(inputLines) > 1 {
		extraH += min(len(inputLines)-1, 2)
	}
	boxH := baseH + extraH

	// Determine layout specific bounds
	var avoidBoxes []TextBox
	var textStartX, textStartY int

	textStartX = (width - boxW) / 2
	textStartY = (height - boxH) / 2

	qHeight := len(questionLines)
	inHeight := len(inputLines)
	inY := textStartY + qHeight + 3

	switch t.layoutStyle {
	case LayoutDialog:
		// Standard box
	case LayoutAvoidance, LayoutFinale:
		// Rectangular bounding text boxes for hard avoidance
		avoidBoxes = append(avoidBoxes, TextBox{X: textStartX, Y: textStartY + 1, W: boxW, H: qHeight + 2})
		avoidBoxes = append(avoidBoxes, TextBox{X: textStartX, Y: inY, W: boxW, H: inHeight + 2})
		if hint != "" {
			avoidBoxes = append(avoidBoxes, TextBox{X: textStartX, Y: inY + 2 + inHeight, W: boxW, H: 3})
		}
	case LayoutDiamondSplat:
		// Calculate circular splat zones that drift and bob slightly
		drift1 := math.Sin(t.timeOffset*0.8) * 2.5
		drift2 := math.Cos(t.timeOffset*0.6) * 2.0
		drift3 := math.Sin(t.timeOffset*0.9) * 1.5

		// Question Bubble
		qCenterY := float64(textStartY+1) + float64(qHeight)/2.0 + drift1
		qRadius := math.Max(float64(boxW)/2.2, float64(qHeight)*1.5)
		avoidBoxes = append(avoidBoxes, TextBox{
			CenterX: float64(width)/2.0 + drift2,
			CenterY: qCenterY,
			Radius:  qRadius,
			X:       textStartX, Y: textStartY + 1 + int(drift1), W: boxW, H: qHeight + 2,
		})

		// Input Bubble
		inCenterY := float64(inY) + float64(inHeight)/2.0 + drift3
		inRadius := math.Max(float64(boxW)/2.5, float64(inHeight)*1.5)
		avoidBoxes = append(avoidBoxes, TextBox{
			CenterX: float64(width)/2.0 - drift1,
			CenterY: inCenterY,
			Radius:  inRadius,
			X:       textStartX, Y: inY + int(drift3), W: boxW, H: inHeight + 2,
		})

		// Hint Bubble
		if hint != "" {
			hCenterY := float64(inY+2+inHeight) + 1.0 + drift2
			avoidBoxes = append(avoidBoxes, TextBox{
				CenterX: float64(width)/2.0 + drift3,
				CenterY: hCenterY,
				Radius:  float64(innerW) / 2.5,
				X:       textStartX, Y: inY + 2 + inHeight + int(drift2), W: boxW, H: 3,
			})
		}
	}

	// Calculate and draw particles
	for _, p := range t.particles {
		if t.mode == ModeVector {
			flowX, flowY := flowField(p.X, p.Y, t.timeOffset)
			p.VX = p.VX*0.95 + flowX*0.05
			p.VY = p.VY*0.95 + flowY*0.05

			// Avoidance logic for vector mode
			if t.layoutStyle == LayoutAvoidance || t.layoutStyle == LayoutFinale {
				for _, box := range avoidBoxes {
					// Inflate box slightly for padding
					bx1 := float64(box.X - 2)
					by1 := float64(box.Y - 1)
					bx2 := float64(box.X + box.W + 2)
					by2 := float64(box.Y + box.H + 1)

					// Simple repulsion
					if p.X > bx1 && p.X < bx2 && p.Y > by1 && p.Y < by2 {
						cx := (bx1 + bx2) / 2
						cy := (by1 + by2) / 2

						dx := p.X - cx
						dy := p.Y - cy

						// Push outward faster based on distance to center
						p.VX += (dx * 0.3)
						p.VY += (dy * 0.3)
					}
				}
			} else if t.layoutStyle == LayoutDiamondSplat {
				// Circular bubble repulsion
				for _, splat := range avoidBoxes {
					// Squish factor for terminal fonts (taller than wide)
					distY := math.Abs(p.Y-splat.CenterY) * 2.0
					distX := math.Abs(p.X - splat.CenterX)
					effectiveDist := math.Hypot(distX, distY)

					if effectiveDist < splat.Radius {
						// Repel
						dx := p.X - splat.CenterX
						dy := p.Y - splat.CenterY
						if dx == 0 && dy == 0 {
							dx = 0.1
						}
						norm := math.Hypot(dx, dy)
						p.VX += (dx / norm) * 0.8
						p.VY += (dy / norm) * 0.4
					}
				}
			}

			p.X += p.VX
			p.Y += p.VY

			// Wrap around
			if p.X < 0 {
				p.X += float64(width)
			}
			if p.X >= float64(width) {
				p.X -= float64(width)
			}
			if p.Y < 0 {
				p.Y += float64(height)
			}
			if p.Y >= float64(height) {
				p.Y -= float64(height)
			}

			p.Char = getStreakChar(p.VX, p.VY)

		} else {
			// Shape target motion
			distToTarget := math.Hypot(p.TargetX-p.X, p.TargetY-p.Y)
			dynamicTargetY := p.TargetY

			if t.mode == ModeWaves {
				waveSpeedOffset := float64(p.ShapeGroup) * 2.0
				dynamicTargetY += math.Sin(p.TargetX*0.05+t.timeOffset+waveSpeedOffset) * (float64(height) * 0.1)
			}

			// Add a subtle rotation/bobbing to shapes over time
			if t.mode == ModeDiamond || t.mode == ModeRings {
				// Gentle global swim
				dynamicTargetY += math.Sin(p.TargetX*0.02+t.timeOffset*0.5) * 1.5
			}

			noiseX := (rand.Float64() - 0.5) * 0.5
			noiseY := (rand.Float64() - 0.5) * 0.5

			moveFactor := 0.08
			if distToTarget < 2.0 {
				moveFactor = 0.02
			}

			newX := p.X + (p.TargetX-p.X)*moveFactor + noiseX
			newY := p.Y + (dynamicTargetY-p.Y)*moveFactor + noiseY

			// Avoidance logic for shape mode
			if t.layoutStyle == LayoutAvoidance || t.layoutStyle == LayoutFinale {
				for _, box := range avoidBoxes {
					bx1 := float64(box.X - 2)
					by1 := float64(box.Y - 1)
					bx2 := float64(box.X + box.W + 2)
					by2 := float64(box.Y + box.H + 1)

					if newX > bx1 && newX < bx2 && newY > by1 && newY < by2 {
						// Nudge particles out of the box borders
						if newX < (bx1+bx2)/2 {
							newX = bx1
						} else {
							newX = bx2
						}
					}
				}
			} else if t.layoutStyle == LayoutDiamondSplat {
				// If a particle's target is inside the void, push it to the edge of the void
				for _, splat := range avoidBoxes {
					distY := math.Abs(newY-splat.CenterY) * 2.0
					distX := math.Abs(newX - splat.CenterX)
					effectiveDist := math.Hypot(distX, distY)

					if effectiveDist < splat.Radius {
						// Push outward to radius
						angle := math.Atan2(newY-splat.CenterY, newX-splat.CenterX)
						// Reverse the Y squish logic for resolution
						newX = splat.CenterX + math.Cos(angle)*splat.Radius
						newY = splat.CenterY + (math.Sin(angle)*splat.Radius)/2.0
					}
				}
			}

			p.X = newX
			p.Y = newY

			p.Char = '.'
			if distToTarget > float64(height)*0.4 {
				p.Char = '-'
			} else if rand.Float64() > 0.95 {
				p.Char = '*'
			}
		}

		// Draw to canvas
		ix, iy := int(p.X), int(p.Y)
		if ix >= 0 && ix < width && iy >= 0 && iy < height {
			style := p.Color

			switch t.colorTheme {
			case ColorThemeFrostyBlue:
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
			case ColorThemeMatrix:
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
			case ColorThemeAnimatedRGB:
				r := int(math.Sin(p.X*0.05+t.timeOffset)*127 + 128)
				g := int(math.Sin(p.Y*0.1+t.timeOffset*1.2)*127 + 128)
				b := int(math.Sin((p.X+p.Y)*0.05+t.timeOffset*0.8)*127 + 128)
				hexStr := fmt.Sprintf("#%02x%02x%02x", r, g, b)
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(hexStr))
			case ColorThemeFire:
				// Fire: Bright yellow near bottom, fading to orange/red upward, shifting based on time
				heat := (float64(height) - p.Y) / float64(height) // 0 (top) to 1 (bottom)
				fireOffset := math.Sin(p.X*0.3+t.timeOffset*2.0) * 0.2
				heat += fireOffset
				if heat > 1.0 {
					heat = 1.0
				}
				if heat < 0.0 {
					heat = 0.0
				}

				r := 255
				g := int(255 * (heat * heat))
				b := int(100 * math.Pow(heat, 4))
				hexStr := fmt.Sprintf("#%02x%02x%02x", r, g, b)
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(hexStr))
			case ColorThemeSparkle:
				// Greyscale with occasional spike of white
				baseIntensity := 100 + int(math.Sin(p.X*0.1+p.Y*0.1)*30)
				spike := 0
				if rand.Float64() > 0.98 {
					spike = 155
				}
				v := baseIntensity + spike
				if v > 255 {
					v = 255
				}
				hexStr := fmt.Sprintf("#%02x%02x%02x", v, v, v)
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(hexStr))
			}

			// Only draw if there isn't something more important there
			if c.Grid[iy][ix].Rune == ' ' {
				c.SetChar(ix, iy, p.Char, style)
			}
		}
	}

	// Render Text Layers
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	if t.layoutStyle == LayoutDialog {
		// Draw semi-transparent box (simulated by dense characters or just solid dark gray)
		boxBg := lipgloss.NewStyle().Background(lipgloss.Color("#111111")).Foreground(lipgloss.Color("#444444"))
		for y := textStartY; y < textStartY+boxH; y++ {
			for x := textStartX; x < textStartX+boxW; x++ {
				// Checkerboard pattern for simulated opacity
				if (x+y)%2 == 0 {
					c.SetChar(x, y, '░', boxBg)
				} else {
					c.SetChar(x, y, ' ', lipgloss.NewStyle().Background(lipgloss.Color("#000000")))
				}
			}
		}
		c.DrawBox(textStartX, textStartY, boxW, boxH, textStyle)
	}

	// Draw contents
	if t.layoutStyle == LayoutDiamondSplat {
		// Draw text starting at dynamic bubble coordinates
		if len(avoidBoxes) > 0 {
			qBox := avoidBoxes[0]
			row := qBox.Y
			for _, line := range questionLines {
				if row >= height {
					break
				}
				c.SetString(qBox.X+2, row, line, textStyle)
				row++
			}

			if len(avoidBoxes) > 1 {
				inBox := avoidBoxes[1]
				row = inBox.Y
				for _, line := range inputLines {
					if row >= height {
						break
					}
					c.SetString(inBox.X+2, row, line, textStyle)
					row++
				}
			}

			if hint != "" && len(avoidBoxes) > 2 {
				hBox := avoidBoxes[2]
				row = hBox.Y
				displayHint := truncateToWidth("CLUE: "+hint, innerW)
				hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
				c.SetString(hBox.X+2, row, displayHint, hintStyle)
			}
		}
	} else {
		row := textStartY + 2
		if t.layoutStyle != LayoutDialog {
			row = textStartY // Less padding if no box
		}

		for _, line := range questionLines {
			if row >= height {
				break
			}
			c.SetString(textStartX+2, row, line, textStyle)
			row++
		}

		row++ // gap

		for _, line := range inputLines {
			if row >= height {
				break
			}
			c.SetString(textStartX+2, row, line, textStyle)
			row++
		}

		if hint != "" {
			row++
			displayHint := truncateToWidth("CLUE: "+hint, innerW)
			hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
			c.SetString(textStartX+2, row, displayHint, hintStyle)
		}
	}

	return c.Render()
}

func init() {
	Register(NewAntigravityThemeDialog)
	Register(NewAntigravityThemeDiamond)
	Register(NewAntigravityThemeAvoidance)
	Register(NewAntigravityThemeFinale)
}
