package transition

import (
	"ctf-tool/pkg/game"
	"math/rand"
	"strings"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

type MatrixTransition struct {
	progress float64
	columns  []int
	width    int
	height   int
}

func NewMatrixTransition() Transition {
	return &MatrixTransition{progress: 0}
}

func (t *MatrixTransition) Init() tea.Cmd {
	return tick()
}

func (t *MatrixTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	switch msg.(type) {
	case game.TickMsg:
		t.progress += 0.008

		if t.width > 0 {
			if len(t.columns) != t.width {
				t.columns = make([]int, t.width)
				for i := range t.columns {
					t.columns[i] = rand.Intn(t.height + 20) - 20
				}
			}
			for i := range t.columns {
				t.columns[i]++
				if t.columns[i] > t.height + 10 {
					t.columns[i] = rand.Intn(10) - 10
				}
			}
		}
		return t, tick()

	case tea.KeyMsg:
		t.progress += 0.05
		return t, nil
	}
	return t, nil
}

func (t *MatrixTransition) View(width, height int) string {
	t.width = width
	t.height = height

	var sb strings.Builder
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	runes := []rune("0123456789ABCDEF")

	if len(t.columns) != width {
		t.columns = make([]int, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			head := t.columns[x]
			if y == head {
				sb.WriteString(white.Render(string(runes[rand.Intn(len(runes))])))
			} else if y < head && y > head - 10 {
				sb.WriteString(green.Render(string(runes[rand.Intn(len(runes))])))
			} else {
				sb.WriteRune(' ')
			}
		}
		if y < height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func (t *MatrixTransition) Done() bool {
	return t.progress >= 1.0
}

func init() {
	Register(NewMatrixTransition)
}
