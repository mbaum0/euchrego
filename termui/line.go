package termui

import "fmt"

type LineStyle int

const (
	Solid LineStyle = iota
	Dashed
	Thick
	Double
)

type LineOptions func(*CustomLine) error
type CustomLine struct {
	style  LineStyle
	vRune  string
	hRune  string
	tlRune string
	trRune string
	blRune string
	brRune string
}

func Style(s LineStyle) func(*CustomLine) error {
	return func(l *CustomLine) error {
		l.style = s
		switch s {
		case Solid:
			return nil
		case Dashed:
			l.vRune = "┄"
			l.hRune = "┆"
		case Thick:
			l.vRune = "┃"
			l.hRune = "━"
			l.tlRune = "┏"
			l.trRune = "┓"
			l.blRune = "┗"
			l.brRune = "┛"
		case Double:
			l.vRune = "║"
			l.hRune = "═"
			l.tlRune = "╔"
			l.trRune = "╗"
			l.blRune = "╚"
			l.brRune = "╝"
		}
		return nil
	}
}

func NewCustomLine(options ...func(*CustomLine) error) *CustomLine {
	cl := &CustomLine{Solid, "│", "─", "┌", "┐", "└", "┘"}
	for _, option := range options {
		option(cl)
	}
	return cl
}

// DrawVerticalLine draws a vertical line at the specified coordinates
func (t *TermUI) DrawVerticalLine(x, y, length int, options ...func(*CustomLine) error) error {
	cl := NewCustomLine(options...)
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if line goes out of bounds
	if y+length > t.height {
		return fmt.Errorf("line goes out of bounds")
	}

	for i := 0; i < length; i++ {
		t.Grid[y+i][x] = cl.vRune
	}
	return nil
}

// DrawHorizontalLine draws a horizontal line at the specified coordinates
func (t *TermUI) DrawHorizontalLine(x, y, length int, options ...func(*CustomLine) error) error {
	cl := NewCustomLine(options...)
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if line goes out of bounds
	if x+length > t.width {
		return fmt.Errorf("line goes out of bounds")
	}

	for i := 0; i < length; i++ {
		t.Grid[y][x+i] = cl.hRune
	}
	return nil
}

// DrawRect draws a rectangle at the specified coordinates
func (t *TermUI) DrawRect(x, y, width, height int, options ...func(*CustomLine) error) error {
	cl := NewCustomLine(options...)
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// errors if rectangle goes out of bounds
	if x+width > t.width {
		return fmt.Errorf("rect goes out of width bounds")
	}

	if y+height > t.height {
		return fmt.Errorf("rect goes out of height bounds")
	}

	// draw top and bottom lines
	t.DrawHorizontalLine(x, y, width, options...)
	t.DrawHorizontalLine(x, y+height-1, width, options...)

	// draw left and right lines
	t.DrawVerticalLine(x, y, height, options...)
	t.DrawVerticalLine(x+width-1, y, height, options...)

	// draw corners
	t.DrawChar(cl.tlRune, x, y)
	t.DrawChar(cl.trRune, x+width-1, y)
	t.DrawChar(cl.blRune, x, y+height-1)
	t.DrawChar(cl.brRune, x+width-1, y+height-1)

	return nil
}
