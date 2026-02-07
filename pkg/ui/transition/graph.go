package transition

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/caps"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"math/rand"
	"strings"
)

type Point struct {
	X, Y int
}

type GraphTransition struct {
	progress float64
	nodes    []Point
	width    int
	height   int
}

func NewGraphTransition() Transition {
	return &GraphTransition{progress: 0}
}

func (t *GraphTransition) Init() tea.Cmd {
	return tick()
}

func (t *GraphTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	switch msg.(type) {
	case game.TickMsg:
		t.progress += 0.01
		return t, tick()
	case tea.KeyMsg:
		t.progress += 0.05
		return t, nil
	}
	return t, nil
}

func (t *GraphTransition) View(width, height int) string {
	if width != t.width || height != t.height {
		t.width = width
		t.height = height
		t.nodes = make([]Point, 15)
		for i := range t.nodes {
			t.nodes[i] = Point{
				X: rand.Intn(width),
				Y: rand.Intn(height),
			}
		}
	}

	canvas := make([][]rune, height)
	for y := range canvas {
		canvas[y] = make([]rune, width)
		for x := range canvas[y] {
			canvas[y][x] = ' '
		}
	}

	for _, n := range t.nodes {
		if n.Y < height && n.X < width && n.Y >= 0 && n.X >= 0 {
			canvas[n.Y][n.X] = 'O'
		}
	}

	for i := 0; i < len(t.nodes)-1; i++ {
		if float64(i)/float64(len(t.nodes)) > t.progress {
			continue
		}

		p1 := t.nodes[i]
		p2 := t.nodes[i+1]

		dx := float64(p2.X - p1.X)
		dy := float64(p2.Y - p1.Y)
		steps := float64(int(max(abs(dx), abs(dy))))

		xInc := dx / steps
		yInc := dy / steps

		x := float64(p1.X)
		y := float64(p1.Y)

		for j := 0; j < int(steps); j++ {
			ix, iy := int(x), int(y)
			if iy >= 0 && iy < height && ix >= 0 && ix < width {
				canvas[iy][ix] = '.'
			}
			x += xInc
			y += yInc
		}
	}

	var sb strings.Builder
	for y := 0; y < height; y++ {
		sb.WriteString(string(canvas[y]))
		if y < height-1 {
			sb.WriteRune('\n')
		}
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Render(sb.String())
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (t *GraphTransition) Done() bool {
	return t.progress >= 1.0
}

func (t *GraphTransition) IsCompatible(c caps.Capabilities) bool {
	return c.ColorProfile <= termenv.ANSI
}

func init() {
	Register(NewGraphTransition)
}
