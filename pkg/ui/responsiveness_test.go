package ui

import (
	"ctf-tool/pkg/game"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestResponsiveness_CtrlC(t *testing.T) {
	// Setup Model
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)

	// Send Ctrl+C
	ctrlC := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(ctrlC)

	if cmd == nil {
		t.Errorf("Expected tea.Quit command, got nil")
	}

	// Execute command to verify it is tea.Quit
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("Command did not produce QuitMsg")
	}
}

func TestResponsiveness_InputDuringTransition(t *testing.T) {
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)

	// Force transition state
	m.StartTransition()

	// If transitions are registered and compatible, we should be in StateTransition
	if m.ActiveTransition != nil && m.State != StateTransition {
		t.Errorf("Expected StateTransition, got %v", m.State)
	}

	// Send key input "A" - this should be processed immediately and non-blocking
	keyA := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")}
	start := time.Now()
	_, _ = m.Update(keyA)
	duration := time.Since(start)

	if duration > 10*time.Millisecond {
		t.Errorf("Update took too long (%v), potential blocking detected", duration)
	}
}

func TestModelInit(t *testing.T) {
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)
	cmd := m.Init()
	if cmd == nil {
		t.Errorf("Init should return initial commands")
	}
}

func TestQuestionFormatting_SingleLineThemesUseEllipsis(t *testing.T) {
	cfg := &game.Config{Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}}}
	m := NewModel(cfg)
	m.Width = 40
	m.Height = 20
	m.ActiveTheme = theme.NewGhostInShellTheme()

	formatted, _, truncated, layout := m.fitQuestionForTheme(m.ActiveTheme, "this is a very long question that definitely exceeds one line")
	if layout.multiline {
		t.Fatalf("expected one-line layout for ghost theme")
	}
	if !truncated {
		t.Fatalf("expected truncation for one-line theme")
	}
	if formatted == "" || []rune(formatted)[len([]rune(formatted))-1] != '…' {
		t.Fatalf("expected ellipsis suffix, got %q", formatted)
	}
}

func TestQuestionFormatting_MultilineWraps(t *testing.T) {
	cfg := &game.Config{Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}}}
	m := NewModel(cfg)
	m.Width = 44
	m.Height = 20
	m.ActiveTheme = theme.NewDOSTheme()

	formatted, _, truncated, layout := m.fitQuestionForTheme(m.ActiveTheme, "one two three four five six seven eight nine ten eleven")
	if !layout.multiline {
		t.Fatalf("expected multiline layout for DOS")
	}
	if !strings.Contains(formatted, "\n") {
		t.Fatalf("expected wrapped output with newline, got %q", formatted)
	}
	if truncated {
		t.Fatalf("did not expect truncation for this prompt")
	}
}
