package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayableCards(t *testing.T) {
	p := InitPlayer()

	// test player has at least one trump card
	p.GiveCard(&Card{rank: NINE, suite: CLUB})
	p.GiveCard(&Card{rank: TEN, suite: DIAMOND})
	p.GiveCard(&Card{rank: JACK, suite: SPADE})
	p.GiveCard(&Card{rank: KING, suite: CLUB})
	cards := p.GetPlayableCards(DIAMOND, HEART)
	assert.Equal(t, len(cards), 1, "Expected 1 card to be returned")
	assert.Equal(t, *cards[0], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	p.ReturnCards()

	// test player has at least one lead card and no trump
	p.GiveCard(&Card{rank: NINE, suite: DIAMOND})
	p.GiveCard(&Card{rank: TEN, suite: DIAMOND})
	p.GiveCard(&Card{rank: JACK, suite: SPADE})
	p.GiveCard(&Card{rank: KING, suite: SPADE})
	cards = p.GetPlayableCards(HEART, DIAMOND)
	assert.Equal(t, len(cards), 2, "Expected 2 cards to be returned")
	assert.Equal(t, *cards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *cards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	p.ReturnCards()

	// test player has no trump or lead cards
	p.GiveCard(&Card{rank: NINE, suite: DIAMOND})
	p.GiveCard(&Card{rank: TEN, suite: DIAMOND})
	p.GiveCard(&Card{rank: JACK, suite: SPADE})
	p.GiveCard(&Card{rank: KING, suite: SPADE})
	cards = p.GetPlayableCards(HEART, HEART)
	assert.Equal(t, len(cards), 4, "Expected 4 cards to be returned")
	assert.Equal(t, *cards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *cards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	assert.Equal(t, *cards[2], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
	assert.Equal(t, *cards[3], Card{rank: KING, suite: SPADE}, "Expected king of spades to be returned")
	p.ReturnCards()

	// player has only the left bauer and no lead card
	p.GiveCard(&Card{rank: NINE, suite: DIAMOND})
	p.GiveCard(&Card{rank: TEN, suite: DIAMOND})
	p.GiveCard(&Card{rank: JACK, suite: SPADE})
	p.GiveCard(&Card{rank: KING, suite: DIAMOND})
	cards = p.GetPlayableCards(CLUB, HEART)
	assert.Equal(t, len(cards), 1, "Expected 1 cards to be returned")
	assert.Equal(t, *cards[0], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
	p.ReturnCards()

}
