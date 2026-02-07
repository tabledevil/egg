package caps

import (
	"os"
	"strings"

	"github.com/muesli/termenv"
)

// Capabilities describes the detected features of the current terminal environment.
type Capabilities struct {
	// ColorProfile represents the supported color range.
	ColorProfile termenv.Profile

	// HasUnicode indicates if the terminal likely supports UTF-8 characters.
	HasUnicode bool

	// IsInteractive indicates if we are running in an interactive terminal.
	IsInteractive bool
}

// Detect analyzes the current environment and returns a Capabilities struct.
func Detect() Capabilities {
	return Capabilities{
		ColorProfile:  termenv.ColorProfile(),
		HasUnicode:    detectUnicode(),
		IsInteractive: detectInteractive(),
	}
}

// String returns a human-readable summary of the capabilities.
func (c Capabilities) String() string {
	cp := "Unknown"
	switch c.ColorProfile {
	case termenv.Ascii:
		cp = "Ascii (No Color)"
	case termenv.ANSI:
		cp = "ANSI (16 Colors)"
	case termenv.ANSI256:
		cp = "ANSI256 (256 Colors)"
	case termenv.TrueColor:
		cp = "TrueColor (16m Colors)"
	}

	uni := "No"
	if c.HasUnicode {
		uni = "Yes"
	}

	interactive := "No"
	if c.IsInteractive {
		interactive = "Yes"
	}

	return "Terminal Capabilities:\n" +
		"  Color Profile: " + cp + "\n" +
		"  Unicode Support: " + uni + "\n" +
		"  Interactive: " + interactive
}

func detectUnicode() bool {
	envVars := []string{"LANG", "LC_ALL", "LC_CTYPE"}
	for _, v := range envVars {
		val := strings.ToLower(os.Getenv(v))
		if strings.Contains(val, "utf-8") || strings.Contains(val, "utf8") {
			return true
		}
	}
	// Default to true for modern terminals unless specifically dumb
	if os.Getenv("TERM") != "dumb" {
		return true
	}
	return false
}

func detectInteractive() bool {
	// Check if stdout is a terminal
	f, _ := os.Stdout.Stat()
	return (f.Mode() & os.ModeCharDevice) != 0
}
