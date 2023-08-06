package game

import "fmt"

type GameState interface {
	NextState(game *Game, event Event)
	EnterState()
	GetName() string
	DoState(game *Game, event Event) Event
}

type DefaultGameState struct {
	StateName string
}

func (state *DefaultGameState) GetName() string {
	return state.StateName
}

func (state *DefaultGameState) EnterState() {
	fmt.Println("changed state to:", state.StateName)
}

type initState struct {
	DefaultGameState
}

func NewInitState() *initState {
	state := &initState{}
	state.StateName = "INIT"
	return state
}

func (state *initState) DoState(game *Game, event Event) Event {
	game.Deck = InitDeck()
	game.Players = make([]Player, 0)
	game.Players = append(game.Players,
		InitPlayer("Player 1"),
		InitPlayer("Player 2"),
		InitPlayer("Player 3"),
		InitPlayer("Player 4"))
	game.CurrentTrump = NONE
	game.PlayedCards = make([]*Card, 0)
	game.DealerIndex = 0
	return NewEmptyEvent()
}

func (state *initState) NextState(game *Game, event Event) {

	game.TransitionState(NewDetermineFirstDealerState())
}

type determineFirstDealerState struct {
	DefaultGameState
}

func NewDetermineFirstDealerState() *determineFirstDealerState {
	state := &determineFirstDealerState{}
	state.StateName = "DETERMINE_FIRST_DEALER_STATE"
	return state
}

func (state *determineFirstDealerState) DoState(game *Game, event Event) Event {
	if _, ok := event.(EmptyEvent); ok {
		// shuffle to start looking for jacks
		game.Deck.Shuffle()
	} else if drawnCardEvent, ok := event.(DrawnCardEvent); ok {
		game.PlayedCards = append(game.PlayedCards, drawnCardEvent.DrawnCard)
		// last event was a card draw
		if drawnCardEvent.DrawnCard.rank == JACK {
			game.DealerIndex = len(game.PlayedCards) % 4
			game.Deck.ReturnCards(&game.PlayedCards)
			return NewDrawnJackEvent()
		} else {
			return event
		}
	}
	return NewEmptyEvent()
}

func (state *determineFirstDealerState) NextState(game *Game, event Event) {
	if _, ok := event.(EmptyEvent); ok {
		game.TransitionState(NewDrawForJackState())
		return
	}

	if _, ok := event.(DrawnCardEvent); ok {
		game.TransitionState(NewDrawForJackState())
		return
	}

	if _, ok := event.(DrawnJackEvent); ok {
		game.TransitionState(NewEndGameState()) // todo
		return
	}
}

type drawForJackState struct {
	DefaultGameState
}

func NewDrawForJackState() *drawForJackState {
	state := &drawForJackState{}
	state.StateName = "DRAW_FOR_JACK_STATE"
	return state
}

func (state *drawForJackState) DoState(game *Game, event Event) Event {
	drawnCard := game.Deck.DrawCards(1)[0]
	drawCardEvent := DrawnCardEvent{}
	drawCardEvent.DrawnCard = drawnCard
	return drawCardEvent
}

func (state *drawForJackState) NextState(game *Game, event Event) {
	game.TransitionState(NewDetermineFirstDealerState())
}

type endGameState struct {
	DefaultGameState
}

func NewEndGameState() *endGameState {
	state := &endGameState{}
	state.StateName = "END_GAME_STATE"
	return state
}

func (state *endGameState) DoState(game *Game, event Event) Event {
	return NewEndGameEvent()
}

func (state *endGameState) NextState(game *Game, event Event) {
	game.TransitionState(NewEndGameState())
}
