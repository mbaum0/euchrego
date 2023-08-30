package godeck

import (
	"errors"
	"fmt"
	"math/rand"
)

type Deck struct {
	cards       []Card
	shuffleSeed int64
	rankValues  map[Rank]int
	deckRange   [2]Rank
	preShuffled bool
}

func (d *Deck) populate() {
	for rank := d.deckRange[0]; rank <= d.deckRange[1]; rank++ {
		for suit := Hearts; suit <= Spades; suit++ {
			d.cards = append(d.cards, NewCard(rank, suit))
		}
	}
}

// Range sets the range of the deck. The deck will contain all cards from start to end, inclusive.
func Range(start, end Rank) func(*Deck) error {
	return func(d *Deck) error {
		if start > end {
			return errors.New("start must be less than or equal to end")
		}
		d.deckRange = [2]Rank{start, end}
		return nil
	}
}

// ShuffleSeed sets the shuffle seed to the provided value
func ShuffleSeed(seed int64) func(*Deck) error {
	return func(d *Deck) error {
		d.shuffleSeed = seed
		return nil
	}
}

// RandomShuffleSeed sets the shuffle seed to a random number
func RandomShuffleSeed() func(*Deck) error {
	return func(d *Deck) error {
		d.shuffleSeed = rand.Int63()
		return nil
	}
}

// RankValues sets the rank values to the provided map
func RankValues(rankValues map[Rank]int) func(*Deck) error {
	return func(d *Deck) error {
		// copy rankValues to prevent mutation of the default map
		d.rankValues = make(map[Rank]int)
		for rank, value := range rankValues {
			d.rankValues[rank] = value
		}
		return nil
	}
}

// AceHigh sets the value of an ace to 14
func AceHigh() func(*Deck) error {
	return func(d *Deck) error {
		d.rankValues[Ace] = 14
		return nil
	}
}

// AceLow sets the value of an ace to 1
func AceLow() func(*Deck) error {
	return func(d *Deck) error {
		d.rankValues[Ace] = 1
		return nil
	}
}

// PreShuffled shuffles the deck before returning it
func PreShuffled() func(*Deck) error {
	return func(d *Deck) error {
		d.preShuffled = true
		return nil
	}
}

func NewDeck(options ...func(*Deck) error) (*Deck, error) {
	deck := &Deck{}
	// always use the defaults, as they can be overriden by the options
	defaults := []func(*Deck) error{Range(Two, Ace), RankValues(defaultRankValues), AceHigh(), RandomShuffleSeed()}
	defaults = append(defaults, options...)
	for _, option := range defaults {
		err := option(deck)
		if err != nil {
			return nil, err
		}
	}

	deck.populate()

	if deck.preShuffled {
		deck.Shuffle()
	}

	return deck, nil
}

func (d *Deck) Shuffle() {
	rng := rand.New(rand.NewSource(d.shuffleSeed))
	for i := range d.cards {
		j := rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Length() int {
	return len(d.cards)
}

func (d *Deck) DrawCard() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, errors.New("no cards left in deck")
	}
	index := len(d.cards) - 1
	card := d.cards[index]
	d.cards = d.cards[:index]
	return card, nil
}

func (d *Deck) SetRankValues(rankValues map[Rank]int) {
	// copy rankValues to prevent mutation of the default map
	for rank, value := range rankValues {
		d.rankValues[rank] = value
	}
}

func (d *Deck) SetRankValue(rank Rank, value int) {
	d.rankValues[rank] = value
}

func (d *Deck) ResetRankValues() {
	// copy default rank values to prevent mutation of the default map
	for rank, value := range defaultRankValues {
		d.rankValues[rank] = value
	}
}

func (d *Deck) RankValue(rank Rank) int {
	return d.rankValues[rank]
}

func (d *Deck) DrawCards(numCards int) ([]Card, error) {
	cards := make([]Card, 0)

	for i := 0; i < numCards; i++ {
		card, err := d.DrawCard()
		if err != nil {
			// put cards back in deck
			d.cards = append(d.cards, cards...)
			errMsg := fmt.Sprintf("not enough cards in deck, tried to draw %d cards when only %d are in the deck", numCards, d.Length())
			return nil, errors.New(errMsg)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (d *Deck) ReturnCards(cards []Card) {
	d.cards = append(d.cards, cards...) // add cards back to the deck
}

func (d *Deck) ReturnCard(card Card) {
	d.cards = append(d.cards, card)
}
