package termui_test

import (
	"testing"

	"github.com/mbaum0/euchrego/termui"
	"github.com/stretchr/testify/assert"
)

func TestNewTermUI(t *testing.T) {
	display, err := termui.NewTermUI()
	assert.Nil(t, err)
	assert.Equal(t, 80, display.Width())
	assert.Equal(t, 24, display.Height())
}

func TestNewTermUIWithSize(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	assert.Equal(t, 10, display.Width())
	assert.Equal(t, 10, display.Height())
}

func TestReset(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)

	// draw a % in the top left corner
	display.DrawRune('%', 0, 0)
	assert.Equal(t, "%", display.Grid[0][0])

	// clear the display
	display.Reset()
	emptyGrid := make([][]string, 10)
	for i := 0; i < 10; i++ {
		emptyGrid[i] = make([]string, 10)
	}
	// fill empty grid with spaces
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			emptyGrid[i][j] = " "
		}
	}
	assert.Equal(t, emptyGrid, display.Grid)
}

func TestDrawRune(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	display.DrawRune('%', 0, 0)
	assert.Equal(t, "%", display.Grid[0][0])
}

func TestDrawRuneOutOfBounds(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawRune('%', 10, 10)
	assert.NotNil(t, err)
}

func TestDrawVerticalLine(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawVerticalLine(0, 0, 10)
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		assert.Equal(t, "│", display.Grid[i][0])
	}
}

func TestDrawVerticalLineOutOfBounds(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawVerticalLine(0, 0, 11)
	assert.NotNil(t, err)
}

func TestDrawHorizontalLine(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawHorizontalLine(0, 0, 10)
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		assert.Equal(t, "─", display.Grid[0][i])
	}
}

func TestDrawHorizontalLineOutOfBounds(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawHorizontalLine(0, 0, 11)
	assert.NotNil(t, err)
}

func TestDrawRect(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawRect(0, 0, 10, 10)
	assert.Nil(t, err)
	// test corners
	assert.Equal(t, "┌", display.Grid[0][0])
	assert.Equal(t, "┐", display.Grid[0][9])
	assert.Equal(t, "└", display.Grid[9][0])
	assert.Equal(t, "┘", display.Grid[9][9])

	// test top and bottom
	for i := 1; i < 9; i++ {
		assert.Equal(t, "─", display.Grid[0][i])
		assert.Equal(t, "─", display.Grid[9][i])
	}

	// test left and right
	for i := 1; i < 9; i++ {
		assert.Equal(t, "│", display.Grid[i][0])
		assert.Equal(t, "│", display.Grid[i][9])
	}
}

func TestDrawRectOutOfBounds(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawRect(0, 0, 11, 11)
	assert.NotNil(t, err)
}

func TestDrawText(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(11, 11))
	assert.Nil(t, err)
	err = display.DrawText("Hello World", 0, 0)
	assert.Nil(t, err)
	// test char array
	expected := []string{"H", "e", "l", "l", "o", " ", "W", "o", "r", "l", "d"}
	assert.Equal(t, expected, display.Grid[0][0:11])
}

func TestDrawTextOutOfBounds(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(10, 10))
	assert.Nil(t, err)
	err = display.DrawText("Hello World", 0, 0)
	assert.NotNil(t, err)
}

func TextBoundFuncs(t *testing.T) {
	display, err := termui.NewTermUI(termui.Size(15, 15))
	assert.Nil(t, err)
	assert.Equal(t, display.Top(), 0)
	assert.Equal(t, display.Bottom(), 14)
	assert.Equal(t, display.Right(), 14)
	assert.Equal(t, display.Left(), 0)
}
