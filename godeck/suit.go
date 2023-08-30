package godeck

type Suit int

const (
	None Suit = iota
	Hearts
	Diamonds
	Clubs
	Spades
)

var suitToString = map[Suit]string{
	None:     "None",
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
	None:     "N",
}

func (s Suit) String() string {
	return suitToString[s]
}

func (s Suit) Symbol() string {
	return suitToSymbol[s]
}
