package game

import "testing"

func TestCheckAnswer(t *testing.T) {
	tests := []struct {
		input    string
		correct  string
		expected bool
	}{
		{"echo", "echo", true},
		{"ECHO", "echo", true},
		{"  echo  ", "echo", true},
		{"ecoh", "echo", true}, // 1 swap, length 4 -> allowed? length 4 > 3, allows 1 typo. Yes.
		{"eco", "echo", true},  // 1 deletion, length 4 -> allowed.
		{"cat", "cat", true},
		{"bat", "cat", false}, // length 3, must be exact
		{"keyboard", "keyboard", true},
		{"keybaord", "keyboard", true}, // 1 swap
		{"keybard", "keyboard", true},  // 1 deletion
		{"koyboard", "keyboard", true}, // 1 sub
		{"keyboar", "keyboard", true},  // 1 deletion
		{"somethingreallylong", "somethingreallylong", true},
		{"somethingrealllong", "somethingreallylong", true}, // 1 typo
		{"somethngrealllong", "somethingreallylong", true},  // 2 typos
	}

	for _, tt := range tests {
		if got := CheckAnswer(tt.input, tt.correct); got != tt.expected {
			t.Errorf("CheckAnswer(%q, %q) = %v, want %v", tt.input, tt.correct, got, tt.expected)
		}
	}
}
