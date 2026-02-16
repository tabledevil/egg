package main

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui"
	"ctf-tool/pkg/ui/caps"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	showcase := flag.Bool("showcase", false, "run UI showcase mode (cycles themes/transitions with placeholder text)")
	list := flag.Bool("list", false, "list supported themes and transitions for the current terminal and exit")
	flag.Parse()

	if *list {
		c := caps.Detect()
		fmt.Println(c.String())
		fmt.Println()

		fmt.Println("Supported themes:")
		themesSupported := 0
		for _, constructor := range theme.Registry {
			t := constructor()
			if aware, ok := t.(theme.CapabilityAware); ok && !aware.IsCompatible(c) {
				continue
			}
			fmt.Printf("- %s: %s\n", t.Name(), t.Description())
			themesSupported++
		}
		if themesSupported == 0 {
			fmt.Println("- (none)")
		}
		fmt.Println()

		fmt.Println("Supported transitions:")
		transSupported := 0
		for _, constructor := range transition.Registry {
			tr := constructor()
			if aware, ok := tr.(transition.CapabilityAware); ok && !aware.IsCompatible(c) {
				continue
			}
			fmt.Printf("- %T\n", tr)
			transSupported++
		}
		if transSupported == 0 {
			fmt.Println("- (none)")
		}
		return
	}

	var config *game.Config
	if *showcase {
		config = &game.Config{
			Questions: []game.Question{
				{
					ID:     1,
					Text:   "Showcase prompt: The quick brown fox jumps over the lazy dog. 0123456789. Symbols: !@#$%^&*()_+-=[]{};':,./<>?",
					Answer: "ok",
					Hint:   "This is a placeholder hint to validate hint rendering.",
				},
			},
			FinalMessage: "SHOWCASE COMPLETE",
			FinalHint:    "Press F1/F2 to cycle, F3 for auto-demo, Ctrl+X/F12 to quit.",
		}
	} else {
		// Load configuration (embedded)
		var err error
		config, err = game.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading game data: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize UI
	model := ui.NewModel(config)
	if *showcase {
		model.EnableShowcase()
	}
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
