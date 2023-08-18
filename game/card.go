package game

import "fmt"

type Suite int

const (
	NONE Suite = iota
	DIAMOND
	CLUB
	HEART
	SPADE
)

func (s Suite) ToString() string {
	var suite string
	switch s {
	case 0:
		suite = "Pass"
	case 1:
		suite = "Diamonds"
	case 2:
		suite = "Clubs"
	case 3:
		suite = "Hearts"
	case 4:
		suite = "Spades"
	default:
		suite = ""
	}
	return suite
}

func SuiteFromChar(s string) Suite {
	switch s {
	case "h":
		return HEART

	case "d":
		return DIAMOND
	case "c":
		return CLUB
	case "s":
		return SPADE
	default:
		return NONE
	}
}

func (r Rank) ToString() string {
	var rank string

	switch r {
	case 0:
		rank = "9"
	case 1:
		rank = "10"
	case 2:
		rank = "Jack"
	case 3:
		rank = "Queen"
	case 4:
		rank = "King"
	case 5:
		rank = "Ace"
	default:
		rank = ""
	}
	return rank
}

func (r Rank) ToChar() string {
	var rank string

	switch r {
	case 0:
		rank = "9"
	case 1:
		rank = "10"
	case 2:
		rank = "J"
	case 3:
		rank = "Q"
	case 4:
		rank = "K"
	case 5:
		rank = "A"
	default:
		rank = ""
	}
	return rank
}

type Rank int

const (
	NINE Rank = iota
	TEN
	JACK
	QUEEN
	KING
	ACE
)

var LeftBauerSuite = map[Suite]Suite{
	DIAMOND: HEART,
	HEART:   DIAMOND,
	CLUB:    SPADE,
	SPADE:   CLUB,
}

type Card struct {
	rank  Rank
	suite Suite
}

func (c *Card) GetRank() Rank {
	return c.rank
}

func (c *Card) GetSuite() Suite {
	return c.suite
}

func (c *Card) ToString() string {

	return fmt.Sprintf("%s of %s", c.rank.ToString(), c.suite.ToString())
}

func IntToSuite(s int) Suite {
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
	panic("invalid int provided")
}

func IntToRank(r int) Rank {
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
	panic("invalid int provided")
}

func (c *Card) GetPlayValue(trump Suite, lead Suite) int {
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

// compare : positive if c1 > c2, 0 if c1 = c2, negative if c1 < c2
func (c *Card) compare(c2 Card, trump Suite, lead Suite) int {
	r1 := c.GetPlayValue(trump, lead)
	r2 := c2.GetPlayValue(trump, lead)

	return r1 - r2
}

// GetWinningCard : returns the card with the highest value
func GetWinningCard(c1 Card, c2 Card, c3 Card, c4 Card, trump Suite, lead Suite) Card {
	winner := c1

	if c2.compare(winner, trump, lead) > 0 {
		winner = c2
	}
	if c3.compare(winner, trump, lead) > 0 {
		winner = c3
	}
	if c4.compare(winner, trump, lead) > 0 {
		winner = c4
	}
	return winner
}

func GetPlayableCards(hand []*Card, trump Suite, lead *Card) []*Card {
	if lead == nil {
		return hand
	}

	var playableCards = make([]*Card, 0)

	// was the left bauer led?
	leftBauerWasLed := lead.suite == LeftBauerSuite[trump] && lead.rank == JACK

	if leftBauerWasLed {
		// the lead suite is actually the trump suite in this case
		lead.suite = trump
	}

	// was trump led?
	trumpWasLed := lead.suite == trump || leftBauerWasLed

	// if trump was led, we must play trump if we have it
	hasTrumpCards := false
	if trumpWasLed {
		for _, c := range hand {
			if c.suite == trump || (c.suite == LeftBauerSuite[trump] && c.rank == JACK) {
				playableCards = append(playableCards, c)
				hasTrumpCards = true
			}
		}
	}

	if hasTrumpCards {
		return playableCards
	}

	// if trump was not led, we must play lead if we have it
	hasLeadCards := false
	// check if any cards match what was lead
	for _, c := range hand {
		if c.suite == lead.suite {
			playableCards = append(playableCards, c)
			hasLeadCards = true
		}
	}

	// if we don't have trump or lead cards, all cards are valid
	if !hasLeadCards && !hasTrumpCards {
		playableCards = append(playableCards, hand...)
	}

	return playableCards
}

func IsCardPlayable(card *Card, hand []*Card, trump Suite, lead *Card) bool {
	playableCards := GetPlayableCards(hand, trump, lead)

	for _, c := range playableCards {
		if c == card {
			return true
		}
	}
	return false
}
