package game

type Game struct {
	State       GameState
	Deck        Deck
	Players     [4]Player
	DealerIndex int
	PlayerIndex int
	TurnedCard  *Card
	PlayedCards []*Card
	Trump       Suite
}

func (g *Game) TransitionState(newState GameState) {
	g.State = newState
	g.State.EnterState()
}

func (g *Game) NextPlayer() {
	g.PlayerIndex = (g.PlayerIndex + 1) % 4
}

func Run() {
	game := Game{}
	game.State = NewInitState()
	for {
		game.State.DoState(&game)

		if game.State.GetName() == EndGame {
			break
		}
	}
}
