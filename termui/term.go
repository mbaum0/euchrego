package termui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type TermUI struct {
	width  int
	height int
	Grid   [][]string
}

// Size sets the size of the grid
func Size(width, height int) func(*TermUI) error {
	return func(t *TermUI) error {
		t.width = width
		t.height = height
		t.Grid = make([][]string, height)
		for i := 0; i < height; i++ {
			t.Grid[i] = make([]string, width)
		}
		return nil
	}
}

// NewTermUI creates a new TermUI with options
func NewTermUI(options ...func(*TermUI) error) (*TermUI, error) {
	t := &TermUI{}
	defaults := []func(*TermUI) error{
		Size(80, 24),
	}
	defaults = append(defaults, options...)
	for _, option := range defaults {
		option(t)
	}
	t.Reset()
	clearTerminal()
	return t, nil
}

func (t *TermUI) Width() int {
	return t.width
}

func (t *TermUI) Height() int {
	return t.height
}

// Render renders the grid to the terminal
func (t *TermUI) Render() {
	moveCursorHome()
	for _, row := range t.Grid {
		// combine the row into a string
		rowString := ""
		for _, cell := range row {
			rowString += cell
		}
		fmt.Print(rowString)

		fmt.Print("\n")
	}
	fmt.Println()
}

// Reset fills the grid with spaces
func (t *TermUI) Reset() {
	for i := 0; i < t.height; i++ {
		for j := 0; j < t.width; j++ {
			t.Grid[i][j] = " "
		}
	}
}

// DrawRune draws a rune at the specified coordinates
func (t *TermUI) DrawRune(x, y int, r rune) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}
	t.Grid[y][x] = string(r)
	return nil
}

// DrawVerticalLine draws a vertical line at the specified coordinates
func (t *TermUI) DrawVerticalLine(x, y, length int) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if line goes out of bounds
	if y+length > t.height {
		return fmt.Errorf("line goes out of bounds")
	}

	for i := 0; i < length; i++ {
		t.Grid[y+i][x] = "│"
	}
	return nil
}

// DrawHorizontalLine draws a horizontal line at the specified coordinates
func (t *TermUI) DrawHorizontalLine(x, y, length int) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if line goes out of bounds
	if x+length > t.width {
		return fmt.Errorf("line goes out of bounds")
	}

	for i := 0; i < length; i++ {
		t.Grid[y][x+i] = "─"
	}
	return nil
}

// DrawRect draws a rectangle at the specified coordinates
func (t *TermUI) DrawRect(x, y, width, height int) error {
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
	t.DrawHorizontalLine(x, y, width)
	t.DrawHorizontalLine(x, y+height-1, width)

	// draw left and right lines
	t.DrawVerticalLine(x, y, height)
	t.DrawVerticalLine(x+width-1, y, height)

	// draw corners
	t.DrawRune(x, y, '┌')
	t.DrawRune(x+width-1, y, '┐')
	t.DrawRune(x, y+height-1, '└')
	t.DrawRune(x+width-1, y+height-1, '┘')

	return nil
}

// DrawText draws text at the specified coordinates
func (t *TermUI) DrawText(x, y int, text string) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if text goes out of bounds
	if x+len(text) > t.width {
		return fmt.Errorf("text goes out of bounds")
	}

	for i, r := range text {
		t.Grid[y][x+i] = string(r)
	}

	return nil
}

// DrawTextCentered draws text centered at the specified coordinates
func (t *TermUI) DrawTextCentered(x, y int, text string) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}

	// error if text goes out of bounds
	if x+len(text)/2 > t.width || x-len(text)/2 < 0 {
		return fmt.Errorf("text goes out of bounds")
	}

	for i, r := range text {
		t.Grid[y][x+i-len(text)/2] = string(r)
	}

	return nil
}

func (t *TermUI) isInBounds(x, y int) bool {
	return x >= 0 && x < t.width && y >= 0 && y < t.height
}

func clearTerminal() {
	// Clear command based on the operating system
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	// Set the correct output device for the command
	cmd.Stdout = os.Stdout

	// Run the command to clear the screen
	cmd.Run()
}

func moveCursorHome() {
	fmt.Print("\033[H")
}
