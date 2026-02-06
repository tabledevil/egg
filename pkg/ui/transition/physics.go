package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/canvas"
	"fmt"
	"math"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 16. Explosion ---

type Explosion struct {
	BaseTransition
	oldLines  []string
	newLines  []string
	particles []particle
	time      float64
	done      bool
}

type particle struct {
	x, y   float64
	vx, vy float64
	char   rune
	color  string
}

func NewExplosion() Transition { return &Explosion{} }

func (t *Explosion) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *Explosion) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.time += 0.1

		// Initialize particles once
		if len(t.particles) == 0 && len(t.oldLines) > 0 {
			// Create particles from center or random
			for i := 0; i < 100; i++ {
				angle := rand.Float64() * 2 * math.Pi
				speed := 1.0 + rand.Float64()*2.0
				t.particles = append(t.particles, particle{
					x:     40, // default center fallback
					y:     12,
					vx:    math.Cos(angle) * speed,
					vy:    math.Sin(angle) * speed,
					char:  '*',
					color: "#FF0000",
				})
			}
		}

		// Update physics
		for i := range t.particles {
			t.particles[i].x += t.particles[i].vx
			t.particles[i].y += t.particles[i].vy
			t.particles[i].vy += 0.1 // Gravity
		}

		if t.time > 5.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *Explosion) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	// Init particles position if not set (first frame view)
	if len(t.particles) > 0 && t.time < 0.2 {
		// Reset to actual center
		for i := range t.particles {
			t.particles[i].x = float64(width / 2)
			t.particles[i].y = float64(height / 2)
		}
	}

	// Fade in New Screen
	alpha := t.time / 3.0
	if alpha > 1.0 { alpha = 1.0 }

	if alpha > 0.5 {
		for y := 0; y < height; y++ {
			if y < len(t.newLines) {
				c.SetString(0, y, t.newLines[y], lipgloss.NewStyle())
			}
		}
	}

	// Draw Particles
	for _, p := range t.particles {
		px, py := int(p.x), int(p.y)
		if px >= 0 && px < width && py >= 0 && py < height {
			c.SetChar(px, py, p.char, lipgloss.NewStyle().Foreground(lipgloss.Color(p.color)))
		}
	}

	return c.Render()
}

func (t *Explosion) Done() bool { return t.done }
func (t *Explosion) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 17. Teleporter ---

type Teleporter struct {
	BaseTransition
	oldLines []string
	newLines []string
	phase    int // 0=dissolve out, 1=dissolve in
	progress float64
	done     bool
}

func NewTeleporter() Transition { return &Teleporter{} }

func (t *Teleporter) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *Teleporter) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.progress += 0.05
		if t.progress >= 1.0 {
			if t.phase == 0 {
				t.phase = 1
				t.progress = 0
			} else {
				t.done = true
			}
		}
		return t, Tick()
	}
	return t, nil
}

func (t *Teleporter) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	sparkle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	lines := t.oldLines
	if t.phase == 1 { lines = t.newLines }

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Probability of showing char vs sparkle vs empty
			prob := 1.0 - t.progress
			if t.phase == 1 { prob = t.progress }

			if rand.Float64() < prob {
				char := ' '
				if y < len(lines) && x < len(lines[y]) {
					char = rune(lines[y][x])
				}
				c.SetChar(x, y, char, lipgloss.NewStyle())
			} else if rand.Float64() < 0.1 {
				c.SetChar(x, y, '*', sparkle)
			}
		}
	}
	return c.Render()
}

func (t *Teleporter) Done() bool { return t.done }
func (t *Teleporter) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}


// --- 18. Film Reel ---

type FilmReel struct {
	BaseTransition
	count int
	done  bool
}

func NewFilmReel() Transition { return &FilmReel{count: 3} }
func (t *FilmReel) SetContent(o, n string) {}

func (t *FilmReel) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		if rand.Float64() < 0.05 { // Slow countdown
			t.count--
			if t.count < 0 { t.done = true }
		}
		return t, Tick()
	}
	return t, nil
}

func (t *FilmReel) View(width, height int) string {
	c := canvas.New(width, height)

	// Sepia tone
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#C0A080"))

	// Big Number
	num := fmt.Sprintf("%d", t.count)
	if t.count == 0 { num = "START" }

	c.Fill(0, 0, width, height, ' ', style)
	c.DrawBox(2, 2, width-4, height-4, style)

	c.SetString(width/2, height/2, num, style.Bold(true))

	// Scratches
	for i := 0; i < 5; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		c.SetChar(x, y, '|', style)
	}

	return c.Render()
}

func (t *FilmReel) Done() bool { return t.done }


// --- 19. Defrag ---

type Defrag struct {
	BaseTransition
	oldLines []string
	newLines []string
	blocks   []int // 0=empty, 1=filled
	cursor   int
	done     bool
}

func NewDefrag() Transition { return &Defrag{} }

func (t *Defrag) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *Defrag) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.cursor += 50
		// Heuristic max
		if t.cursor > 5000 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *Defrag) View(width, height int) string {
	c := canvas.New(width, height)

	// Visualizing defrag as blocks
	// Not actual content transition, just visual

	for i := 0; i < width*height; i++ {
		x := i % width
		y := i / width

		if i < t.cursor {
			// Defragged (Green)
			c.SetChar(x, y, '█', lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))
		} else {
			// Fragmented (Red/Blue/Empty)
			if rand.Float64() < 0.3 {
				c.SetChar(x, y, '▓', lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")))
			} else {
				c.SetChar(x, y, '░', lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF")))
			}
		}
	}

	return c.Render()
}

func (t *Defrag) Done() bool { return t.done }


// --- 20. Wave Distortion ---

type WaveDistortion struct {
	BaseTransition
	oldLines []string
	newLines []string
	offset   float64
	done     bool
}

func NewWaveDistortion() Transition { return &WaveDistortion{} }

func (t *WaveDistortion) SetContent(o, n string) {
	t.oldLines = getLines(o)
	t.newLines = getLines(n)
}

func (t *WaveDistortion) Update(msg tea.Msg) (Transition, tea.Cmd) {
	if _, ok := msg.(game.TickMsg); ok {
		t.offset += 0.2
		if t.offset > 6.0 { t.done = true }
		return t, Tick()
	}
	return t, nil
}

func (t *WaveDistortion) View(width, height int) string {
	c := canvas.New(width, height)
	t.ensureLines(height)

	// Use Old until half way, then New?
	// Or just distort Old then Fade to New?
	// Let's distort Old -> New swap -> distort New

	lines := t.oldLines
	if t.offset > 3.0 { lines = t.newLines }

	for y := 0; y < height; y++ {
		// Sine wave shift x
		shift := int(5.0 * math.Sin(float64(y)*0.5 + t.offset))

		line := ""
		if y < len(lines) { line = lines[y] }

		for x := 0; x < width; x++ {
			srcX := x - shift
			if srcX >= 0 && srcX < len(line) {
				c.SetChar(x, y, rune(line[srcX]), lipgloss.NewStyle())
			}
		}
	}
	return c.Render()
}

func (t *WaveDistortion) Done() bool { return t.done }
func (t *WaveDistortion) ensureLines(h int) {
	for len(t.oldLines) < h { t.oldLines = append(t.oldLines, "") }
	for len(t.newLines) < h { t.newLines = append(t.newLines, "") }
}

func init() {
	Register(NewExplosion)
	Register(NewTeleporter)
	Register(NewFilmReel)
	Register(NewDefrag)
	Register(NewWaveDistortion)
}
