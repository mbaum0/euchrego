package card

func NewEuchreDeck(options ...func(*Deck) error) *Deck {
	d := &Deck{}
	options = append(options, Range(Nine, Ace))
	for _, option := range options {
		option(d)
	}
	return d
}
