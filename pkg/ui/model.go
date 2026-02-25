package ui

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/boot"
	"ctf-tool/pkg/ui/caps"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var sgrTextPattern = regexp.MustCompile(`\[[0-9;]*m`)

const transitionWatchdogTicks = 1800

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return game.TickMsg(t)
	})
}

func isDebugDumpKey(msg tea.KeyMsg) bool {
	key := strings.ToLower(strings.TrimSpace(msg.String()))
	if key == "ctrl+f5" {
		return true
	}
	if msg.Type == tea.KeyF5 || key == "f5" {
		return true
	}
	return false
}

func gameStateName(state GameState) string {
	switch state {
	case StateIntro:
		return "intro"
	case StateQuestion:
		return "question"
	case StateTransition:
		return "transition"
	case StateSuccess:
		return "success"
	default:
		return fmt.Sprintf("unknown(%d)", state)
	}
}

type Model struct {
	Config *game.Config
	State  GameState

	// Terminal capabilities used for compatibility-aware theme selection.
	Caps caps.Capabilities

	// Showcase mode cycles through themes/transitions using a stable placeholder
	// question, rather than game progression.
	Showcase bool

	// Game State
	CurrentQuestionIndex  int
	WrongAnswers          int
	ActiveBoot            boot.Intro
	ActiveTheme           theme.Theme
	ActiveTransition      transition.Transition
	BootStatus            string
	TransitionTickCount   int
	TransitionWatchdogHit bool
	DebugDumpRequested    bool
	DebugDumpTrigger      string

	// Animation State
	TypewriterIndex int
	FinaleTheme     theme.Theme

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

	showcaseThemeCursor      int
	showcaseTransitionCursor int
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
		Caps:   caps.Detect(),
		Input:  ti,
	}
	m.PickRandomBootIntro()
	m.PickRandomTheme()
	return m
}

func (m *Model) EnableShowcase() {
	m.Showcase = true
	m.AutoDemo = true
	m.ShowHint = true
	m.ActiveBoot = nil
	m.BootStatus = ""
	m.TransitionTickCount = 0
	m.TransitionWatchdogHit = false

	// Prefer a stable first question.
	if len(m.Config.Questions) > 0 {
		m.CurrentQuestionIndex = 0
		m.TypewriterIndex = len(m.Config.Questions[0].Text)
	}

	// Skip intro for quick visual inspection.
	m.State = StateQuestion

	// Provide a visible placeholder input so input fields can be evaluated too.
	m.Input.SetValue("hunter2")
}

func (m *Model) PickRandomBootIntro() tea.Cmd {
	if len(boot.Registry) == 0 {
		m.ActiveBoot = nil
		m.BootStatus = "No cinematic boot profiles are registered."
		return nil
	}

	var candidates []boot.Intro
	for _, constructor := range boot.Registry {
		intro := constructor()
		if aware, ok := intro.(boot.CapabilityAware); ok {
			if !aware.IsCompatible(m.Caps) {
				continue
			}
		}
		candidates = append(candidates, intro)
	}

	if len(candidates) == 0 {
		m.ActiveBoot = nil
		m.BootStatus = "Cinematic boot check: unsupported terminal profile, using classic startup."
		return nil
	}

	m.ActiveBoot = candidates[rand.Intn(len(candidates))]
	m.BootStatus = fmt.Sprintf("Cinematic boot check: OK - profile \"%s\"", m.ActiveBoot.Name())
	return nil
}

func (m *Model) PickRandomTheme() tea.Cmd {
	if len(theme.Registry) == 0 {
		m.ActiveTheme = nil
		return nil
	}

	var candidates []theme.Theme
	for _, constructor := range theme.Registry {
		t := constructor()
		if aware, ok := t.(theme.CapabilityAware); ok {
			if !aware.IsCompatible(m.Caps) {
				continue
			}
		}
		candidates = append(candidates, t)
	}

	// If everything opted out (e.g. very limited terminal), fall back to picking
	// something rather than failing the UI entirely.
	if len(candidates) == 0 {
		m.ActiveTheme = theme.Registry[rand.Intn(len(theme.Registry))]()
		return nil
	}

	m.ActiveTheme = candidates[rand.Intn(len(candidates))]
	return nil
}

func (m *Model) pickNextCompatibleTheme() {
	if len(theme.Registry) == 0 {
		m.ActiveTheme = nil
		return
	}

	// Try at most N constructors to find the next compatible one.
	for i := 0; i < len(theme.Registry); i++ {
		idx := (m.showcaseThemeCursor + 1 + i) % len(theme.Registry)
		t := theme.Registry[idx]()
		if aware, ok := t.(theme.CapabilityAware); ok && !aware.IsCompatible(m.Caps) {
			continue
		}
		m.ActiveTheme = t
		m.showcaseThemeCursor = idx
		return
	}

	// If everything opted out, pick something anyway.
	m.ActiveTheme = theme.Registry[(m.showcaseThemeCursor+1)%len(theme.Registry)]()
	m.showcaseThemeCursor = (m.showcaseThemeCursor + 1) % len(theme.Registry)
}

func (m *Model) PickNextTheme() tea.Cmd {
	m.pickNextCompatibleTheme()
	return nil
}

func (m *Model) nextCompatibleTransition() transition.Transition {
	if len(transition.Registry) == 0 {
		return nil
	}

	for i := 0; i < len(transition.Registry); i++ {
		idx := (m.showcaseTransitionCursor + i) % len(transition.Registry)
		candidate := transition.Registry[idx]()
		if aware, ok := candidate.(transition.CapabilityAware); ok && !aware.IsCompatible(m.Caps) {
			continue
		}
		m.showcaseTransitionCursor = (idx + 1) % len(transition.Registry)
		return candidate
	}

	idx := m.showcaseTransitionCursor % len(transition.Registry)
	m.showcaseTransitionCursor = (idx + 1) % len(transition.Registry)
	return transition.Registry[idx]()
}

func (m *Model) randomCompatibleTransition() transition.Transition {
	if len(transition.Registry) == 0 {
		return nil
	}

	var constructors []transition.Constructor
	for _, constructor := range transition.Registry {
		candidate := constructor()
		if aware, ok := candidate.(transition.CapabilityAware); ok && !aware.IsCompatible(m.Caps) {
			continue
		}
		constructors = append(constructors, constructor)
	}

	if len(constructors) == 0 {
		return transition.Registry[rand.Intn(len(transition.Registry))]()
	}

	return constructors[rand.Intn(len(constructors))]()
}

func (m *Model) StartShowcaseTransition() tea.Cmd {
	if len(m.Config.Questions) == 0 {
		return nil
	}

	// 1. Capture Old View (fully visible text).
	q := m.Config.Questions[m.CurrentQuestionIndex]
	displayQ := q
	displayQ.Text = q.Text
	hint := ""
	if m.ShowHint {
		hint = q.Hint
	}
	oldView := m.safeThemeView(&displayQ, m.themeInputValue(), hint)

	// 2. Pick next compatible theme.
	m.pickNextCompatibleTheme()

	// 3. Capture New View (same question, different theme).
	newView := m.safeThemeView(&displayQ, m.themeInputValue(), hint)

	// 4. Create next transition (sequential).
	m.State = StateTransition
	m.TransitionTickCount = 0
	m.TransitionWatchdogHit = false
	if len(transition.Registry) > 0 {
		m.ActiveTransition = m.nextCompatibleTransition()
		if m.ActiveTransition == nil {
			m.State = StateQuestion
			return nil
		}
		m.ActiveTransition.SetContent(oldView, newView)
		return m.ActiveTransition.Init()
	}

	m.State = StateQuestion
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
	if m.ShowHint {
		hint = q.Hint
	}

	oldView := m.safeThemeView(&displayQ, m.themeInputValue(), hint)

	// 2. Advance State
	m.CurrentQuestionIndex++
	// Check for game completion
	if m.CurrentQuestionIndex >= len(m.Config.Questions) {
		m.State = StateSuccess
		m.Input.Reset()

		// Initialize Finale theme if available
		m.FinaleTheme = nil
		for _, constructor := range theme.Registry {
			tmp := constructor()
			if tmp.Name() == "Antigravity (Finale)" {
				m.FinaleTheme = tmp
				if initCmd := m.FinaleTheme.Init(); initCmd != nil {
					return initCmd
				}
				break
			}
		}

		return nil
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
	m.TransitionTickCount = 0
	m.TransitionWatchdogHit = false
	if len(transition.Registry) > 0 {
		m.ActiveTransition = m.randomCompatibleTransition()
		if m.ActiveTransition == nil {
			m.State = StateQuestion
			return nil
		}
		m.ActiveTransition.SetContent(oldView, newView)
		return m.ActiveTransition.Init()
	}

	m.State = StateQuestion
	return nil
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, tick())
	if m.ActiveBoot != nil {
		cmds = append(cmds, m.ActiveBoot.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if isDebugDumpKey(msg) {
			m.DebugDumpRequested = true
			trigger := strings.TrimSpace(msg.String())
			if trigger == "" {
				trigger = fmt.Sprintf("%v", msg.Type)
			}
			m.DebugDumpTrigger = trigger
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlX, tea.KeyEsc, tea.KeyF12:
			return m, tea.Quit
		case tea.KeyF1:
			// Theme cycling
			if m.Showcase {
				cmds = append(cmds, m.PickNextTheme())
			} else {
				cmds = append(cmds, m.PickRandomTheme())
			}
		case tea.KeyF2:
			// Transition cycling
			if m.Showcase {
				cmds = append(cmds, m.StartShowcaseTransition())
			} else {
				// Force Transition (stay on same Q)
				m.CurrentQuestionIndex--
				if m.CurrentQuestionIndex < 0 {
					m.CurrentQuestionIndex = 0
				}
				cmds = append(cmds, m.StartTransition())
			}
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
					if m.Showcase {
						cmds = append(cmds, m.StartShowcaseTransition())
					} else {
						cmds = append(cmds, m.StartTransition())
					}
				}
			}
		}
	}

	// State Machine
	switch m.State {
	case StateIntro:
		if m.ActiveBoot != nil {
			var bCmd tea.Cmd
			nextBoot, bCmd := m.ActiveBoot.Update(msg)
			if nextBoot != nil {
				m.ActiveBoot = nextBoot
			}
			cmds = append(cmds, bCmd)
		}

		if _, ok := msg.(game.TickMsg); ok {
			if m.ActiveBoot == nil {
				cmds = append(cmds, tick())
			}
		}
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			m.State = StateQuestion
			m.TypewriterIndex = 0
		}

	case StateTransition:
		if _, ok := msg.(game.TickMsg); ok {
			m.TransitionTickCount++
			if m.TransitionTickCount > transitionWatchdogTicks {
				m.TransitionWatchdogHit = true
				m.ActiveTransition = nil
				m.State = StateQuestion
				m.TypewriterIndex = 0
				break
			}
		}

		if m.ActiveTransition != nil {
			var tCmd tea.Cmd
			nextTransition, tCmd := m.ActiveTransition.Update(msg)
			if nextTransition != nil {
				m.ActiveTransition = nextTransition
			}
			cmds = append(cmds, tCmd)

			if m.ActiveTransition.Done() {
				m.ActiveTransition = nil
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
				m.Input.Reset()
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
		if m.FinaleTheme != nil {
			var tCmd tea.Cmd
			nextTheme, tCmd := m.FinaleTheme.Update(msg)
			if nextTheme != nil {
				m.FinaleTheme = nextTheme
			}
			cmds = append(cmds, tCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Width == 0 {
		return "Loading..."
	}

	switch m.State {
	case StateIntro:
		if m.ActiveBoot != nil {
			if bootView := m.safeBootView(); bootView != "" {
				statusStyle := lipgloss.NewStyle().
					Width(m.Width).
					Align(lipgloss.Center).
					Foreground(lipgloss.Color("#8DF7D9"))

				action := "[PRESS ENTER TO SKIP BOOT]"
				if m.ActiveBoot.Done() {
					action = "[PRESS ENTER TO CONTINUE]"
				}

				footer := statusStyle.Render(fmt.Sprintf("%s\n%s", m.BootStatus, action))
				return lipgloss.JoinVertical(lipgloss.Left, bootView, footer)
			}
		}

		style := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#00FF00"))

		classic := "SYSTEM BOOT SEQUENCE INITIATED...\n\n[PRESS ENTER TO HACK THE PLANET]"
		if m.BootStatus != "" {
			classic += "\n\n" + m.BootStatus
		}
		return style.Render(classic)

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
		// Draw finale background if available
		bg := ""
		if m.FinaleTheme != nil {
			// Mock a question so that the Avoidance layout has bounding boxes to avoid
			mockQ := &game.Question{
				Text: m.Config.FinalMessage,
			}
			bg = m.FinaleTheme.View(m.Width, m.Height, mockQ, m.Config.FinalHint, "")
		} else {
			bg = lipgloss.NewStyle().Width(m.Width).Height(m.Height).Render("")
		}

		style := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#00FF00")).
			Padding(2)

		// The Avoidance theme renders the text itself, so we don't strictly need to overlay it
		// here if the theme did it, but the theme expects standard hints.
		// Since we mocked the question and hint, the Finale theme will just draw it.
		if m.FinaleTheme != nil {
			return bg
		}

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

func (m *Model) safeBootView() (view string) {
	if m.ActiveBoot == nil || m.Width <= 0 || m.Height <= 0 {
		return ""
	}

	height := m.Height - 2
	if height < 3 {
		height = m.Height
	}

	defer func() {
		if recover() != nil {
			view = ""
		}
	}()

	return m.ActiveBoot.View(m.Width, height)
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

func typeName(v any) string {
	if v == nil {
		return "(none)"
	}
	return fmt.Sprintf("%T", v)
}

func trimForDebug(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func (m Model) DebugSnapshot() string {
	inputPreview := ansi.Strip(m.Input.Value())
	inputPreview = sgrTextPattern.ReplaceAllString(inputPreview, "")
	inputPreview = strings.ReplaceAll(inputPreview, "\n", " ")
	inputPreview = strings.ReplaceAll(inputPreview, "\r", " ")

	activeThemeName := "(none)"
	if m.ActiveTheme != nil {
		activeThemeName = m.ActiveTheme.Name()
	}

	activeBootName := "(none)"
	activeBootDone := false
	if m.ActiveBoot != nil {
		activeBootName = m.ActiveBoot.Name()
		activeBootDone = m.ActiveBoot.Done()
	}

	qID := -1
	qText := ""
	if m.Config != nil && m.CurrentQuestionIndex >= 0 && m.CurrentQuestionIndex < len(m.Config.Questions) {
		q := m.Config.Questions[m.CurrentQuestionIndex]
		qID = q.ID
		qText = trimForDebug(q.Text, 96)
	}

	var b strings.Builder
	b.WriteString("=== EGG DEBUG SNAPSHOT ===\n")
	b.WriteString("timestamp: ")
	b.WriteString(time.Now().Format(time.RFC3339Nano))
	b.WriteString("\n")
	b.WriteString("trigger: ")
	b.WriteString(trimForDebug(m.DebugDumpTrigger, 64))
	b.WriteString("\n")
	b.WriteString("state: ")
	b.WriteString(gameStateName(m.State))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("terminal: width=%d height=%d\n", m.Width, m.Height))
	b.WriteString(fmt.Sprintf("caps: color=%v unicode=%t interactive=%t\n", m.Caps.ColorProfile, m.Caps.HasUnicode, m.Caps.IsInteractive))
	b.WriteString(fmt.Sprintf("boot: name=%q done=%t status=%q\n", activeBootName, activeBootDone, m.BootStatus))
	b.WriteString(fmt.Sprintf("theme: name=%q type=%s\n", activeThemeName, typeName(m.ActiveTheme)))
	b.WriteString(fmt.Sprintf("transition: type=%s ticks=%d watchdog_hit=%t watchdog_limit=%d\n", typeName(m.ActiveTransition), m.TransitionTickCount, m.TransitionWatchdogHit, transitionWatchdogTicks))
	b.WriteString(fmt.Sprintf("progress: question_index=%d question_id=%d wrong_answers=%d hint_visible=%t typewriter_index=%d\n", m.CurrentQuestionIndex, qID, m.WrongAnswers, m.ShowHint, m.TypewriterIndex))
	b.WriteString(fmt.Sprintf("question_preview: %q\n", qText))
	b.WriteString(fmt.Sprintf("input: len=%d value=%q\n", len([]rune(inputPreview)), trimForDebug(inputPreview, 96)))
	b.WriteString(fmt.Sprintf("modes: showcase=%t auto_demo=%t\n", m.Showcase, m.AutoDemo))
	b.WriteString("=== END SNAPSHOT ===")

	return b.String()
}
