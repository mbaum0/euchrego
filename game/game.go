package game

type Game struct {
	State              GameState
	Deck               Deck
	Players            []Player
	DealerIndex        int
	CurrentTrump       Suite
	PlayedCards        []*Card
	PotentialTrumpCard *Card
	BurnedTrumpSuite   Suite
	CurrentPlayerIndex int
}

func (g *Game) TransitionState(newState GameState) {
	g.State = newState
	g.State.EnterState()
}

func Run() {
	game := Game{}
	game.State = NewInitState()
	var event Event = nil
	for {
		event = game.State.DoState(&game, event)
		game.State.NextState(&game, event)

		if _, ok := event.(EndGameEvent); ok {
			break
		}
	}
}
