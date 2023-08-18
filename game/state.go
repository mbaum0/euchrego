package game

import "fmt"

type StateName string

const (
	InitGame            StateName = "InitGame"
	DrawForDealer       StateName = "DrawForDealer"
	ResetDeckAndShuffle StateName = "ResetDeckAndShuffle"
	DealCards           StateName = "DealCards"
	RevealTopCard       StateName = "RevealTopCard"
	TrumpSelectionOne   StateName = "TrumpSelectionOne"
	PlayerPickupTrump   StateName = "PlayerPickupTrump"
	PlayerExchangeTrump StateName = "PlayerExchangeTrump"
	TrumpSelectionTwo   StateName = "TrumpSelectionTwo"
	ScrewDealer         StateName = "ScrewDealer"
	StartRound          StateName = "StartRound"
	GetPlayerCard       StateName = "GetPlayerCard"
	CheckValidCard      StateName = "CheckValidCard"
	PlayCard            StateName = "PlayCard"
	GetTrickWinner      StateName = "GetTrickWinner"
	GivePoints          StateName = "GivePoints"
	CheckForWinner      StateName = "CheckForWinner"
	EndGame             StateName = "EndGame"
)

type GameState interface {
	EnterState()
	DoState(game *Game)
	GetName() StateName
}

type NamedState struct {
	Name StateName
}

func (state *NamedState) GetName() StateName {
	return state.Name
}

func (state *NamedState) EnterState() {
	fmt.Printf("\n<-- Entered %s State -->\n", state.Name)
}

// ============================ InitGameState ============================
type InitGameState struct {
	NamedState
}

func NewInitState() *InitGameState {
	return &InitGameState{NamedState{Name: InitGame}}
}

func (state *InitGameState) DoState(game *Game) {
	game.Deck = InitDeck()
	game.Deck.Shuffle()
	game.Players[0] = InitPlayer("Player 1", 0)
	game.Players[1] = InitPlayer("Player 2", 1)
	game.Players[2] = InitPlayer("Player 3", 2)
	game.Players[3] = InitPlayer("Player 4", 3)
	game.DealerIndex = 0
	game.PlayerIndex = 0
	game.TurnedCard = nil
	game.Trump = NONE
	game.PlayedCards = make([]*Card, 0)

	game.TransitionState(NewDrawForDealerState())
}

// ============================ DrawForDealerState ============================
type DrawForDealerState struct {
	NamedState
}

func NewDrawForDealerState() *DrawForDealerState {
	return &DrawForDealerState{NamedState{Name: DrawForDealer}}
}

func (state *DrawForDealerState) DoState(game *Game) {
	// draw for a jack
	game.PlayedCards = append(game.PlayedCards, game.Deck.pop())
	lastIndex := len(game.PlayedCards) - 1

	// print drawn card
	fmt.Printf("%s was drawn!\n", game.PlayedCards[lastIndex].ToString())

	if game.PlayedCards[lastIndex].rank == JACK {
		// got trump. Set dealer and continue
		game.DealerIndex = game.PlayerIndex
		game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
		game.TransitionState(NewResetDeckAndShuffleState())
		fmt.Printf("Player %d is dealer\n", game.DealerIndex)
	} else {
		// no trump. Continue drawing
		game.NextPlayer()
	}

}

// ============================ ResetDeckAndShuffleState ============================
type ResetDeckAndShuffleState struct {
	NamedState
}

func NewResetDeckAndShuffleState() *ResetDeckAndShuffleState {
	return &ResetDeckAndShuffleState{NamedState{Name: ResetDeckAndShuffle}}
}

func (state *ResetDeckAndShuffleState) DoState(game *Game) {
	// reset deck
	game.Deck.ReturnCards(&game.PlayedCards)
	game.Deck.Shuffle()

	game.TransitionState(NewDealCardsState())
}

// ============================ DealCardsState ============================
type DealCardsState struct {
	NamedState
}

func NewDealCardsState() *DealCardsState {
	return &DealCardsState{NamedState{Name: DealCards}}
}

func (state *DealCardsState) DoState(game *Game) {
	// deal cards in standard euchre fashion

	dealerIndex := game.DealerIndex
	firstPlayerIndex := (dealerIndex + 1) % 4
	secondPlayerIndex := (dealerIndex + 2) % 4
	thirdPlayerIndex := (dealerIndex + 3) % 4

	game.Players[firstPlayerIndex].GiveCards(game.Deck.DrawCards(3))
	game.Players[secondPlayerIndex].GiveCards(game.Deck.DrawCards(2))
	game.Players[thirdPlayerIndex].GiveCards(game.Deck.DrawCards(3))
	game.Players[dealerIndex].GiveCards(game.Deck.DrawCards(2))

	game.Players[firstPlayerIndex].GiveCards(game.Deck.DrawCards(2))
	game.Players[secondPlayerIndex].GiveCards(game.Deck.DrawCards(3))
	game.Players[thirdPlayerIndex].GiveCards(game.Deck.DrawCards(2))
	game.Players[dealerIndex].GiveCards(game.Deck.DrawCards(3))

	game.TransitionState(NewRevealTopCardState())
}

// ============================ RevealTopCardState ============================
type RevealTopCardState struct {
	NamedState
}

func NewRevealTopCardState() *RevealTopCardState {
	return &RevealTopCardState{NamedState{Name: RevealTopCard}}
}

func (state *RevealTopCardState) DoState(game *Game) {
	game.TurnedCard = game.Deck.pop()

	// print out name of turned card
	fmt.Printf("%s was turned\n", game.TurnedCard.ToString())
	game.TransitionState(NewTrumpSelectionOneState())
}

// ============================ TrumpSelectionOneState ============================
type TrumpSelectionOneState struct {
	NamedState
}

func NewTrumpSelectionOneState() *TrumpSelectionOneState {
	return &TrumpSelectionOneState{NamedState{Name: TrumpSelectionOne}}
}

func (state *TrumpSelectionOneState) DoState(game *Game) {
	// ask player if they want trump
	player := game.Players[game.PlayerIndex]
	pickedUp := GetTrumpSelectionOneInput(player, *game.TurnedCard)

	// if picked up, we want to ask the dealer if they want the turned card
	if pickedUp {
		game.TransitionState(NewPlayerPickupTrumpState())
		return
	}

	// if this player was the dealer, we will move on to Trump Selection Two
	if game.PlayerIndex == game.DealerIndex {
		game.TransitionState(NewTrumpSelectionTwoState())
		return
	}

	// otherwise, move on to the next player
	game.NextPlayer()
}

// ============================ PlayerPickupTrumpState ============================
type PlayerPickupTrumpState struct {
	NamedState
}

func NewPlayerPickupTrumpState() *PlayerPickupTrumpState {
	return &PlayerPickupTrumpState{NamedState{Name: PlayerPickupTrump}}
}

func (state *PlayerPickupTrumpState) DoState(game *Game) {
	dealer := game.Players[game.DealerIndex]
	// ask the dealer if they want the turned card
	wantsIt := GetDealerWantsToPickItUp(dealer, *game.TurnedCard)

	// if dealer wants it, prompt them for a discard and give them the new card
	if wantsIt {
		game.TransitionState(NewPlayerExchangeTrumpState())
		return
	}
	game.TransitionState(NewStartRoundState())

}

// ============================ PlayerExhangeTrumpState ============================
type PlayerExchangeTrumpState struct {
	NamedState
}

func NewPlayerExchangeTrumpState() *PlayerExchangeTrumpState {
	return &PlayerExchangeTrumpState{NamedState{Name: PlayerExchangeTrump}}
}

func (state *PlayerExchangeTrumpState) DoState(game *Game) {
	dealer := game.Players[game.DealerIndex]
	burnCard := GetDealersBurnCard(dealer)
	game.Deck.ReturnCard(&burnCard)
	dealer.GiveCard(game.TurnedCard)
	game.TurnedCard = nil

	game.TransitionState(NewStartRoundState())
	return
}

// ============================ TrumpSelectionTwoState ============================
type TrumpSelectionTwoState struct {
	NamedState
}

func NewTrumpSelectionTwoState() *TrumpSelectionTwoState {
	return &TrumpSelectionTwoState{NamedState{Name: TrumpSelectionTwo}}
}

func (state *TrumpSelectionTwoState) DoState(game *Game) {
	player := game.Players[game.PlayerIndex]

	// if the player is the dealer, they must select a suite
	if game.PlayerIndex == game.DealerIndex {
		game.TransitionState(NewScrewDealerState())
		return
	}

	// otherwise, let the next player pick a suite if they want
	selectedSuite := GetTrumpSelectionTwoInput(player, *game.TurnedCard)

	// if the player selected a suite, set it as trump
	if selectedSuite != NONE {
		game.Trump = selectedSuite
		game.Deck.ReturnCard(&game.TurnedCard)
		game.TransitionState(NewStartRoundState())
		return
	}

	// move on to the next player
	game.NextPlayer()
}

// ============================ ScrewDealerState ============================
type ScrewDealerState struct {
	NamedState
}

func NewScrewDealerState() *ScrewDealerState {
	return &ScrewDealerState{NamedState{Name: ScrewDealer}}
}

func (state *ScrewDealerState) DoState(game *Game) {
	player := game.Players[game.PlayerIndex]
	selectedSuite := GetScrewTheDealerInput(player, *game.TurnedCard)

	game.Trump = selectedSuite
	game.Deck.ReturnCard(&game.TurnedCard)
	game.TransitionState(NewStartRoundState())
}

// ============================ StartRoundState ============================
type StartRoundState struct {
	NamedState
}

func NewStartRoundState() *StartRoundState {

	return &StartRoundState{NamedState{Name: StartRound}}
}

func (state *StartRoundState) DoState(game *Game) {
	game.PlayerIndex = game.DealerIndex + 1
	game.TransitionState(NewGameOverState())
}

// ============================ GameOverState ============================
type GameOverState struct {
	NamedState
}

func NewGameOverState() *GameOverState {
	return &GameOverState{NamedState{Name: EndGame}}
}

func (state *GameOverState) DoState(game *Game) {
	return
}
