package main

import (
	"math/rand"
	"time"
)

type Deck []Card

func InitDeck() Deck {
	var deck = make([]Card, 24)
	for i := 0; i < len(deck); i++ {
		deck[i].rank = GetRank(i % 6)
		deck[i].suite = GetSuite(i / 6)
	}
	return deck
}

func (d Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	n := len(d)
	for i := n; i >= 1; i-- {
		j := rng.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
}
