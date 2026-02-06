package ui

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/theme"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"strings"
	"time"
)

type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Model struct {
	Config *game.Config
	State  GameState

	// Game State
	CurrentQuestionIndex int
	WrongAnswers         int
	CurrentTheme         theme.Theme

	// Animation State
	TransitionProgress float64
	TransitionText     string
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

func (m *Model) PickRandomTheme() {
	if len(theme.Registry) > 0 {
		m.CurrentTheme = theme.Registry[rand.Intn(len(theme.Registry))]
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		if m.State == StateIntro {
			if msg.Type == tea.KeyEnter {
				m.State = StateQuestion
				m.TypewriterIndex = 0
				return m, nil
			}
		} else if m.State == StateTransition {
			m.TransitionProgress += 0.05
			if m.TransitionProgress >= 1.0 {
				m.State = StateQuestion
				m.TransitionProgress = 0
				m.TypewriterIndex = 0 // Reset for next question
				return m, nil
			}
		} else if m.State == StateQuestion {
			switch msg.Type {
			case tea.KeyEnter:
				// Validate Answer
				currentQ := m.Config.Questions[m.CurrentQuestionIndex]
				if game.CheckAnswer(m.Input.Value(), currentQ.Answer) {
					// Correct!
					m.CurrentQuestionIndex++
					m.Input.Reset()
					m.ShowHint = false
					m.WrongAnswers = 0

					if m.CurrentQuestionIndex >= len(m.Config.Questions) {
						m.State = StateSuccess
					} else {
						// Start Transition
						m.State = StateTransition
						m.TransitionProgress = 0
						m.TransitionText = fmt.Sprintf("DECRYPTING NODE %02d KEY...", m.CurrentQuestionIndex+1)
						m.PickRandomTheme()
						return m, nil
					}
				} else {
					// Wrong
					m.WrongAnswers++
					if m.WrongAnswers >= 1 {
						m.ShowHint = true
					}
					// Reset input or shake? nah, keep value so they can fix typo
				}
			}
		}

	case TickMsg:
		if m.State == StateTransition {
			m.TransitionProgress += 0.01
			if m.TransitionProgress >= 1.0 {
				m.State = StateQuestion
				m.TransitionProgress = 0
				m.TypewriterIndex = 0
			}
			return m, tick()
		} else if m.State == StateQuestion {
			// Typewriter effect
			currentQ := m.Config.Questions[m.CurrentQuestionIndex]
			if m.TypewriterIndex < len(currentQ.Text) {
				m.TypewriterIndex++
				// Speed up typing by incrementing more if desired, or just depend on tick rate
			}
			return m, tick()
		}
		return m, tick()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	// Update components
	if m.State == StateQuestion {
		m.Input, cmd = m.Input.Update(msg)
	}

	return m, cmd
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
		width := 40
		filled := int(m.TransitionProgress * float64(width))
		if filled > width {
			filled = width
		}

		bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

		style := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#FF00FF"))

		return style.Render(fmt.Sprintf("%s\n\n[%s]\n\n(Mash keys to accelerate)", m.TransitionText, bar))

	case StateQuestion:
		q := m.Config.Questions[m.CurrentQuestionIndex]

		// Apply Typewriter Effect
		visibleText := q.Text
		if m.TypewriterIndex < len(q.Text) {
			visibleText = q.Text[:m.TypewriterIndex] + "█" // Cursor
		}

		displayQ := q
		displayQ.Text = visibleText

		hint := ""
		if m.ShowHint {
			hint = q.Hint
		}

		if m.CurrentTheme != nil {
			return m.CurrentTheme.Render(m.Width, m.Height, &displayQ, m.Input.View(), hint)
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
