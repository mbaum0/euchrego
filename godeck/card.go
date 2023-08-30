package godeck

type Card struct {
	rank Rank
	suit Suit
}

func EmptyCard() Card {
	return Card{rank: Null, suit: None}
}

func NewCard(rank Rank, suit Suit) Card {
	return Card{rank: rank, suit: suit}
}

func (c Card) Rank() Rank {
	return c.rank
}

func (c Card) Suit() Suit {
	return c.suit
}

func (c Card) String() string {
	return c.rank.String() + " of " + c.suit.String()
}
