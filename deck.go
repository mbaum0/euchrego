package main

import (
	"math/rand"
	"time"
)

type Deck []*Card

func InitDeck() Deck {
	var deck = make([]*Card, 0)
	for i := 0; i < 24; i++ {
		c := Card{rank: GetRank(i % 6), suite: GetSuite(i / 6)}
		deck = append(deck, &c)
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
