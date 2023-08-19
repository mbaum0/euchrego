package game

import (
	"fmt"
	"time"
)

type Game struct {
	State              GameState
	Deck               Deck
	Players            [4]*Player
	DealerIndex        int
	PlayerIndex        int
	TurnedCard         *Card
	PlayedCards        []*Card
	Trump              Suite
	OrderedPlayerIndex int // the player who ordered it up
	logs               []string
}

func NewGame() Game {
	game := Game{}
	game.State = NewInitState()
	game.PlayedCards = nil
	game.logs = make([]string, 0)
	return game
}

func (g *Game) Log(format string, args ...interface{}) {
	g.logs = append(g.logs, fmt.Sprintf(format, args...))
}

func (g *Game) TransitionState(newState GameState) {
	g.State = newState
	g.State.EnterState()
}

func (g *Game) PlayCard(card *Card) {
	g.PlayedCards = append(g.PlayedCards, card)
}

func (g *Game) ReturnPlayedCards() {
	g.Deck.ReturnCards(&g.PlayedCards)
}

func (g *Game) NextPlayer() {
	g.PlayerIndex = (g.PlayerIndex + 1) % 4
}

func Run() {
	game := NewGame()
	display := NewTextDisplay(220, 60)
	game.State = NewInitState()
	for {
		display.DrawBoard(&game)
		game.State.DoState(&game)
		// delay for 1 second for animation
		time.Sleep(1 * time.Second)

		if game.State.GetName() == EndGame {
			break
		}
	}
}
