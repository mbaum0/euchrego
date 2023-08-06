package game

type Suite int

const (
	NONE Suite = iota
	DIAMOND
	CLUB
	HEART
	SPADE
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
