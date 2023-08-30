package card_test

import (
	"testing"

	"github.com/mbaum0/euchrego/card"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	// test both trump
	d := card.NewEuchreDeck()
	c1 := card.NewCard(card.Ten, card.Diamonds)
	c2 := card.NewCard(card.Nine, card.Diamonds)
	res := d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer
	c1 = card.NewCard(card.Jack, card.Diamonds)
	c2 = card.NewCard(card.Ace, card.Diamonds)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer, c2 is left
	c1 = card.NewCard(card.Jack, card.Diamonds)
	c2 = card.NewCard(card.Jack, card.Hearts)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 isn't a bauer and c2 is the left bauer
	c1 = card.NewCard(card.Nine, card.Hearts)
	c2 = card.NewCard(card.Jack, card.Diamonds)
	res = d.CompareCards(c2, c1, card.Hearts, card.Hearts)
	assert.Positive(t, res, "expected c2 to be greater than c1")

	// test c1 is trump, c2 is not
	c1 = card.NewCard(card.Nine, card.Diamonds)
	c2 = card.NewCard(card.Nine, card.Hearts)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c2 is trump, c1 is not
	c1 = card.NewCard(card.Nine, card.Hearts)
	c2 = card.NewCard(card.Nine, card.Diamonds)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test neither are trump but both lead
	c1 = card.NewCard(card.Ten, card.Hearts)
	c2 = card.NewCard(card.Nine, card.Hearts)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c1 lead, c2 off-suite
	c1 = card.NewCard(card.Ten, card.Hearts)
	c2 = card.NewCard(card.Nine, card.Spades)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c2 lead, c1 off-suite
	c1 = card.NewCard(card.Ten, card.Spades)
	c2 = card.NewCard(card.Nine, card.Hearts)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test both are off-suite
	c1 = card.NewCard(card.Ten, card.Spades)
	c2 = card.NewCard(card.Nine, card.Clubs)
	res = d.CompareCards(c1, c2, card.Diamonds, card.Hearts)
	assert.Zero(t, res, "expected c2 to be equal to c1")
}

func TestGetPlayableCards(t *testing.T) {
	d := card.NewEuchreDeck()
	// test there was no lead card
	var cards = make([]*card.Card, 0)
	c1 := card.NewCard(card.Nine, card.Clubs)
	c2 := card.NewCard(card.Ten, card.Diamonds)
	c3 := card.NewCard(card.Jack, card.Spades)
	c4 := card.NewCard(card.King, card.Clubs)
	c5 := card.NewCard(card.Ace, card.Clubs)
	cards = append(cards, &c1, &c2, &c3, &c4, &c5)

	playableCards := d.GetPlayableCards(cards, card.Diamonds, nil)
	assert.Equal(t, 5, len(playableCards), "Expected 5 cards to be returned")
	cards = cards[:0]

	// test player has no lead cards and one trump card

	c1 = card.NewCard(card.Nine, card.Clubs)
	c2 = card.NewCard(card.Ten, card.Diamonds)
	c3 = card.NewCard(card.Jack, card.Spades)
	c4 = card.NewCard(card.King, card.Clubs)
	c5 = card.NewCard(card.Ace, card.Clubs)
	cards = append(cards, &c1, &c2, &c3, &c4, &c5)

	leadCard := card.NewCard(card.Nine, card.Hearts)
	playableCards = d.GetPlayableCards(cards, card.Diamonds, &leadCard)
	assert.Equal(t, 5, len(playableCards), "Expected 5 card to be returned")
	cards = cards[:0]

	// test player has at least one lead card and no trump
	leadCard = card.NewCard(card.Ace, card.Diamonds)

	c1 = card.NewCard(card.Nine, card.Diamonds)
	c2 = card.NewCard(card.Ten, card.Diamonds)
	c3 = card.NewCard(card.Jack, card.Spades)
	c4 = card.NewCard(card.King, card.Spades)
	c5 = card.NewCard(card.Ace, card.Clubs)

	cards = append(cards, &c1, &c2, &c3, &c4, &c5)

	playableCards = d.GetPlayableCards(cards, card.Hearts, &leadCard)
	assert.Equal(t, len(playableCards), 2, "Expected 2 cards to be returned")
	assert.Contains(t, playableCards, &c1, "Expected nine of diamonds to be returned")
	assert.Contains(t, playableCards, &c2, "Expected ten of diamonds to be returned")
	cards = cards[:0]

	// test player has no trump or lead cards
	leadCard = card.NewCard(card.Ace, card.Hearts)
	c1 = card.NewCard(card.Nine, card.Diamonds)
	c2 = card.NewCard(card.Ten, card.Diamonds)
	c3 = card.NewCard(card.Jack, card.Spades)
	c4 = card.NewCard(card.King, card.Spades)
	c5 = card.NewCard(card.Ace, card.Clubs)
	cards = append(cards, &c1, &c2, &c3, &c4, &c5)

	playableCards = d.GetPlayableCards(cards, card.Hearts, &leadCard)
	assert.Equal(t, len(playableCards), 5, "Expected 5 cards to be returned")
	cards = cards[:0]

	// player has only the left bauer, no lead cards, and trump was led
	leadCard = card.NewCard(card.Ace, card.Clubs)
	c1 = card.NewCard(card.Nine, card.Diamonds)
	c2 = card.NewCard(card.Ten, card.Diamonds)
	c3 = card.NewCard(card.Jack, card.Spades)
	c4 = card.NewCard(card.King, card.Diamonds)
	c5 = card.NewCard(card.Ace, card.Diamonds)
	cards = append(cards, &c1, &c2, &c3, &c4, &c5)

	playableCards = d.GetPlayableCards(cards, card.Clubs, &leadCard)
	assert.Equal(t, len(playableCards), 1, "Expected 1 cards to be returned")
	assert.Contains(t, playableCards, &c3, "Expected jack of spades to be returned")
}
