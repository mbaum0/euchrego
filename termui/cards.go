package termui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mbaum0/euchrego/godeck"
)

func getCardArt(c godeck.Card) [][]rune {
	cardRows := 9
	cardCols := 11
	cardArt := make([][]rune, cardRows)
	for i := 0; i < cardRows; i++ {
		cardArt[i] = make([]rune, cardCols)
	}

	if c == godeck.EmptyCard() {
		cardArt[0] = []rune("┌─────────┐")
		cardArt[1] = []rune("│         │")
		cardArt[2] = []rune("│         │")
		cardArt[3] = []rune("│         │")
		cardArt[4] = []rune("│         │")
		cardArt[5] = []rune("│         │")
		cardArt[6] = []rune("│         │")
		cardArt[7] = []rune("│         │")
		cardArt[8] = []rune("└─────────┘")
		return cardArt
	}

	var suitSymbol string
	switch c.Suit() {
	case godeck.Hearts:
		suitSymbol = "♥ ♥ ♥"
	case godeck.Diamonds:
		suitSymbol = "♦ ♦ ♦"
	case godeck.Clubs:
		suitSymbol = "♣ ♣ ♣"
	case godeck.Spades:
		suitSymbol = "♠ ♠ ♠"
	case godeck.None:
		suitSymbol = "   "
	}

	rank := c.Rank()
	rankChar := rank.Symbol()

	cardArt[0] = []rune("┌─────────┐")
	cardArt[1] = []rune(fmt.Sprintf("│ %s       │", rankChar))
	cardArt[2] = []rune("│         │")
	cardArt[3] = []rune("│         │")
	cardArt[4] = []rune(fmt.Sprintf("│  %s  │", suitSymbol))
	cardArt[5] = []rune("│         │")
	cardArt[6] = []rune("│         │")
	cardArt[7] = []rune(fmt.Sprintf("│       %s │", rankChar))
	cardArt[8] = []rune("└─────────┘")

	// 10 is a special case because it has two characters
	if rank == godeck.Ten {
		cardArt[1] = []rune("│ 10      │")
		cardArt[7] = []rune("│     10  │")
	}

	return cardArt
}

func (t *TermUI) DrawCard(x, y int, c godeck.Card) {
	cardArt := getCardArt(c)

	colorWay := color.New(color.FgWhite).SprintFunc()

	switch c.Suit() {
	case godeck.Hearts:
		colorWay = color.New(color.FgRed).SprintFunc()
	case godeck.Diamonds:
		colorWay = color.New(color.FgMagenta).SprintFunc()
	case godeck.Clubs:
		colorWay = color.New(color.FgYellow).SprintFunc()
	case godeck.Spades:
		colorWay = color.New(color.FgGreen).SprintFunc()
	}

	for i, row := range cardArt {
		for j, cell := range row {
			t.Grid[y+i][x+j] = colorWay(string(cell))
		}
	}
}

// used for clearing out old cards
func (t *TermUI) DrawBlankCard(x, y int) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 11; j++ {
			t.Grid[y+i][x+j] = " "
		}
	}
}

func (t *TermUI) DrawHand(x, y int, cards []godeck.Card, enumerate bool, compress bool) {
	gap := 12
	if compress {
		gap = 5
	}
	for i := 0; i < 5; i++ {
		t.DrawBlankCard(x+i*gap, y)
	}
	for i, card := range cards {
		t.DrawCard(x+i*gap, y, card)
	}

	// draw the index of the card beneath each card
	if enumerate {
		for i := range cards {
			correction := 0
			if compress {
				correction = 2
			}
			t.DrawText(fmt.Sprintf("(%d)", i), x+4+i*gap-correction, y+9)
		}
	}
}
