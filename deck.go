package main

type Deck []Card

func InitDeck() Deck {
	var deck = make([]Card, 24)
	for i := 0; i < len(deck); i++ {
		deck[i].rank = GetRank(i % 6)
		deck[i].suite = GetSuite(i / 6)
	}
	return deck
}
