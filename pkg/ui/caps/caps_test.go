package caps

import (
	"strings"
	"testing"
	"time"
)

func TestAnimationFrameIntervalBounds(t *testing.T) {
	interval := AnimationFrameInterval()
	if interval != 33*time.Millisecond && interval != 55*time.Millisecond {
		t.Fatalf("unexpected animation interval: %v", interval)
	}
}

func TestBootFrameIntervalBounds(t *testing.T) {
	interval := BootFrameInterval()
	if interval != 45*time.Millisecond && interval != 66*time.Millisecond {
		t.Fatalf("unexpected boot interval: %v", interval)
	}
}

func TestCapabilitiesStringContainsBandwidthLine(t *testing.T) {
	c := Capabilities{LowBandwidthTTY: true}
	if got := c.String(); got == "" || !strings.Contains(got, "Low Bandwidth Mode") {
		t.Fatalf("capabilities string missing low bandwidth line: %q", got)
	}
}
