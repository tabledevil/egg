package main

import (
	"ctf-tool/pkg/game"
	"ctf-tool/pkg/ui"
	"ctf-tool/pkg/ui/boot"
	"ctf-tool/pkg/ui/caps"
	"ctf-tool/pkg/ui/theme"
	"ctf-tool/pkg/ui/transition"
	"ctf-tool/pkg/web"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	showcase := flag.Bool("showcase", false, "run UI showcase mode (cycles themes/transitions with placeholder text)")
	list := flag.Bool("list", false, "list supported boot profiles, themes, and transitions for this terminal")
	webMode := flag.Bool("web", false, "serve the CTF tool as a web terminal instead of running in the current terminal")
	port := flag.Int("port", 8080, "port for the web terminal server (used with -web)")
	flag.Parse()

	// --- Web terminal mode ---
	if *webMode {
		self, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot resolve own binary: %v\n", err)
			os.Exit(1)
		}
		// Build args to pass to the child process (everything except -web/-port).
		var childArgs []string
		if *showcase {
			childArgs = append(childArgs, "-showcase")
		}
		addr := fmt.Sprintf(":%d", *port)
		if err := web.Serve(addr, self, childArgs); err != nil {
			fmt.Fprintf(os.Stderr, "web server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *list {
		c := caps.Detect()
		fmt.Println(c.String())
		fmt.Println()

		fmt.Println("Boot profiles:")
		bootSupported := 0
		for _, constructor := range boot.Registry {
			b := constructor()
			supported := true
			if aware, ok := b.(boot.CapabilityAware); ok && !aware.IsCompatible(c) {
				supported = false
			}
			status := "supported"
			if !supported {
				status = "unsupported"
			}
			fmt.Printf("- %s: %s [%s]\n", b.Name(), b.Description(), status)
			if supported {
				bootSupported++
			}
		}
		if bootSupported == 0 {
			fmt.Println("- (none; classic startup intro will be used)")
		}
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

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(ui.Model); ok && m.DebugDumpRequested {
		fmt.Println()
		fmt.Println(m.DebugSnapshot())
	}
}
