package termui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"unicode/utf8"

	"github.com/fatih/color"
)

type TermUI struct {
	width        int
	height       int
	Grid         [][]string
	inputEnabled bool
	reader       *bufio.Reader
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

// EnableInput enables input capture
func EnableInput() func(*TermUI) error {
	return func(t *TermUI) error {
		t.inputEnabled = true
		t.reader = bufio.NewReader(os.Stdin)
		t.height -= 2 // shrink bound box for input prompt
		return nil
	}
}

type Option func(*TermUI) error

// NewTermUI creates a new TermUI with options
func NewTermUI(options ...Option) (*TermUI, error) {
	t := &TermUI{}
	defaults := []Option{
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

func (t *TermUI) Top() int {
	return 0
}

func (t *TermUI) Bottom() int {
	return t.height - 1
}

func (t *TermUI) Left() int {
	return 0
}

func (t *TermUI) Right() int {
	return t.width - 1
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

	if t.inputEnabled {
		moveCursorUp(2)
		clearLine()
		fmt.Print(" > ")
	}
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
func (t *TermUI) DrawRune(r rune, x, y int) error {
	// error if coordinates are out of bounds
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}
	t.Grid[y][x] = string(r)
	return nil
}

// DrawRune draws a char at the specified coordinates
func (t *TermUI) DrawChar(c string, x, y int) error {
	lC := utf8.RuneCountInString(c)
	if lC > 1 {
		return fmt.Errorf("got a string. expected a single char")
	}
	if !t.isInBounds(x, y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}
	t.Grid[y][x] = c
	return nil
}

func (t *TermUI) DrawText(text string, x int, y int, options ...func(*CustomText) error) error {
	customText := &CustomText{}
	customText.text = text
	customText.x = x
	customText.y = y
	defaults := []func(*CustomText) error{Color(White), Width(len(text)), Justify(Left)}
	defaults = append(defaults, options...)
	for _, option := range defaults {
		option(customText)
	}

	if !t.isInBounds(customText.x, customText.y) {
		return fmt.Errorf("coordinates (%d, %d) are out of bounds", customText.x, customText.y)
	}

	// error if text goes out of bounds
	if customText.x+len(customText.text) > t.width {
		return fmt.Errorf("text goes out of bounds")
	}

	var c *color.Color
	switch customText.color {
	case Green:
		c = color.New(color.FgGreen)
	case Red:
		c = color.New(color.FgRed)
	case Yellow:
		c = color.New(color.FgYellow)
	case White:
		c = color.New(color.FgWhite)
	case Blue:
		c = color.New(color.FgBlue)
	default:
		c = color.New(color.FgWhite)
	}

	for i, r := range customText.text {
		t.Grid[customText.y][customText.x+i] = c.Sprint(string(r))
	}

	return nil
}

func (t *TermUI) DrawTitle(text string) error {
	t.DrawText(text, t.Left()+4, t.Top())
	return nil
}

func (t *TermUI) PollForInput() string {
	input, err := t.reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return input[:len(input)-1]
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

func clearLine() {
	fmt.Print("\x1b[K")
}

func moveCursorUp(lines int) {
	escapeSequence := fmt.Sprintf("\x1b[%dA", lines)
	fmt.Print(escapeSequence)
}
