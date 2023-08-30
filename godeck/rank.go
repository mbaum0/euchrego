package godeck

type Rank int

const (
	Null Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var rankToString = map[Rank]string{
	Null:  "Null",
	Two:   "Two",
	Three: "Three",
	Four:  "Four",
	Five:  "Five",
	Six:   "Six",
	Seven: "Seven",
	Eight: "Eight",
	Nine:  "Nine",
	Ten:   "Ten",
	Jack:  "Jack",
	Queen: "Queen",
	King:  "King",
	Ace:   "Ace",
}

var rankToSymbol = map[Rank]string{
	Null:  "0",
	Two:   "2",
	Three: "3",
	Four:  "4",
	Five:  "5",
	Six:   "6",
	Seven: "7",
	Eight: "8",
	Nine:  "9",
	Ten:   "10",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

var defaultRankValues = map[Rank]int{
	Null:  0,
	Two:   2,
	Three: 3,
	Four:  4,
	Five:  5,
	Six:   6,
	Seven: 7,
	Eight: 8,
	Nine:  9,
	Ten:   10,
	Jack:  11,
	Queen: 12,
	King:  13,
	Ace:   14,
}

// String returns the string representation of a given rank.
func (r Rank) String() string {
	return rankToString[r]
}

// Symbol returns the symbol representation of a given rank.
func (r Rank) Symbol() string {
	return rankToSymbol[r]
}
