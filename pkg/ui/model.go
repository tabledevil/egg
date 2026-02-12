package ui

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var sgrTextPattern = regexp.MustCompile(`\[[0-9;]*m`)

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

type Model struct {
	Config *game.Config
	State  GameState

	// Game State
	CurrentQuestionIndex int
	WrongAnswers         int
	ActiveTheme          theme.Theme
	ActiveTransition     transition.Transition

	// Animation State
	TypewriterIndex    int

	// UI Components
	Input textinput.Model

	// Dimensions
	Width  int
	Height int

	// Feedback
	ShowHint bool

	// Demo
	AutoDemo bool
	DemoTick int
}

func NewModel(config *game.Config) Model {
	rand.Seed(time.Now().UnixNano())

	ti := textinput.New()
	ti.Placeholder = "Type answer..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	m := Model{
		Config: config,
		State:  StateIntro,
		Input:  ti,
	}
	m.PickRandomTheme()
	return m
}

func (m *Model) PickRandomTheme() tea.Cmd {
	if len(theme.Registry) > 0 {
		constructor := theme.Registry[rand.Intn(len(theme.Registry))]
		m.ActiveTheme = constructor()
	}
	return nil
}

func (m *Model) StartTransition() tea.Cmd {
	// 1. Capture Old View
	q := m.Config.Questions[m.CurrentQuestionIndex]
	visibleText := q.Text
	if m.TypewriterIndex < len(q.Text) {
		visibleText = q.Text[:m.TypewriterIndex] + "█"
	}
	displayQ := q
	displayQ.Text = visibleText
	hint := ""
	if m.ShowHint { hint = q.Hint }

	oldView := m.safeThemeView(&displayQ, m.themeInputValue(), hint)

	// 2. Advance State
	m.CurrentQuestionIndex++
	// Wrap around for demo/endless feel, or success
	if m.CurrentQuestionIndex >= len(m.Config.Questions) {
		m.CurrentQuestionIndex = 0 // Loop for demo purposes
		// m.State = StateSuccess
		// return nil
	}

	// 3. Pick New Theme
	m.PickRandomTheme()

	// 4. Capture New View (Preview)
	m.Input.Reset()
	m.ShowHint = false
	m.WrongAnswers = 0
	m.TypewriterIndex = 0

	newQ := m.Config.Questions[m.CurrentQuestionIndex]
	newDisplayQ := newQ
	newDisplayQ.Text = "█"

	newView := m.safeThemeView(&newDisplayQ, m.themeInputValue(), "")

	// 5. Create Transition
	m.State = StateTransition
	if len(transition.Registry) > 0 {
		constructor := transition.Registry[rand.Intn(len(transition.Registry))]
		m.ActiveTransition = constructor()
		m.ActiveTransition.SetContent(oldView, newView)
		return m.ActiveTransition.Init()
	}

	m.State = StateQuestion
	return nil
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, tick())
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlX, tea.KeyEsc, tea.KeyF12:
			return m, tea.Quit
		case tea.KeyF1:
			// Next Theme
			cmds = append(cmds, m.PickRandomTheme())
		case tea.KeyF2:
			// Force Transition (stay on same Q)
			m.CurrentQuestionIndex--
			if m.CurrentQuestionIndex < 0 { m.CurrentQuestionIndex = 0 }
			cmds = append(cmds, m.StartTransition())
		case tea.KeyF3:
			m.AutoDemo = !m.AutoDemo
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	// Auto Demo Logic
	if m.AutoDemo {
		if _, ok := msg.(game.TickMsg); ok {
			m.DemoTick++
			if m.DemoTick > 150 { // ~5 seconds
				m.DemoTick = 0
				if m.State == StateQuestion {
					cmds = append(cmds, m.StartTransition())
				}
			}
		}
	}

	// State Machine
	switch m.State {
	case StateIntro:
		if _, ok := msg.(game.TickMsg); ok {
			cmds = append(cmds, tick())
		}
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			m.State = StateQuestion
			m.TypewriterIndex = 0
		}

	case StateTransition:
		if m.ActiveTransition != nil {
			var tCmd tea.Cmd
			nextTransition, tCmd := m.ActiveTransition.Update(msg)
			if nextTransition != nil {
				m.ActiveTransition = nextTransition
			}
			cmds = append(cmds, tCmd)

			if m.ActiveTransition.Done() {
				m.State = StateQuestion
				m.TypewriterIndex = 0
			}
		} else {
			m.State = StateQuestion
		}

	case StateQuestion:
		// Theme Update
		if m.ActiveTheme != nil {
			var tCmd tea.Cmd
			nextTheme, tCmd := m.ActiveTheme.Update(msg)
			if nextTheme != nil {
				m.ActiveTheme = nextTheme
			}
			cmds = append(cmds, tCmd)
		}

		// Input Handling
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			currentQ := m.Config.Questions[m.CurrentQuestionIndex]
			if game.CheckAnswer(m.Input.Value(), currentQ.Answer) {
				cmds = append(cmds, m.StartTransition())
			} else {
				m.WrongAnswers++
				if m.WrongAnswers >= 1 {
					m.ShowHint = true
				}
			}
		} else {
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Typewriter logic
		if _, ok := msg.(game.TickMsg); ok {
			currentQ := m.Config.Questions[m.CurrentQuestionIndex]
			if m.TypewriterIndex < len(currentQ.Text) {
				m.TypewriterIndex++
			}
		}

	case StateSuccess:
		// Success state logic
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Width == 0 {
		return "Loading..."
	}

	switch m.State {
	case StateIntro:
		style := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#00FF00"))

		return style.Render("SYSTEM BOOT SEQUENCE INITIATED...\n\n[PRESS ENTER TO HACK THE PLANET]")

	case StateTransition:
		if m.ActiveTransition != nil {
			return m.safeTransitionView()
		}
		return "Loading next level..."

	case StateQuestion:
		q := m.Config.Questions[m.CurrentQuestionIndex]

		visibleText := q.Text
		if m.TypewriterIndex < len(q.Text) {
			visibleText = q.Text[:m.TypewriterIndex] + "█"
		}

		displayQ := q
		displayQ.Text = visibleText

		hint := ""
		if m.ShowHint {
			hint = q.Hint
		}

		content := "Error: No Theme Selected"
		if m.ActiveTheme != nil {
			if rendered := m.safeThemeView(&displayQ, m.themeInputValue(), hint); rendered != "" {
				content = rendered
			}
		}

		// Overlay Demo Status
		if m.AutoDemo {
			content = lipgloss.JoinVertical(lipgloss.Left, content, lipgloss.NewStyle().Background(lipgloss.Color("#FF0000")).Render(" DEMO MODE "))
		}
		return content

	case StateSuccess:
		style := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(2)

		return style.Render(fmt.Sprintf("ACCESS GRANTED\n\n%s\n\n%s", m.Config.FinalMessage, m.Config.FinalHint))
	}

	return ""
}

func (m *Model) safeThemeView(q *game.Question, inputView, hint string) (view string) {
	if m.ActiveTheme == nil || m.Width <= 0 || m.Height <= 0 {
		return ""
	}

	defer func() {
		if recover() != nil {
			view = ""
		}
	}()

	return m.ActiveTheme.View(m.Width, m.Height, q, inputView, hint)
}

func (m *Model) safeTransitionView() (view string) {
	if m.ActiveTransition == nil || m.Width <= 0 || m.Height <= 0 {
		return "Loading next level..."
	}

	defer func() {
		if recover() != nil {
			view = "Loading next level..."
		}
	}()

	return m.ActiveTransition.View(m.Width, m.Height)
}

func (m *Model) themeInputValue() string {
	plain := ansi.Strip(m.Input.Value())
	plain = sgrTextPattern.ReplaceAllString(plain, "")
	plain = strings.ReplaceAll(plain, "\n", " ")
	plain = strings.ReplaceAll(plain, "\r", " ")
	return plain
}
