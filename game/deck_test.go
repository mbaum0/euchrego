package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDeck(t *testing.T) {
	deck := InitDeck(0)

	assert.Equal(t, len(deck.cards), 24, "Deck should contain 24 cards")

	// make sure we have 6 of each suite
	hearts := 0
	clubs := 0
	diamonds := 0
	spades := 0
	for _, c := range deck.cards {
		switch c.suite {
		case HEART:
			hearts += 1
		case CLUB:
			clubs += 1
		case DIAMOND:
			diamonds += 1
		case SPADE:
			spades += 1
		}
	}

	assert.Equal(t, hearts, 6, "Expected 6 hearts")
	assert.Equal(t, clubs, 6, "Expected 6 clubs")
	assert.Equal(t, diamonds, 6, "Expected 6 diamonds")
	assert.Equal(t, spades, 6, "Expected 6 spades")
}

func TestDrawCards(t *testing.T) {
	deck := InitDeck(0)
	assert.Len(t, deck.cards, 24, "Deck should contain 24 cards")
	cards := deck.DrawCards(2)
	assert.Len(t, deck.cards, 22, "Deck should contain 22 cards")
	cards = append(cards, deck.DrawCards(3)...)
	assert.Len(t, deck.cards, 19, "Deck should contain 19 cards")
	cards = append(cards, deck.DrawCards(2)...)
	assert.Len(t, deck.cards, 17, "Deck should contain 17 cards")
	cards = append(cards, deck.DrawCards(3)...)
	assert.Len(t, deck.cards, 14, "Deck should contain 14 cards")
	cards = append(cards, deck.DrawCards(2)...)
	assert.Len(t, deck.cards, 12, "Deck should contain 12 cards")
	cards = append(cards, deck.DrawCards(3)...)
	assert.Len(t, deck.cards, 9, "Deck should contain 9 cards")
	cards = append(cards, deck.DrawCards(2)...)
	assert.Len(t, deck.cards, 7, "Deck should contain 7 cards")
	cards = append(cards, deck.DrawCards(3)...)
	assert.Len(t, deck.cards, 4, "Deck should contain 4 cards")

	// we should have dealt 20 cards
	assert.Len(t, cards, 20, "Expected to have dealt 20 cards")

	deck.ReturnCards(&cards)
	assert.Len(t, deck.cards, 24, "Expected deck to have 24 cards")
}

func TestShuffle(t *testing.T) {
	deck := InitDeck(0)
	deck.Shuffle()
}
