package game

import (
	"math/rand"
	"time"
)

type Deck struct {
	cards []*Card
}

func InitDeck() *Deck {
	deck := Deck{}
	var cards = make([]*Card, 0)
	for i := 0; i < 24; i++ {
		c := Card{rank: IntToRank(i % 6), suite: IntToSuite(i / 6)}
		cards = append(cards, &c)
	}
	deck.cards = cards
	return &deck
}

func (d *Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

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

func (d *Deck) ReturnCards(cards []*Card) {
	d.cards = append(d.cards, cards...)
}
