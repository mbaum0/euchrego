package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayableCards(t *testing.T) {
	p := InitPlayer("Player 1")

	// test player has at least one trump card
	var cards = make([]*Card, 0)
	cards = append(cards, &Card{rank: NINE, suite: CLUB})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: CLUB})
	p.GiveCards(cards)
	playableCards := p.GetPlayableCards(DIAMOND, HEART)
	assert.Equal(t, len(playableCards), 1, "Expected 1 card to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	p.ReturnCards()
	cards = cards[:0]

	// test player has at least one lead card and no trump
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: SPADE})
	p.GiveCards(cards)
	playableCards = p.GetPlayableCards(HEART, DIAMOND)
	assert.Equal(t, len(playableCards), 2, "Expected 2 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *playableCards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	p.ReturnCards()
	cards = cards[:0]

	// test player has no trump or lead cards
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: SPADE})
	p.GiveCards(cards)
	playableCards = p.GetPlayableCards(HEART, HEART)
	assert.Equal(t, len(playableCards), 4, "Expected 4 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *playableCards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	assert.Equal(t, *playableCards[2], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
	assert.Equal(t, *playableCards[3], Card{rank: KING, suite: SPADE}, "Expected king of spades to be returned")
	p.ReturnCards()
	cards = cards[:0]

	// player has only the left bauer and no lead card
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: DIAMOND})
	p.GiveCards(cards)
	playableCards = p.GetPlayableCards(CLUB, HEART)
	assert.Equal(t, len(playableCards), 1, "Expected 1 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
	p.ReturnCards()
}
