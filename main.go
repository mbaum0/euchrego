package main

import (
	"fmt"

	"github.com/mbaum0/euchrego/game"
)

type Game struct {
	playedCards []*game.Card
	players     []*game.Player
}

// func (g *Game) doPlayerTurn(p *game.Player) {
// 	if len(g.playedCards) == 0 {
// 		// new deal, we can play anything
// 		p.GetPlayableCards()
// 	}
// }

func (g *Game) determineTrump(flippedSuite game.Suite) game.Suite {
	var trumpSuite game.Suite
	for _, player := range g.players {
		fmt.Printf("%s pick trump!\n", player.GetName())
		player.PrintHand()
		trumpSuite = game.GetSuiteInput(flippedSuite, game.NONE)
		if trumpSuite == flippedSuite {
			return flippedSuite
		}
	}
	// if we got here, no trump was picked.
	suites := []game.Suite{game.CLUB, game.DIAMOND, game.HEART, game.SPADE}
	allowedSuites := make([]game.Suite, 0)
	for _, s := range suites {
		if s != flippedSuite {
			allowedSuites = append(allowedSuites, s)
		}
	}

	allowedSuitesWithNone := make([]game.Suite, 0)
	allowedSuitesWithNone = append(allowedSuitesWithNone, allowedSuites...)
	allowedSuitesWithNone = append(allowedSuitesWithNone, game.NONE)

	for i, player := range g.players {
		if i == 3 {
			// last player must pick!
			fmt.Printf("%s must pick trump!\n", player.GetName())
			player.PrintHand()
			trumpSuite = game.GetSuiteInput(allowedSuites...)
		} else {
			fmt.Printf("%s pick trump!\n", player.GetName())
			player.PrintHand()
			trumpSuite = game.GetSuiteInput(allowedSuitesWithNone...)
			if trumpSuite != game.NONE {
				return trumpSuite
			}
		}
	}
	return trumpSuite
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

	game.players = append(game.players, &player1, &player2, &player3, &player4)

	// flip card for trump
	flippedCard := deck.DrawCards(1)[0]
	fmt.Printf("%s was drawn!\n", flippedCard.Info())

	for teamOneScore < 10 && teamTwoScore < 10 {
		trump := game.determineTrump(flippedCard.GetSuite())
		fmt.Printf("%s are trump!\n", trump.GetString())
	}
}
