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
	"sort"
	"time"
	"unicode/utf8"
)

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
	TypewriterIndex int

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
type questionLayout struct {
	maxLineChars int
	maxLines     int
	multiline    bool
}

func (m *Model) pickBestThemeForQuestion(question string) {
	if len(theme.Registry) == 0 {
		return
	}

	type candidate struct {
		theme theme.Theme
		score int
	}

	cands := make([]candidate, 0, len(theme.Registry))
	for _, constructor := range theme.Registry {
		t := constructor()
		_, shown, truncated, layout := m.fitQuestionForTheme(t, question)
		score := shown
		if layout.multiline {
			score += 3
		}
		if !truncated {
			score += 2000
		}
		cands = append(cands, candidate{theme: t, score: score})
	}

	rand.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	sort.SliceStable(cands, func(i, j int) bool { return cands[i].score > cands[j].score })
	m.ActiveTheme = cands[0].theme
}

func (m *Model) questionLayoutForTheme(t theme.Theme) questionLayout {
	minWidth := func(v int) int {
		if v < 8 {
			return 8
		}
		return v
	}

	switch t.(type) {
	case *theme.MatrixTheme:
		boxW := min(60, m.Width-4)
		boxH := min(15, m.Height-4)
		return questionLayout{maxLineChars: minWidth(boxW - 4), maxLines: max(1, boxH-8), multiline: true}
	case *theme.VHSTheme:
		return questionLayout{maxLineChars: minWidth(m.Width - 8), maxLines: max(1, m.Height/2), multiline: true}
	case *theme.BBSTheme:
		return questionLayout{maxLineChars: minWidth(m.Width - 14), maxLines: 2, multiline: true}
	case *theme.C64Theme:
		return questionLayout{maxLineChars: minWidth(m.Width - 15), maxLines: 2, multiline: true}
	case *theme.DOSTheme:
		winW := min(60, m.Width-2)
		return questionLayout{maxLineChars: minWidth(winW - 4), maxLines: 2, multiline: true}
	case *theme.GameboyTheme:
		return questionLayout{maxLineChars: 34, maxLines: 2, multiline: true}
	case *theme.NESTheme:
		boxW := min(50, m.Width-4)
		return questionLayout{maxLineChars: minWidth(boxW - 4), maxLines: 2, multiline: true}
	default:
		return questionLayout{maxLineChars: minWidth(m.Width - 12), maxLines: 1, multiline: false}
	}
}

func (m *Model) fitQuestionForTheme(t theme.Theme, text string) (string, int, bool, questionLayout) {
	layout := m.questionLayoutForTheme(t)
	if layout.multiline {
		lines, truncated := wrapText(text, layout.maxLineChars, layout.maxLines)
		shown := utf8.RuneCountInString(strings.Join(lines, ""))
		return strings.Join(lines, "\n"), shown, truncated, layout
	}

	line, truncated := truncateWithEllipsis(text, layout.maxLineChars)
	return line, utf8.RuneCountInString(line), truncated, layout
}

func wrapText(text string, maxWidth, maxLines int) ([]string, bool) {
	if maxWidth < 1 {
		maxWidth = 1
	}
	if maxLines < 1 {
		maxLines = 1
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}, false
	}

	lines := make([]string, 0, maxLines)
	current := ""
	appendLine := func(line string) bool {
		if len(lines) >= maxLines {
			return false
		}
		lines = append(lines, line)
		return true
	}

	for _, word := range words {
		runes := []rune(word)
		for len(runes) > maxWidth {
			if current != "" {
				if !appendLine(current) {
					return appendEllipsis(lines, maxWidth), true
				}
				current = ""
			}
			if !appendLine(string(runes[:maxWidth])) {
				return appendEllipsis(lines, maxWidth), true
			}
			runes = runes[maxWidth:]
		}
		word = string(runes)

		candidate := word
		if current != "" {
			candidate = current + " " + word
		}
		if utf8.RuneCountInString(candidate) <= maxWidth {
			current = candidate
			continue
		}

		if !appendLine(current) {
			return appendEllipsis(lines, maxWidth), true
		}
		current = word
	}

	if current != "" {
		if !appendLine(current) {
			return appendEllipsis(lines, maxWidth), true
		}
	}

	return lines, false
}

func appendEllipsis(lines []string, maxWidth int) []string {
	if len(lines) == 0 {
		return []string{"…"}
	}
	lines[len(lines)-1] = forceEllipsis(lines[len(lines)-1], maxWidth)
	return lines
}

func forceEllipsis(text string, maxWidth int) string {
	if maxWidth <= 1 {
		return "…"
	}
	r := []rune(text)
	if len(r) >= maxWidth {
		return string(r[:maxWidth-1]) + "…"
	}
	return text + "…"
}

func truncateWithEllipsis(text string, maxWidth int) (string, bool) {
	r := []rune(text)
	if len(r) <= maxWidth {
		return text, false
	}
	if maxWidth <= 1 {
		return "…", true
	}
	return string(r[:maxWidth-1]) + "…", true
}

	if m.ActiveTheme != nil {
		visibleText, _, _, _ = m.fitQuestionForTheme(m.ActiveTheme, visibleText)
	}
	if m.ShowHint {
		hint = q.Hint
	}
	}
	// 3. Pick New Theme based on available terminal space (re-evaluated at transition time).
	m.pickBestThemeForQuestion(m.Config.Questions[m.CurrentQuestionIndex].Text)

			if m.CurrentQuestionIndex < 0 {
				m.CurrentQuestionIndex = 0
			}
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

	oldView := ""
	if m.ActiveTheme != nil {
		oldView = m.ActiveTheme.View(m.Width, m.Height, &displayQ, m.Input.View(), hint)
	}

	// 2. Advance State
	m.CurrentQuestionIndex++
	// Wrap around for demo/endless feel, or success
	if m.CurrentQuestionIndex >= len(m.Config.Questions) {
		m.CurrentQuestionIndex = 0 // Loop for demo purposes
		// m.State = StateSuccess
		// return nil
	}

	// 3. Pick New Theme
	themeCmd := m.PickRandomTheme()

	// 4. Capture New View (Preview)
	m.Input.Reset()
	m.ShowHint = false
	m.WrongAnswers = 0
	m.TypewriterIndex = 0

	newQ := m.Config.Questions[m.CurrentQuestionIndex]
	newDisplayQ := newQ
	newDisplayQ.Text = "█"

	newView := ""
	if m.ActiveTheme != nil {
		newView = m.ActiveTheme.View(m.Width, m.Height, &newDisplayQ, m.Input.View(), "")
	}

	// 5. Create Transition
	m.State = StateTransition
	if len(transition.Registry) > 0 {
		constructor := transition.Registry[rand.Intn(len(transition.Registry))]
		m.ActiveTransition = constructor()
		m.ActiveTransition.SetContent(oldView, newView)
		return tea.Batch(themeCmd, m.ActiveTransition.Init())
	}

	m.State = StateQuestion
	return themeCmd
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
			cmds = append(cmds, tick())
		}

	case StateSuccess:
		// Success state logic
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
		if m.ActiveTheme != nil {
			visibleText, _, _, _ = m.fitQuestionForTheme(m.ActiveTheme, visibleText)
		}
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

		content := "Error: No Theme Selected"
		if m.ActiveTheme != nil {
			content = m.ActiveTheme.View(m.Width, m.Height, &displayQ, m.Input.View(), hint)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
