package card

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

var suitToString = map[Suit]string{
	Hearts:   "Hearts",
	Diamonds: "Diamonds",
	Clubs:    "Clubs",
	Spades:   "Spades",
}

var suitToSymbol = map[Suit]string{
	Hearts:   "♥",
	Diamonds: "♦",
	Clubs:    "♣",
	Spades:   "♠",
}

func (s Suit) String() string {
	return suitToString[s]
}

func (s Suit) Symbol() string {
	return suitToSymbol[s]
}
