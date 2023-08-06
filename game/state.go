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
		InitPlayer("Player 1", 0),
		InitPlayer("Player 2", 1),
		InitPlayer("Player 3", 2),
		InitPlayer("Player 4", 3))
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
		game.TransitionState(NewShuffleState())
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
	drawCardEvent := NewDrawnCardEvent(drawnCard, &game.Players[len(game.PlayedCards)%4])
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

type shuffleState struct {
	DefaultGameState
}

func NewShuffleState() *shuffleState {
	state := &shuffleState{}
	state.StateName = "SHUFFLE_STATE"
	return state
}

func (state *shuffleState) DoState(game *Game, event Event) Event {
	game.Deck.Shuffle()
	return NewEmptyEvent()
}

func (state *shuffleState) NextState(game *Game, event Event) {
	game.TransitionState(NewDealCardsState())
}

type dealCardsState struct {
	DefaultGameState
}

func NewDealCardsState() *dealCardsState {
	state := &dealCardsState{}
	state.StateName = "DEAL_CARDS_STATE"
	return state
}

func (state *dealCardsState) DoState(game *Game, event Event) Event {
	// if we finished dealing
	if len(game.Players[game.DealerIndex].hand) == 5 {
		return NewFinishedDealingEvent()
	}

	// start dealing to first player if we got an empty event
	if _, ok := event.(EmptyEvent); ok {
		dealtCards := game.Deck.DrawCards(3) // first deal is three cards
		currentPlayerIndex := game.DealerIndex + 1
		if currentPlayerIndex > 3 {
			currentPlayerIndex = 0
		}
		game.Players[currentPlayerIndex].GiveCards(dealtCards)
		return NewDealtCardsEvent(dealtCards, &game.Players[currentPlayerIndex])
	}

	// if we're in the middle of dealing
	if event, ok := event.(DealtCardsEvent); ok {

		// if we've finished dealing
		if len(game.Deck.cards) == 4 {
			return NewFinishedDealingEvent()
		}

		currentPlayerIndex := event.player.index + 1
		if currentPlayerIndex > 3 {
			currentPlayerIndex = 0
		}
		currentPlayer := &game.Players[currentPlayerIndex]

		numToDeal := 2
		// if this is current players first set of cards, use previous deal to determine amount
		if len(currentPlayer.hand) == 0 {
			if len(event.DealtCards) == 2 {
				numToDeal = 3
			}
		} else {
			// otherwise, deal remainder of cards
			numToDeal = 5 - len(currentPlayer.hand)
		}
		dealtCards := game.Deck.DrawCards(numToDeal)
		currentPlayer.GiveCards(dealtCards)
		return NewDealtCardsEvent(dealtCards, currentPlayer)
	}

	panic("Got bad event!")
}

func (state *dealCardsState) NextState(game *Game, event Event) {
	if _, ok := event.(FinishedDealingEvent); ok {
		game.TransitionState(NewEndGameState())

	} else {
		game.TransitionState(NewDealCardsState())
	}

}
