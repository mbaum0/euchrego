package main

import (
	"fmt"

	"github.com/mbaum0/euchrego/game"
)

type Game struct {
	playedCards []*game.Card
}

// func (g *Game) doPlayerTurn(p *game.Player) {
// 	if len(g.playedCards) == 0 {
// 		// new deal, we can play anything
// 		p.GetPlayableCards()
// 	}
// }

func (g *Game) determineTrump() {
	s := game.GetSuiteInput(game.DIAMOND, game.CLUB, game.HEART, game.SPADE, game.NONE)
	fmt.Printf("you entered: %d\n", s)
}

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

	teamOneScore := 0
	teamTwoScore := 0

	var playedCards = make([]*game.Card, 0)
	game := Game{playedCards: playedCards}

	for teamOneScore < 10 && teamTwoScore < 10 {
		game.determineTrump()
	}
}
