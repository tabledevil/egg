package canvas

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cell represents a single character on the screen with its style
type Cell struct {
	Rune  rune
	Style lipgloss.Style
}

// Canvas represents a 2D grid of cells
type Canvas struct {
	Width  int
	Height int
	Grid   [][]Cell
}

// New creates a new canvas of the given dimensions
func New(width, height int) *Canvas {
	grid := make([][]Cell, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			grid[y][x] = Cell{Rune: ' ', Style: lipgloss.NewStyle()}
		}
	}
	return &Canvas{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

// Clear resets the canvas to spaces
func (c *Canvas) Clear() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Grid[y][x] = Cell{Rune: ' ', Style: lipgloss.NewStyle()}
		}
	}
}

// SetChar sets a single character at the given coordinates
func (c *Canvas) SetChar(x, y int, char rune, style lipgloss.Style) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Grid[y][x] = Cell{Rune: char, Style: style}
	}
}

// SetString writes a string starting at the given coordinates
func (c *Canvas) SetString(x, y int, text string, style lipgloss.Style) {
	if y < 0 || y >= c.Height {
		return
	}

	col := x
	for _, char := range text {
		if col >= 0 && col < c.Width {
			c.Grid[y][col] = Cell{Rune: char, Style: style}
		}
		col++
	}
}

// Render converts the grid to a string for Bubble Tea
func (c *Canvas) Render() string {
	var b strings.Builder

	for y := 0; y < c.Height; y++ {
		rowBuilder := strings.Builder{}
		for x := 0; x < c.Width; x++ {
			cell := c.Grid[y][x]
			// Optimization: If style is empty, just write the rune
			// This might need checking if lipgloss.Style is "zero" value effectively
			// For now, we render every cell.
			rowBuilder.WriteString(cell.Style.Render(string(cell.Rune)))
		}
		if y < c.Height-1 {
			rowBuilder.WriteRune('\n')
		}
		b.WriteString(rowBuilder.String())
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
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			c.SetChar(x+dx, y+dy, char, style)
		}
	}
}
