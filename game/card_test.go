package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	// test both trump
	c1 := Card{rank: TEN, suite: DIAMOND}
	c2 := Card{rank: NINE, suite: DIAMOND}
	res := c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer
	c1 = Card{rank: JACK, suite: DIAMOND}
	c2 = Card{rank: ACE, suite: DIAMOND}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test both trump, c1 is right bauer, c2 is left
	c1 = Card{rank: JACK, suite: DIAMOND}
	c2 = Card{rank: JACK, suite: HEART}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c1 is trump, c2 is not
	c1 = Card{rank: NINE, suite: DIAMOND}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test c2 is trump, c1 is not
	c1 = Card{rank: NINE, suite: HEART}
	c2 = Card{rank: NINE, suite: DIAMOND}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test neither are trump but both lead
	c1 = Card{rank: TEN, suite: HEART}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c1 lead, c2 off-suite
	c1 = Card{rank: TEN, suite: HEART}
	c2 = Card{rank: NINE, suite: SPADE}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Positive(t, res, "expected c1 to be greater than c2")

	// test neither are trump, c2 lead, c1 off-suite
	c1 = Card{rank: TEN, suite: SPADE}
	c2 = Card{rank: NINE, suite: HEART}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Negative(t, res, "expected c2 to be greater than c1")

	// test both are off-suite
	c1 = Card{rank: TEN, suite: SPADE}
	c2 = Card{rank: NINE, suite: CLUB}
	res = c1.Compare(c2, DIAMOND, HEART)
	assert.Zero(t, res, "expected c2 to be equal to c1")
}
