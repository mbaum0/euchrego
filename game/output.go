package game

import (
	"fmt"
	"strings"
)

func (c Card) GetCardArt() string {
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
	default:
		return ""
	}

	rank := c.GetRank()

	var upperRank string
	var lowerRank string

	if rank == TEN {
		upperRank = "\b10"
		lowerRank = "\b10"
	} else {
		upperRank = "\b" + rank.ToChar() + " "
		lowerRank = rank.ToChar()
	}

	cardArt := fmt.Sprintf(`
┌─────────┐
│  %s      │
│         │
│         │
│  %s  │
│         │
│         │
│       %s │
└─────────┘
`, upperRank, suitSymbol, lowerRank)

	return cardArt
}

func GetHandArt(cards []*Card, enumerate bool) string {

	var cardArts = make([]string, 0)
	for _, c := range cards {
		cardArts = append(cardArts, c.GetCardArt())
	}

	var builder strings.Builder
	rows := len(strings.Split(cardArts[0], "\n")) // cards have the same number of rows
	for row := 0; row < rows; row++ {
		for _, cardArt := range cardArts {
			lines := strings.Split(cardArt, "\n")
			builder.WriteString(lines[row] + " ")
		}
		if row != rows-1 {
			builder.WriteString("\n")
		} else {
			builder.WriteString("\r")
		}

	}

	if enumerate {
		for i := range cards {
			startSpaces := strings.Repeat(" ", 4)
			endSpaces := strings.Repeat(" ", 5)
			builder.WriteString(fmt.Sprintf("%s(%d)%s", startSpaces, i, endSpaces))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
