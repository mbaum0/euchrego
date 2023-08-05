package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDeck(t *testing.T) {
	deck := InitDeck()

	assert.Equal(t, len(deck), 24, "Deck should contain 24 cards")

	// make sure we have 6 of each suite
	hearts := 0
	clubs := 0
	diamonds := 0
	spades := 0
	for _, c := range deck {
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
