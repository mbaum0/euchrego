package main

import "fmt"

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

type Deck []Card

func (c Card) Info() (string, string) {
	// translate card to human readable info
	var cardName, cardSuit string

	switch c.rank {
	case 11:
		cardName = "J"
	case 12:
		cardName = "Q"
	case 13:
		cardName = "K"
	case 14:
		cardName = "A"
	default:
		cardName = fmt.Sprint(c.rank)
	}

	// translate grate to human readable info
	switch c.suite {
	case 0:
		cardSuit = "♦"
	case 1:
		cardSuit = "♣"
	case 2:
		cardSuit = "♥"
	case 3:
		cardSuit = "♠"
	default:
	}

	return cardName, cardSuit
}

func GetCardRanking(c Card, trump Suite, lead Suite) int {
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

// IsCardEqual : positive if c1 > c2, 0 if c1 = c2, negative if c1 < c2
func IsCardEqual(c1 Card, c2 Card, trump Suite, lead Suite) int {
	r1 := GetCardRanking(c1, trump, lead)
	r2 := GetCardRanking(c2, trump, lead)

	return r1 - r2
}
