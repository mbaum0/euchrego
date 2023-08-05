package main

import "github.com/mbaum0/euchrego/game"

func main() {
	player1 := game.InitPlayer("Player 1")
	player2 := game.InitPlayer("Player 2")
	player3 := game.InitPlayer("Player 3")
	player4 := game.InitPlayer("Player 4")
	deck := game.InitDeck()

	deck.Shuffle()

	player1.GiveCards(deck.DrawCards(2))
	player2.GiveCards(deck.DrawCards(3))
	player3.GiveCards(deck.DrawCards(2))
	player4.GiveCards(deck.DrawCards(3))

	player1.GiveCards(deck.DrawCards(3))
	player2.GiveCards(deck.DrawCards(2))
	player3.GiveCards(deck.DrawCards(3))
	player4.GiveCards(deck.DrawCards(2))

	player1.PrintHand()
	player2.PrintHand()
	player3.PrintHand()
	player4.PrintHand()
}
