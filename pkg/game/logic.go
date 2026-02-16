package game

import (
	"strings"
	"unicode"
)

// NormalizeString is the legacy answer normalization: lowercase + trim.
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// normalizeEnhanced removes separators/punctuation while keeping meaningful
// alphanumeric symbols, so answers like "port-22" and "port 22" can match.
func normalizeEnhanced(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
		case r == '#', r == '+':
			b.WriteRune(r)
		}
	}

	return b.String()
}

func fuzzyMatch(input, correct string) bool {
	if input == correct {
		return true
	}
	if input == "" || correct == "" {
		return false
	}

	dist := LevenshteinDistance(input, correct)
	length := len([]rune(correct))

	// Keep short answers strict to avoid accidental passes.
	if length <= 3 {
		return dist == 0
	} else if length <= 6 {
		return dist <= 1
	}
	return dist <= 2
}

// LevenshteinDistance calculates the Optimal String Alignment distance (restricted Damerau-Levenshtein)
func LevenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= m; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
			// Transposition check
			if i > 1 && j > 1 && r1[i-1] == r2[j-2] && r1[i-2] == r2[j-1] {
				matrix[i][j] = min(matrix[i][j], matrix[i-2][j-2]+1)
			}
		}
	}
	return matrix[n][m]
}

func min(vals ...int) int {
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// CheckAnswer validates the user input against the correct answer with fuzzy matching
func CheckAnswer(input, correct string) bool {
	legacyInput := NormalizeString(input)
	legacyCorrect := NormalizeString(correct)
	if fuzzyMatch(legacyInput, legacyCorrect) {
		return true
	}

	enhancedInput := normalizeEnhanced(input)
	enhancedCorrect := normalizeEnhanced(correct)
	return fuzzyMatch(enhancedInput, enhancedCorrect)
}
