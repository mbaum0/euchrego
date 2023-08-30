package godeck_test

import (
	"testing"

	"github.com/mbaum0/euchrego/godeck"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	// test both trump
	d := godeck.NewEuchreDeck()
	c1 := godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c2 := godeck.NewCard(godeck.Nine, godeck.Diamonds)
	res := d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer
	c1 = godeck.NewCard(godeck.Jack, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Ace, godeck.Diamonds)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer, c2 is left
	c1 = godeck.NewCard(godeck.Jack, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Jack, godeck.Hearts)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 isn't a bauer and c2 is the left bauer
	c1 = godeck.NewCard(godeck.Nine, godeck.Hearts)
	c2 = godeck.NewCard(godeck.Jack, godeck.Diamonds)
	res = d.CompareCards(c2, c1, godeck.Hearts, godeck.Hearts)
	assert.Positive(t, res, "expected c2 to be greater than c1")

	// test c1 is trump, c2 is not
	c1 = godeck.NewCard(godeck.Nine, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Nine, godeck.Hearts)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c2 is trump, c1 is not
	c1 = godeck.NewCard(godeck.Nine, godeck.Hearts)
	c2 = godeck.NewCard(godeck.Nine, godeck.Diamonds)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test neither are trump but both lead
	c1 = godeck.NewCard(godeck.Ten, godeck.Hearts)
	c2 = godeck.NewCard(godeck.Nine, godeck.Hearts)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c1 lead, c2 off-suite
	c1 = godeck.NewCard(godeck.Ten, godeck.Hearts)
	c2 = godeck.NewCard(godeck.Nine, godeck.Spades)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c2 lead, c1 off-suite
	c1 = godeck.NewCard(godeck.Ten, godeck.Spades)
	c2 = godeck.NewCard(godeck.Nine, godeck.Hearts)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test both are off-suite
	c1 = godeck.NewCard(godeck.Ten, godeck.Spades)
	c2 = godeck.NewCard(godeck.Nine, godeck.Clubs)
	res = d.CompareCards(c1, c2, godeck.Diamonds, godeck.Hearts)
	assert.Zero(t, res, "expected c2 to be equal to c1")
}

func TestGetPlayableCards(t *testing.T) {
	d := godeck.NewEuchreDeck()
	// test there was no lead card
	var cards = make([]godeck.Card, 0)
	c1 := godeck.NewCard(godeck.Nine, godeck.Clubs)
	c2 := godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c3 := godeck.NewCard(godeck.Jack, godeck.Spades)
	c4 := godeck.NewCard(godeck.King, godeck.Clubs)
	c5 := godeck.NewCard(godeck.Ace, godeck.Clubs)
	cards = append(cards, c1, c2, c3, c4, c5)

	playableCards := d.GetPlayableCards(cards, godeck.Diamonds, godeck.EmptyCard())
	assert.Equal(t, 5, len(playableCards), "Expected 5 cards to be returned")
	cards = cards[:0]

	// test player has no lead cards and one trump card

	c1 = godeck.NewCard(godeck.Nine, godeck.Clubs)
	c2 = godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c3 = godeck.NewCard(godeck.Jack, godeck.Spades)
	c4 = godeck.NewCard(godeck.King, godeck.Clubs)
	c5 = godeck.NewCard(godeck.Ace, godeck.Clubs)
	cards = append(cards, c1, c2, c3, c4, c5)

	leadCard := godeck.NewCard(godeck.Nine, godeck.Hearts)
	playableCards = d.GetPlayableCards(cards, godeck.Diamonds, leadCard)
	assert.Equal(t, 5, len(playableCards), "Expected 5 card to be returned")
	cards = cards[:0]

	// test player has at least one lead card and no trump
	leadCard = godeck.NewCard(godeck.Ace, godeck.Diamonds)

	c1 = godeck.NewCard(godeck.Nine, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c3 = godeck.NewCard(godeck.Jack, godeck.Spades)
	c4 = godeck.NewCard(godeck.King, godeck.Spades)
	c5 = godeck.NewCard(godeck.Ace, godeck.Clubs)

	cards = append(cards, c1, c2, c3, c4, c5)

	playableCards = d.GetPlayableCards(cards, godeck.Hearts, leadCard)
	assert.Equal(t, len(playableCards), 2, "Expected 2 cards to be returned")
	assert.Contains(t, playableCards, c1, "Expected nine of diamonds to be returned")
	assert.Contains(t, playableCards, c2, "Expected ten of diamonds to be returned")
	cards = cards[:0]

	// test player has no trump or lead cards
	leadCard = godeck.NewCard(godeck.Ace, godeck.Hearts)
	c1 = godeck.NewCard(godeck.Nine, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c3 = godeck.NewCard(godeck.Jack, godeck.Spades)
	c4 = godeck.NewCard(godeck.King, godeck.Spades)
	c5 = godeck.NewCard(godeck.Ace, godeck.Clubs)
	cards = append(cards, c1, c2, c3, c4, c5)

	playableCards = d.GetPlayableCards(cards, godeck.Hearts, leadCard)
	assert.Equal(t, len(playableCards), 5, "Expected 5 cards to be returned")
	cards = cards[:0]

	// player has only the left bauer, no lead cards, and trump was led
	leadCard = godeck.NewCard(godeck.Ace, godeck.Clubs)
	c1 = godeck.NewCard(godeck.Nine, godeck.Diamonds)
	c2 = godeck.NewCard(godeck.Ten, godeck.Diamonds)
	c3 = godeck.NewCard(godeck.Jack, godeck.Spades)
	c4 = godeck.NewCard(godeck.King, godeck.Diamonds)
	c5 = godeck.NewCard(godeck.Ace, godeck.Diamonds)
	cards = append(cards, c1, c2, c3, c4, c5)

	playableCards = d.GetPlayableCards(cards, godeck.Clubs, leadCard)
	assert.Equal(t, len(playableCards), 1, "Expected 1 cards to be returned")
	assert.Contains(t, playableCards, c3, "Expected jack of spades to be returned")
}
