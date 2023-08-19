package game

import (
	"math/rand"
)

type Deck struct {
	cards       []*Card
	ShuffleSeed int64
}

func InitDeck(shuffleSeed int64) Deck {
	deck := Deck{}
	var cards = make([]*Card, 0)
	for i := 0; i < 24; i++ {
		c := Card{rank: IntToRank(i % 6), suite: IntToSuite(i / 6)}
		cards = append(cards, &c)
	}
	deck.cards = cards
	deck.ShuffleSeed = shuffleSeed
	return deck
}

func (d *Deck) Shuffle() {
	rng := rand.New(rand.NewSource(d.ShuffleSeed))

	n := len(d.cards)
	for i := n - 1; i >= 1; i-- {
		j := rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Length() int {
	return len(d.cards)
}

func (d *Deck) pop() *Card {
	index := len(d.cards) - 1
	card := d.cards[index]
	d.cards = d.cards[:index]
	return card
}

func (d *Deck) DrawCards(numCards int) []*Card {
	cards := make([]*Card, 0)

	for i := 0; i < numCards; i++ {
		cards = append(cards, d.pop())
	}
	return cards
}

func (d *Deck) ReturnCards(cards *[]*Card) {
	d.cards = append(d.cards, *cards...) // add cards back to the deck
	*cards = make([]*Card, 0)            // remove cards from the input array
}

func (d *Deck) ReturnCard(card *Card) {
	d.cards = append(d.cards, card)
}
