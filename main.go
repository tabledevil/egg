package main

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	// Load configuration (embedded)
	config, err := game.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading game data: %v\n", err)
		os.Exit(1)
	}

	// Initialize UI
	model := ui.NewModel(config)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
