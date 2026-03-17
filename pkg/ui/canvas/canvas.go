package canvas

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// defaultStyle is a shared blank style reused for all empty cells,
// avoiding repeated lipgloss.NewStyle() allocations.
var defaultStyle = lipgloss.NewStyle()

// defaultStyleID is the reserved registry ID for the default (blank) style.
const defaultStyleID = 0

// Cell represents a single character on the screen with its style
type Cell struct {
	Rune    rune
	Style   lipgloss.Style
	styleID int // assigned by Canvas.registerStyle for fast equality comparison
}

// Canvas represents a 2D grid of cells
type Canvas struct {
	Width  int
	Height int
	Grid   [][]Cell

	// Style registry: maps a style fingerprint to a compact integer ID so that
	// Render() can coalesce adjacent same-styled cells with a simple int comparison.
	styles   []lipgloss.Style
	styleMap map[string]int
}

// New creates a new canvas of the given dimensions
func New(width, height int) *Canvas {
	c := &Canvas{
		Width:  width,
		Height: height,
		Grid:   make([][]Cell, height),
		styles: []lipgloss.Style{defaultStyle},
		styleMap: map[string]int{
			defaultStyle.Render("X"): defaultStyleID,
		},
	}
	blank := Cell{Rune: ' ', Style: defaultStyle, styleID: defaultStyleID}
	for y := 0; y < height; y++ {
		c.Grid[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			c.Grid[y][x] = blank
		}
	}
	return c
}

// Clear resets the canvas to spaces
func (c *Canvas) Clear() {
	blank := Cell{Rune: ' ', Style: defaultStyle, styleID: defaultStyleID}
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Grid[y][x] = blank
		}
	}
}

// registerStyle returns a stable integer ID for the given style within this
// canvas frame.  Styles that produce identical ANSI output share the same ID.
func (c *Canvas) registerStyle(s lipgloss.Style) int {
	key := s.Render("X")
	if id, ok := c.styleMap[key]; ok {
		return id
	}
	id := len(c.styles)
	c.styles = append(c.styles, s)
	c.styleMap[key] = id
	return id
}

// SetChar sets a single character at the given coordinates
func (c *Canvas) SetChar(x, y int, char rune, style lipgloss.Style) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Grid[y][x] = Cell{Rune: char, Style: style, styleID: c.registerStyle(style)}
	}
}

// SetString writes a string starting at the given coordinates
func (c *Canvas) SetString(x, y int, text string, style lipgloss.Style) {
	if y < 0 || y >= c.Height {
		return
	}

	sid := c.registerStyle(style)

	startX := x
	col := x
	row := y
	for _, char := range text {
		switch char {
		case '\r':
			col = startX
			continue
		case '\n':
			row++
			col = startX
			if row >= c.Height {
				return
			}
			continue
		}

		if col >= 0 && col < c.Width && row >= 0 && row < c.Height {
			c.Grid[row][col] = Cell{Rune: char, Style: style, styleID: sid}
		}
		col++
	}
}

// Render converts the grid to a string for Bubble Tea.
//
// Run-length coalescing: consecutive cells on the same row that share the same
// styleID are batched into a single lipgloss.Render() call.  This typically
// reduces the number of ANSI escape sequences by 60-90%.
func (c *Canvas) Render() string {
	var b strings.Builder
	b.Grow(c.Width * c.Height * 2) // conservative pre-alloc

	var run strings.Builder
	run.Grow(c.Width)

	for y := 0; y < c.Height; y++ {
		if c.Width == 0 {
			if y < c.Height-1 {
				b.WriteRune('\n')
			}
			continue
		}

		row := c.Grid[y]
		runID := row[0].styleID
		run.Reset()
		run.WriteRune(row[0].Rune)

		for x := 1; x < c.Width; x++ {
			cell := row[x]
			if cell.styleID == runID {
				run.WriteRune(cell.Rune)
			} else {
				// Flush current run
				if runID == defaultStyleID {
					b.WriteString(run.String())
				} else {
					b.WriteString(c.styles[runID].Render(run.String()))
				}
				run.Reset()
				run.WriteRune(cell.Rune)
				runID = cell.styleID
			}
		}
		// Flush final run of the row
		if runID == defaultStyleID {
			b.WriteString(run.String())
		} else {
			b.WriteString(c.styles[runID].Render(run.String()))
		}

		if y < c.Height-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// DrawBox draws a border box
func (c *Canvas) DrawBox(x, y, w, h int, style lipgloss.Style) {
	// Top/Bottom
	for i := 0; i < w; i++ {
		c.SetChar(x+i, y, '─', style)
		c.SetChar(x+i, y+h-1, '─', style)
	}
	// Sides
	for i := 0; i < h; i++ {
		c.SetChar(x, y+i, '│', style)
		c.SetChar(x+w-1, y+i, '│', style)
	}
	// Corners
	c.SetChar(x, y, '┌', style)
	c.SetChar(x+w-1, y, '┐', style)
	c.SetChar(x, y+h-1, '└', style)
	c.SetChar(x+w-1, y+h-1, '┘', style)
}

// Fill fills a rectangle with a character
func (c *Canvas) Fill(x, y, w, h int, char rune, style lipgloss.Style) {
	sid := c.registerStyle(style)
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < c.Width && py >= 0 && py < c.Height {
				c.Grid[py][px] = Cell{Rune: char, Style: style, styleID: sid}
			}
		}
	}
}
