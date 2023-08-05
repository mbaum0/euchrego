package game

import (
	"fmt"
	"strings"
)

type Suite rune

const (
	DIAMOND Suite = '♦'
	CLUB    Suite = '♣'
	HEART   Suite = '♥'
	SPADE   Suite = '♠'
)

type Rank int

const (
	NINE Rank = iota
	TEN
	JACK
	QUEEN
	KING
	ACE
)

type Card struct {
	rank  Rank  // 9,10,jack,queen,king,ace
	suite Suite // ♦ = 0 , ♣ = 1 , ♥ = 2 , ♠ = 3
}

var LeftBauerSuite = map[Suite]Suite{
	DIAMOND: HEART,
	HEART:   DIAMOND,
	CLUB:    SPADE,
	SPADE:   CLUB,
}

func GetRank(r int) Rank {
	switch r {
	case 0:
		return NINE
	case 1:
		return TEN
	case 2:
		return JACK
	case 3:
		return QUEEN
	case 4:
		return KING
	case 5:
		return ACE
	}
	panic("Invalid input received. Rank values can be [0,5]")
}

func GetSuite(s int) Suite {
	switch s {
	case 0:
		return DIAMOND
	case 1:
		return CLUB
	case 2:
		return HEART
	case 3:
		return SPADE
	}
	panic("Invalid input received. Suite values can be [0,3]")
}

func (c *Card) GetRank() string {
	var cardRank string

	switch c.rank {
	case 2:
		cardRank = "J"
	case 3:
		cardRank = "Q"
	case 4:
		cardRank = "K"
	case 5:
		cardRank = "A"
	default:
		cardRank = fmt.Sprint(c.rank + 9)
	}
	return cardRank
}

func (c Card) Info() string {
	// translate card to human readable info
	var cardRank, cardSuite string

	switch c.rank {
	case 2:
		cardRank = "J"
	case 3:
		cardRank = "Q"
	case 4:
		cardRank = "K"
	case 5:
		cardRank = "A"
	default:
		cardRank = fmt.Sprint(c.rank + 9)
	}

	// translate grate to human readable info
	switch c.suite {
	case DIAMOND:
		cardSuite = "♦"
	case CLUB:
		cardSuite = "♣"
	case HEART:
		cardSuite = "♥"
	case SPADE:
		cardSuite = "♠"
	default:
	}

	return fmt.Sprintf("%s of %s", cardRank, cardSuite)
}

func (c *Card) GetRanking(trump Suite, lead Suite) int {
	value := 1 // start at 1 because we have some value if trump or lead

	if c.suite != trump && c.suite != lead {
		return 0
	}

	if c.suite == trump {
		value += 10

		// if right bauer
		if c.rank == JACK {
			value += 5
		}
	}

	// check for left bauer
	if c.suite == LeftBauerSuite[trump] && c.rank == JACK {
		value += 10 // 10 pts for left suite
		value += 4  // 4 pts for being left bauer
	}

	value += int(c.rank)
	return value
}

// Compare : positive if c1 > c2, 0 if c1 = c2, negative if c1 < c2
func (c *Card) Compare(c2 Card, trump Suite, lead Suite) int {
	r1 := c.GetRanking(trump, lead)
	r2 := c2.GetRanking(trump, lead)

	return r1 - r2
}

// GetWinningCard : returns the card with the highest value
func GetWinningCard(c1 Card, c2 Card, c3 Card, c4 Card, trump Suite, lead Suite) Card {
	winner := c1

	if c2.Compare(winner, trump, lead) > 0 {
		winner = c2
	}
	if c3.Compare(winner, trump, lead) > 0 {
		winner = c3
	}
	if c4.Compare(winner, trump, lead) > 0 {
		winner = c4
	}
	return winner
}

func (c *Card) GetCardArt() string {
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

	upperRank := rank
	lowerRank := rank

	if rank == "10" {
		upperRank = "\b10"
		lowerRank = "\b10"
	} else {
		upperRank = "\b" + rank + " "
		lowerRank = rank
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

func GetHandArt(cards []*Card) string {

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
		builder.WriteString("\n")
	}
	return builder.String()
}
