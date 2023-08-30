package game

import "github.com/mbaum0/euchrego/godeck"

type GameBoard struct {
	Deck               *godeck.EuchreDeck
	Players            [4]*Player
	DealerIndex        int
	PlayerIndex        int
	TurnedCard         godeck.Card
	PlayedCards        []godeck.Card
	Trump              godeck.Suit
	OrderedPlayerIndex int // the player who ordered it up
	logs               []string
	RandSeed           int64
}

func NewGameBoard() GameBoard {
	board := GameBoard{}
	board.PlayedCards = nil
	board.logs = make([]string, 0)
	board.OrderedPlayerIndex = -1
	board.DealerIndex = 0
	board.PlayerIndex = 0
	board.RandSeed = int64(1)
	board.Players[0] = InitPlayer("Player 1", 0)
	board.Players[1] = InitPlayer("Player 2", 1)
	board.Players[2] = InitPlayer("Player 3", 2)
	board.Players[3] = InitPlayer("Player 4", 3)
	board.TurnedCard = godeck.EmptyCard()
	board.Trump = godeck.None
	board.PlayedCards = make([]godeck.Card, 0)
	return board
}

func (g *GameBoard) Log(format string, args ...interface{}) {
	g.logs = append(g.logs, format)
}

func (g *GameBoard) PlayCard(card godeck.Card) {
	g.PlayedCards = append(g.PlayedCards, card)
}

func (g *GameBoard) ReturnPlayedCards() {
	g.Deck.ReturnCards(g.PlayedCards)
	// clear played cards
	g.PlayedCards = g.PlayedCards[:0]
}

func (g *GameBoard) NextPlayer() {
	g.PlayerIndex = (g.PlayerIndex + 1) % 4
}
