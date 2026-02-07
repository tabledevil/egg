package main

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui"
	"ctf-tool/pkg/ui/caps"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--profile" || os.Args[1] == "-p") {
		runProfiler()
		return
	}

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

func runProfiler() {
	c := caps.Detect()
	fmt.Println(lipgloss.NewStyle().Bold(true).Render("=== TERMINAL PROFILER ==="))
	fmt.Println(c.String())
	fmt.Println()

	fmt.Println(lipgloss.NewStyle().Underline(true).Render("Checking Themes:"))
	passCount := 0
	failCount := 0

	for _, constructor := range theme.Registry {
		th := constructor()
		name := th.Name()
		compatible := th.IsCompatible(c)
		status := " [PASS] "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

		if !compatible {
			status = " [FAIL] "
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
			failCount++
		} else {
			passCount++
		}

		fmt.Printf("%s %s\n", style.Render(status), name)
		if !compatible {
			fmt.Printf("    Reason: Incompatible with current capabilities (ColorProfile=%v)\n", c.ColorProfile)
		}
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Underline(true).Render("Checking Transitions:"))

	for _, constructor := range transition.Registry {
		tr := constructor()
		compatible := tr.IsCompatible(c)
		status := " [PASS] "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

		if !compatible {
			status = " [FAIL] "
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		}

		// Transitions don't have names in interface, try to identify by type or generic name
		fmt.Printf("%s Transition\n", style.Render(status))
	}

	fmt.Println()
	fmt.Printf("Summary: %d Compatible, %d Incompatible\n", passCount, failCount)
	if failCount > 0 {
		fmt.Println("Some themes are disabled in this environment.")
	} else {
		fmt.Println("All systems operational.")
	}
}
