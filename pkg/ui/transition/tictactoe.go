package transition

import (
	"ctf-tool/pkg/game"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"math/rand"
	"fmt"
)

type TicTacToeTransition struct {
	progress float64
	board    [9]string // " ", "X", "O"
	cursor   int
	turn     string // "X" (Player) or "O" (AI)
	winner   string
	aiTimer  int
}

func NewTicTacToeTransition() Transition {
	t := &TicTacToeTransition{
		progress: 0,
		turn: "X",
	}
	for i := range t.board {
		t.board[i] = " "
	}
	return t
}

func (t *TicTacToeTransition) Init() tea.Cmd {
	return tick()
}

func (t *TicTacToeTransition) Update(msg tea.Msg) (Transition, tea.Cmd) {
	switch msg := msg.(type) {
	case game.TickMsg:
		t.progress += 0.005
		t.aiTimer++

		if t.turn == "O" && t.aiTimer > 10 { // AI delay
			empty := []int{}
			for i, v := range t.board {
				if v == " " { empty = append(empty, i) }
			}
			if len(empty) > 0 {
				move := empty[rand.Intn(len(empty))]
				t.board[move] = "O"
				t.checkWin()
				t.turn = "X"
			}
			t.aiTimer = 0
		} else if t.turn == "X" && t.aiTimer > 30 {
			// Auto-play for player
			empty := []int{}
			for i, v := range t.board {
				if v == " " { empty = append(empty, i) }
			}
			if len(empty) > 0 {
				move := empty[rand.Intn(len(empty))]
				t.board[move] = "X"
				t.checkWin()
				t.turn = "O"
			}
			t.aiTimer = 0
		}

		return t, tick()

	case tea.KeyMsg:
		if t.turn == "X" && t.winner == "" {
			switch msg.String() {
			case "up":
				t.cursor -= 3
			case "down":
				t.cursor += 3
			case "left":
				t.cursor -= 1
			case "right":
				t.cursor += 1
			case "enter", " ":
				if t.cursor >= 0 && t.cursor < 9 && t.board[t.cursor] == " " {
					t.board[t.cursor] = "X"
					t.checkWin()
					t.turn = "O"
					t.aiTimer = 0
				}
			}
			if t.cursor < 0 { t.cursor += 9 }
			if t.cursor > 8 { t.cursor -= 9 }
		} else {
			t.progress += 0.05
		}
	}
	return t, nil
}

func (t *TicTacToeTransition) checkWin() {
	wins := [][]int{
		{0,1,2}, {3,4,5}, {6,7,8},
		{0,3,6}, {1,4,7}, {2,5,8},
		{0,4,8}, {2,4,6},
	}
	for _, w := range wins {
		if t.board[w[0]] != " " && t.board[w[0]] == t.board[w[1]] && t.board[w[1]] == t.board[w[2]] {
			t.winner = t.board[w[0]]
			t.progress = 1.0
		}
	}
	full := true
	for _, v := range t.board {
		if v == " " { full = false }
	}
	if full {
		t.winner = "Draw"
		t.progress = 1.0
	}
}

func (t *TicTacToeTransition) View(width, height int) string {
	var s strings.Builder

	s.WriteString("BYPASS SECURITY PROTOCOL: WIN OR WAIT...\n\n")

	for i := 0; i < 9; i++ {
		cell := t.board[i]
		if i == t.cursor {
			cell = "[" + cell + "]"
		} else {
			cell = " " + cell + " "
		}

		s.WriteString(cell)
		if (i+1) % 3 == 0 {
			if i < 8 {
				s.WriteString("\n---+---+---\n")
			}
		} else {
			s.WriteString("|")
		}
	}

	bar := int(t.progress * 20)
	if bar > 20 { bar = 20 }
	s.WriteString(fmt.Sprintf("\n\nBRUTE FORCE: [%s%s]", strings.Repeat("#", bar), strings.Repeat("-", 20-bar)))

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s.String())
}

func (t *TicTacToeTransition) Done() bool {
	return t.progress >= 1.0
}

func init() {
	Register(NewTicTacToeTransition)
}
