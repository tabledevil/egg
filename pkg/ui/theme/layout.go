package theme

import "strings"

type layoutBox struct {
	x int
	y int
	w int
	h int
}

func boundedSpan(total, margin, minSpan, maxSpan int) int {
	if total <= 0 {
		return 1
	}
	if margin < 0 {
		margin = 0
	}
	if minSpan < 1 {
		minSpan = 1
	}

	span := total - (margin * 2)
	if span < 1 {
		span = 1
	}
	if maxSpan > 0 && span > maxSpan {
		span = maxSpan
	}
	if span < minSpan {
		span = minSpan
	}
	if span > total {
		span = total
	}
	if span < 1 {
		span = 1
	}
	return span
}

func centeredStart(total, span int) int {
	if total <= 0 || span >= total {
		return 0
	}
	if span < 1 {
		span = 1
	}
	return (total - span) / 2
}

func centeredBox(totalWidth, totalHeight, boxWidth, boxHeight int) layoutBox {
	if totalWidth <= 0 {
		totalWidth = 1
	}
	if totalHeight <= 0 {
		totalHeight = 1
	}
	if boxWidth < 1 {
		boxWidth = 1
	}
	if boxHeight < 1 {
		boxHeight = 1
	}
	if boxWidth > totalWidth {
		boxWidth = totalWidth
	}
	if boxHeight > totalHeight {
		boxHeight = totalHeight
	}

	return layoutBox{
		x: centeredStart(totalWidth, boxWidth),
		y: centeredStart(totalHeight, boxHeight),
		w: boxWidth,
		h: boxHeight,
	}
}

func (b layoutBox) inset(padX, padY int) layoutBox {
	if padX < 0 {
		padX = 0
	}
	if padY < 0 {
		padY = 0
	}

	inner := layoutBox{
		x: b.x + padX,
		y: b.y + padY,
		w: b.w - (padX * 2),
		h: b.h - (padY * 2),
	}

	if inner.w < 1 {
		inner.w = 1
	}
	if inner.h < 1 {
		inner.h = 1
	}

	return inner
}

func wrapAndClamp(prefix, text string, width, maxLines int) []string {
	if width <= 0 || maxLines <= 0 {
		return nil
	}

	var lines []string
	if prefix == "" {
		lines = wrapText(text, width)
	} else {
		lines = wrapLabeled(prefix, text, width)
	}

	return clampLines(lines, maxLines, width)
}

func remainingRows(totalRows, usedRows, minRows int) int {
	if totalRows <= 0 {
		return 0
	}
	if minRows < 0 {
		minRows = 0
	}

	remaining := totalRows - usedRows
	if remaining < minRows {
		remaining = minRows
	}
	if remaining > totalRows {
		remaining = totalRows
	}
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

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
