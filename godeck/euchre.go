package godeck

var LeftBauerSuite = map[Suit]Suit{
	Diamonds: Hearts,
	Hearts:   Diamonds,
	Clubs:    Spades,
	Spades:   Clubs,
}

type EuchreDeck struct {
	Deck
}

func NewEuchreDeck(options ...func(*Deck) error) *EuchreDeck {
	d := &Deck{}
	options = append(options, Range(Nine, Ace))
	for _, option := range options {
		option(d)
	}

	ed := EuchreDeck{Deck: *d}
	return &ed
}

func (d *EuchreDeck) isCardLeftBauer(c Card, trump Suit) bool {
	return c.suit == LeftBauerSuite[trump] && c.rank == Jack
}

func (d *EuchreDeck) getCardValue(c Card, trump Suit, lead Suit) int {
	value := 1 // start at 1 because we have some value if trump or lead

	if c.suit != trump && c.suit != lead && !d.isCardLeftBauer(c, trump) {
		return 0
	}

	if c.suit == trump {
		value += 10

		// if right bauer
		if c.rank == Jack {
			value += 5
		}
	}

	// check for left bauer
	if d.isCardLeftBauer(c, trump) {
		value += 10 // 10 pts for left suite
		value += 4  // 4 pts for being left bauer
	}

	value += int(c.rank)
	return value
}

func (d *EuchreDeck) CompareCards(c1 Card, c2 Card, trump Suit, lead Suit) int {
	r1 := d.getCardValue(c1, trump, lead)
	r2 := d.getCardValue(c2, trump, lead)

	return r1 - r2
}

func (d *EuchreDeck) GetWinningCard(c1 Card, c2 Card, c3 Card, c4 Card, trump Suit, lead Suit) Card {
	winner := c1

	if d.CompareCards(c2, winner, trump, lead) > 0 {
		winner = c2
	}
	if d.CompareCards(c3, winner, trump, lead) > 0 {
		winner = c3
	}
	if d.CompareCards(c4, winner, trump, lead) > 0 {
		winner = c4
	}
	return winner
}

func (d *EuchreDeck) GetPlayableCards(hand []Card, trump Suit, lead Card) []Card {
	if lead == EmptyCard() {
		return hand
	}

	var playableCards = make([]Card, 0)

	// was the left bauer led?
	leftBauerWasLed := d.isCardLeftBauer(lead, trump)

	if leftBauerWasLed {
		// the lead suite is actually the trump suite in this case
		lead.suit = trump
	}

	// was trump led?
	trumpWasLed := lead.suit == trump || leftBauerWasLed

	// if trump was led, we must play trump if we have it
	hasTrumpCards := false
	if trumpWasLed {
		for _, c := range hand {
			if c.suit == trump || d.isCardLeftBauer(c, trump) {
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
		// left bauer is a different suite than it shows
		if d.isCardLeftBauer(c, trump) {
			continue
		}

		if c.suit == lead.suit {
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

func (d *EuchreDeck) IsCardPlayable(card Card, hand []Card, trump Suit, lead Card) bool {
	playableCards := d.GetPlayableCards(hand, trump, lead)

	for _, c := range playableCards {
		if c == card {
			return true
		}
	}
	return false
}
