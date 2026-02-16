package theme

import "strings"

func runeLen(s string) int {
	return len([]rune(s))
}

func truncateToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	return string(r[:width])
}

func truncateWithEllipsis(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}

func appendEllipsisWithinWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) >= width {
		if width == 1 {
			return "…"
		}
		return string(r[:width-1]) + "…"
	}
	return s + "…"
}

func splitLongToken(token string, width int) []string {
	if width <= 0 {
		return nil
	}
	r := []rune(token)
	if len(r) == 0 {
		return []string{""}
	}

	parts := make([]string, 0, (len(r)/width)+1)
	for start := 0; start < len(r); start += width {
		end := start + width
		if end > len(r) {
			end = len(r)
		}
		parts = append(parts, string(r[start:end]))
	}
	return parts
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return nil
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	paragraphs := strings.Split(text, "\n")
	var out []string

	for i, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}

		line := ""
		for _, word := range words {
			if line == "" {
				if runeLen(word) <= width {
					line = word
					continue
				}

				parts := splitLongToken(word, width)
				if len(parts) == 0 {
					continue
				}
				if len(parts) > 1 {
					out = append(out, parts[:len(parts)-1]...)
				}
				line = parts[len(parts)-1]
				continue
			}

			candidate := line + " " + word
			if runeLen(candidate) <= width {
				line = candidate
				continue
			}

			out = append(out, line)
			if runeLen(word) <= width {
				line = word
				continue
			}

			parts := splitLongToken(word, width)
			if len(parts) == 0 {
				line = ""
				continue
			}
			if len(parts) > 1 {
				out = append(out, parts[:len(parts)-1]...)
			}
			line = parts[len(parts)-1]
		}

		if line != "" {
			out = append(out, line)
		}

		if i < len(paragraphs)-1 && paragraph == "" {
			continue
		}
	}

	if len(out) == 0 {
		return []string{""}
	}

	return out
}

func wrapLabeled(prefix, text string, width int) []string {
	if width <= 0 {
		return nil
	}

	prefixLen := runeLen(prefix)
	if prefixLen >= width {
		return []string{truncateWithEllipsis(prefix, width)}
	}

	contentWidth := width - prefixLen
	wrapped := wrapText(text, contentWidth)
	if len(wrapped) == 0 {
		return []string{truncateToWidth(prefix, width)}
	}

	indent := strings.Repeat(" ", prefixLen)
	out := make([]string, 0, len(wrapped))
	out = append(out, truncateToWidth(prefix+wrapped[0], width))
	for _, line := range wrapped[1:] {
		out = append(out, truncateToWidth(indent+line, width))
	}
	return out
}

func clampLines(lines []string, maxLines, width int) []string {
	if maxLines <= 0 || width <= 0 {
		return nil
	}
	if len(lines) == 0 {
		return []string{}
	}

	limit := maxLines
	if len(lines) < limit {
		limit = len(lines)
	}

	out := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, truncateToWidth(lines[i], width))
	}

	if len(lines) > maxLines && len(out) > 0 {
		out[len(out)-1] = appendEllipsisWithinWidth(out[len(out)-1], width)
	}

	return out
}

func sliceRunes(s string, start, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) == 0 {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if start >= len(r) {
		return ""
	}
	end := start + width
	if end > len(r) {
		end = len(r)
	}
	return string(r[start:end])
}
