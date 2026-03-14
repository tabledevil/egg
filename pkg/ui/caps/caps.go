package caps

import (
	"os"
	"strings"
	"sync"
	"time"

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

	// LowBandwidthTTY indicates a TTY session that benefits from fewer redraws,
	// e.g. ttyd/web terminals or explicit low-bandwidth mode.
	LowBandwidthTTY bool
}

// Detect analyzes the current environment and returns a Capabilities struct.
func Detect() Capabilities {
	return Capabilities{
		ColorProfile:    termenv.ColorProfile(),
		HasUnicode:      detectUnicode(),
		IsInteractive:   detectInteractive(),
		LowBandwidthTTY: detectLowBandwidthTTY(),
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

	lowBandwidth := "No"
	if c.LowBandwidthTTY {
		lowBandwidth = "Yes"
	}

	return "Terminal Capabilities:\n" +
		"  Color Profile: " + cp + "\n" +
		"  Unicode Support: " + uni + "\n" +
		"  Interactive: " + interactive + "\n" +
		"  Low Bandwidth Mode: " + lowBandwidth
}

var lowBandwidthOnce sync.Once
var lowBandwidthEnabled bool

// AnimationFrameInterval returns the preferred frame interval for animations.
func AnimationFrameInterval() time.Duration {
	if detectLowBandwidthTTY() {
		return 55 * time.Millisecond
	}
	return 33 * time.Millisecond
}

// BootFrameInterval returns the preferred frame interval for intro animations.
func BootFrameInterval() time.Duration {
	if detectLowBandwidthTTY() {
		return 66 * time.Millisecond
	}
	return 45 * time.Millisecond
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

func detectLowBandwidthTTY() bool {
	lowBandwidthOnce.Do(func() {
		for _, key := range []string{"EGG_LOW_BANDWIDTH", "LOW_BANDWIDTH_TTY"} {
			if envTruthy(os.Getenv(key)) {
				lowBandwidthEnabled = true
				return
			}
		}

		for _, key := range []string{"TTYD", "WETTY", "WEBTTY", "WEB_TERMINAL", "TT_REPLAY"} {
			if os.Getenv(key) != "" {
				lowBandwidthEnabled = true
				return
			}
		}

		for _, key := range []string{"TERMINAL_EMULATOR", "TERM_PROGRAM"} {
			value := strings.ToLower(os.Getenv(key))
			if strings.Contains(value, "ttyd") || strings.Contains(value, "wetty") || strings.Contains(value, "web") {
				lowBandwidthEnabled = true
				return
			}
		}
	})

	return lowBandwidthEnabled
}

func envTruthy(value string) bool {
	v := strings.TrimSpace(strings.ToLower(value))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
