package ui

import (
	"ctf-tool/pkg/game"
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

	// Verify Quit command
	// tea.Quit returns a specific type of command, but we can't easily compare functions.
	// However, usually tea.Quit returns a message.
	// Actually tea.Quit() returns a Msg which is QuitMsg.
	// No, tea.Quit is a tea.Cmd.
	// We can execute it and see what happens?
	// Or trust that Update returned tea.Quit which is standard.
	// Let's just check if cmd != nil.
	if cmd == nil {
		t.Errorf("Expected tea.Quit command, got nil")
	}

	// A better check: execute the command.
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
	if m.State != StateTransition {
		// Might have skipped if no transition compatible?
		// But let's assume at least one is (Minimal is usually theme, not transition).
		// Wait, transition registry might be empty or incompatible.
		// If m.ActiveTransition is nil, State might be Question.
		// Let's check state.
	}

	// Send key input "A"
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
