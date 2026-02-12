package ui

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui/theme"
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

func TestResponsiveness_QuitKeys(t *testing.T) {
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)

	tests := []tea.KeyType{tea.KeyCtrlX, tea.KeyF12}
	for _, keyType := range tests {
		key := tea.KeyMsg{Type: keyType}
		_, cmd := m.Update(key)
		if cmd == nil {
			t.Fatalf("expected tea.Quit command for key %v, got nil", keyType)
		}

		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); !ok {
			t.Fatalf("command for key %v did not produce QuitMsg", keyType)
		}
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

func TestResponsiveness_ThemeUpdateCannotDropActiveTheme(t *testing.T) {
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)
	m.State = StateQuestion
	m.ActiveTheme = theme.NewDOSTheme()

	_, _ = m.Update(game.TickMsg(time.Now()))
	if m.ActiveTheme == nil {
		t.Fatalf("active theme unexpectedly became nil after tick update")
	}
}

func TestResponsiveness_ViewDoesNotLeakInputAnsiArtifacts(t *testing.T) {
	cfg := &game.Config{
		Questions: []game.Question{{ID: 1, Text: "Q1", Answer: "A1"}},
	}
	m := NewModel(cfg)
	m.State = StateQuestion
	m.ActiveTheme = theme.NewDOSTheme()

	m.Input.SetValue("[7mT[0m\nX\rY")
	sanitized := m.themeInputValue()
	if sanitized != "T X Y" {
		t.Fatalf("unexpected sanitized input value: %q", sanitized)
	}

	m.Input.SetValue("")
	if m.themeInputValue() != "" {
		t.Fatalf("empty input should not render placeholder text")
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
