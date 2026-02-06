package ui

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"time"
)

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
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
		return m.ActiveTheme.Init()
	}
	return nil
}

func (m *Model) StartTransition() tea.Cmd {
	m.State = StateTransition
	if len(transition.Registry) > 0 {
		constructor := transition.Registry[rand.Intn(len(transition.Registry))]
		m.ActiveTransition = constructor()
		return m.ActiveTransition.Init()
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, tick())
	if m.ActiveTheme != nil {
		cmds = append(cmds, m.ActiveTheme.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	// State Machine
	switch m.State {
	case StateIntro:
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			m.State = StateQuestion
			m.TypewriterIndex = 0
			if m.ActiveTheme != nil {
				cmds = append(cmds, m.ActiveTheme.Init())
			}
		}

	case StateTransition:
		if m.ActiveTransition != nil {
			var tCmd tea.Cmd
			m.ActiveTransition, tCmd = m.ActiveTransition.Update(msg)
			cmds = append(cmds, tCmd)

			if m.ActiveTransition.Done() {
				m.State = StateQuestion
				m.TypewriterIndex = 0
				cmds = append(cmds, m.PickRandomTheme())
			}
		} else {
			m.State = StateQuestion
		}

	case StateQuestion:
		// Theme Update
		if m.ActiveTheme != nil {
			var tCmd tea.Cmd
			m.ActiveTheme, tCmd = m.ActiveTheme.Update(msg)
			cmds = append(cmds, tCmd)
		}

		// Input Handling
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			currentQ := m.Config.Questions[m.CurrentQuestionIndex]
			if game.CheckAnswer(m.Input.Value(), currentQ.Answer) {
				m.CurrentQuestionIndex++
				m.Input.Reset()
				m.ShowHint = false
				m.WrongAnswers = 0

				if m.CurrentQuestionIndex >= len(m.Config.Questions) {
					m.State = StateSuccess
				} else {
					cmds = append(cmds, m.StartTransition())
				}
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
			cmds = append(cmds, tick())
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
			return m.ActiveTransition.View(m.Width, m.Height)
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

		if m.ActiveTheme != nil {
			return m.ActiveTheme.View(m.Width, m.Height, &displayQ, m.Input.View(), hint)
		}
		return "Error: No Theme Selected"

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
