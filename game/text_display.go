package game

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type TextDisplay struct {
	width  int
	height int
	Grid   [][]rune
}

func NewTextDisplay(width, height int) *TextDisplay {
	t := TextDisplay{}
	t.width = width
	t.height = height
	t.Grid = make([][]rune, height)
	for i := 0; i < height; i++ {
		t.Grid[i] = make([]rune, width)
	}
	return &t
}

func (t *TextDisplay) Render() {
	clearTerminal()
	for _, row := range t.Grid {
		for _, cell := range row {
			fmt.Print(string(cell))
		}
		fmt.Print("\n")
	}
}

func (t *TextDisplay) ClearDisplay() {
	for i := 0; i < t.height; i++ {
		for j := 0; j < t.width; j++ {
			t.Grid[i][j] = ' '
		}
	}
}

func (t *TextDisplay) DrawVerticalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.Grid[y+i][x] = '|'
	}
}

func (t *TextDisplay) DrawHorizontalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.Grid[y][x+i] = '-'
	}
}

func (t *TextDisplay) DrawCard(x, y int, card Card) {
	cardArt := getCardArt(card)
	for i, row := range cardArt {
		for j, cell := range row {
			t.Grid[y+i][x+j] = cell
		}
	}
}

func getCardArt(c Card) [][]rune {
	cardRows := 9
	cardCols := 11
	cardArt := make([][]rune, cardRows)
	for i := 0; i < cardRows; i++ {
		cardArt[i] = make([]rune, cardCols)
	}

	var suitSymbol string
	switch c.suite {
	case HEART:
		suitSymbol = "♥ ♥ ♥"
	case DIAMOND:
		suitSymbol = "♦ ♦ ♦"
	case CLUB:
		suitSymbol = "♣ ♣ ♣"
	case SPADE:
		suitSymbol = "♠ ♠ ♠"
	case NONE:
		suitSymbol = "   "
	}

	rank := c.GetRank()
	rankChar := rank.ToChar()

	cardArt[0] = []rune("┌─────────┐")
	cardArt[1] = []rune(fmt.Sprintf("│  %s      │", rankChar))
	cardArt[2] = []rune("│         │")
	cardArt[3] = []rune("│         │")
	cardArt[4] = []rune(fmt.Sprintf("│  %s  │", suitSymbol))
	cardArt[5] = []rune("│         │")
	cardArt[6] = []rune("│         │")
	cardArt[7] = []rune(fmt.Sprintf("│       %s │", rankChar))
	cardArt[8] = []rune("└─────────┘")

	// 10 is a special case because it has two characters
	if rank == TEN {
		cardArt[1] = []rune("│  10     │")
		cardArt[7] = []rune("│     10  │")
	}

	return cardArt
}

func (t *TextDisplay) DrawText(x, y int, text string) {
	for i, c := range text {
		t.Grid[y][x+i] = c
	}
}

func (t *TextDisplay) DrawPlayerHand(x, y int, player Player) {
	cards := player.hand
	for i, card := range cards {
		t.DrawCard(x+i*12, y, *card)
	}
}

func (t *TextDisplay) DrawPlayerHands(game *Game) {
	player1 := game.Players[0]
	player2 := game.Players[1]
	player3 := game.Players[2]
	player4 := game.Players[3]

	t.DrawText(1, 2, player1.name)
	t.DrawPlayerHand(0, 3, player1)

	t.DrawText(1, 14, player2.name)
	t.DrawPlayerHand(0, 15, player2)

	t.DrawText(1, 26, player3.name)
	t.DrawPlayerHand(0, 27, player3)

	t.DrawText(1, 38, player4.name)
	t.DrawPlayerHand(0, 39, player4)
}

func (t *TextDisplay) DrawDealerArrow(game *Game) {
	dealerIndex := game.DealerIndex

	y := 5 + 12*dealerIndex
	x := 60
	t.DrawText(x, y, "<-- Dealer")
}

func (t *TextDisplay) DrawTurnArrow(game *Game) {
	playerIndex := game.PlayerIndex

	y := 6 + 12*playerIndex
	x := 60
	t.DrawText(x, y, "<-- Turn")
}

func (t *TextDisplay) DrawBoard(game *Game) {
	t.ClearDisplay()
	t.DrawPlayerHands(game)
	t.DrawDealerArrow(game)
	t.DrawTurnArrow(game)
	t.Render()
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
