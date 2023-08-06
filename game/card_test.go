package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Testcompare(t *testing.T) {
	// test both trump
	c1 := Card{rank: TEN, suite: DIAMOND}
	c2 := Card{rank: NINE, suite: DIAMOND}
	res := c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer
	c1 = Card{rank: JACK, suite: DIAMOND}
	c2 = Card{rank: ACE, suite: DIAMOND}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer, c2 is left
	c1 = Card{rank: JACK, suite: DIAMOND}
	c2 = Card{rank: JACK, suite: HEART}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c1 is trump, c2 is not
	c1 = Card{rank: NINE, suite: DIAMOND}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c2 is trump, c1 is not
	c1 = Card{rank: NINE, suite: HEART}
	c2 = Card{rank: NINE, suite: DIAMOND}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test neither are trump but both lead
	c1 = Card{rank: TEN, suite: HEART}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c1 lead, c2 off-suite
	c1 = Card{rank: TEN, suite: HEART}
	c2 = Card{rank: NINE, suite: SPADE}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c2 lead, c1 off-suite
	c1 = Card{rank: TEN, suite: SPADE}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test both are off-suite
	c1 = Card{rank: TEN, suite: SPADE}
	c2 = Card{rank: NINE, suite: CLUB}
	res = c1.compare(c2, DIAMOND, HEART)
	assert.Zero(t, res, "expected c2 to be equal to c1")
}

func TestGetPlayableCards(t *testing.T) {

	// test player has at least one trump card
	var cards = make([]*Card, 0)
	cards = append(cards, &Card{rank: NINE, suite: CLUB})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: CLUB})

	playableCards := GetPlayableCards(cards, DIAMOND, HEART)
	assert.Equal(t, len(playableCards), 1, "Expected 1 card to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	cards = cards[:0]

	// test player has at least one lead card and no trump
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: SPADE})

	playableCards = GetPlayableCards(cards, HEART, DIAMOND)
	assert.Equal(t, len(playableCards), 2, "Expected 2 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *playableCards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	cards = cards[:0]

	// test player has no trump or lead cards
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: SPADE})

	playableCards = GetPlayableCards(cards, HEART, HEART)
	assert.Equal(t, len(playableCards), 4, "Expected 4 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: NINE, suite: DIAMOND}, "Expected nine of diamonds to be returned")
	assert.Equal(t, *playableCards[1], Card{rank: TEN, suite: DIAMOND}, "Expected ten of diamonds to be returned")
	assert.Equal(t, *playableCards[2], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
	assert.Equal(t, *playableCards[3], Card{rank: KING, suite: SPADE}, "Expected king of spades to be returned")
	cards = cards[:0]

	// player has only the left bauer and no lead card
	cards = append(cards, &Card{rank: NINE, suite: DIAMOND})
	cards = append(cards, &Card{rank: TEN, suite: DIAMOND})
	cards = append(cards, &Card{rank: JACK, suite: SPADE})
	cards = append(cards, &Card{rank: KING, suite: DIAMOND})
	playableCards = GetPlayableCards(cards, CLUB, HEART)
	assert.Equal(t, len(playableCards), 1, "Expected 1 cards to be returned")
	assert.Equal(t, *playableCards[0], Card{rank: JACK, suite: SPADE}, "Expected jack of spades to be returned")
}
