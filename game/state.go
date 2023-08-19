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
	DealerPickupTrump   StateName = "DealerPickupTrump"
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

type StateMachine struct {
	CurrentState GameState
}

func NewStateMachine() StateMachine {
	sm := StateMachine{}
	sm.CurrentState = NewInitState()
	return sm
}

func (sm *StateMachine) TransitionState(newStateName StateName) {
	// check if current state is allowed to transition to the newState
	if !sm.CurrentState.CanTransitionTo(newStateName) {
		errMsg := fmt.Sprintf("Cannot transition from %s to %s", sm.CurrentState.GetName(), newStateName)
		panic(errMsg)
	}

	// create the new state and enter it
	switch newStateName {
	case InitGame:
		sm.CurrentState = NewInitState()
	case DrawForDealer:
		sm.CurrentState = NewDrawForDealerState()
	case ResetDeckAndShuffle:
		sm.CurrentState = NewResetDeckAndShuffleState()
	case DealCards:
		sm.CurrentState = NewDealCardsState()
	case RevealTopCard:
		sm.CurrentState = NewRevealTopCardState()
	case TrumpSelectionOne:
		sm.CurrentState = NewTrumpSelectionOneState()
	case DealerPickupTrump:
		sm.CurrentState = NewDealerPickupTrumpState()
	case TrumpSelectionTwo:
		sm.CurrentState = NewTrumpSelectionTwoState()
	case ScrewDealer:
		sm.CurrentState = NewScrewDealerState()
	case StartRound:
		sm.CurrentState = NewStartRoundState()
	case GetPlayerCard:
		sm.CurrentState = NewGetPlayerCardState()
	case CheckValidCard:
		sm.CurrentState = NewCheckValidCardState()
	case PlayCard:
		sm.CurrentState = NewPlayCardState()
	case GetTrickWinner:
		sm.CurrentState = NewGetTrickWinnerState()
	case GivePoints:
		sm.CurrentState = NewGivePointsState()
	case CheckForWinner:
		sm.CurrentState = NewCheckForWinnerState()
	case EndGame:
		sm.CurrentState = NewEndGameState()
	}
}

func (sm *StateMachine) Step(game *Game) {
	nextStateName := sm.CurrentState.DoState(game)
	sm.TransitionState(nextStateName)
}

type GameState interface {
	EnterState()
	DoState(game *Game) StateName
	GetName() StateName
	CanTransitionTo(newState StateName) bool
}

type NamedState struct {
	Name               StateName
	PossibleNextStates []StateName
}

func (state *NamedState) GetName() StateName {
	return state.Name
}

func (state *NamedState) EnterState() {
	//fmt.Printf("\n<-- Entered %s State -->\n", state.Name)
}

func (state *NamedState) CanTransitionTo(newStateName StateName) bool {
	for _, nextStateName := range state.PossibleNextStates {
		if nextStateName == newStateName {
			return true
		}
	}
	return false
}

// ============================ InitGameState ============================
type InitGameState struct {
	NamedState
}

func NewInitState() *InitGameState {
	gs := InitGameState{NamedState{Name: InitGame}}
	gs.PossibleNextStates = []StateName{DrawForDealer}
	return &gs
}

func (state *InitGameState) DoState(game *Game) StateName {
	game.Deck = InitDeck(game.RandSeed)
	game.Deck.Shuffle()
	return DrawForDealer
}

// ============================ DrawForDealerState ============================
type DrawForDealerState struct {
	NamedState
}

func NewDrawForDealerState() *DrawForDealerState {
	gs := DrawForDealerState{NamedState{Name: DrawForDealer}}
	gs.PossibleNextStates = []StateName{DrawForDealer, ResetDeckAndShuffle}
	return &gs
}

func (state *DrawForDealerState) DoState(game *Game) StateName {
	// draw for a jack
	game.PlayedCards = append(game.PlayedCards, game.Deck.pop())
	lastIndex := len(game.PlayedCards) - 1

	// print drawn card
	game.Log("%s was drawn", game.PlayedCards[lastIndex].ToString())

	if game.PlayedCards[lastIndex].rank == JACK {
		// got trump. Set dealer and continue
		game.DealerIndex = game.PlayerIndex
		game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
		dealer := game.Players[game.DealerIndex]
		game.Log("%s is dealer", dealer.name)
		game.Deck.ReturnCards(&game.PlayedCards)
		return ResetDeckAndShuffle
	}

	// no trump. Continue drawing
	game.NextPlayer()
	return DrawForDealer
}

// ============================ ResetDeckAndShuffleState ============================
type ResetDeckAndShuffleState struct {
	NamedState
}

func NewResetDeckAndShuffleState() *ResetDeckAndShuffleState {
	gs := ResetDeckAndShuffleState{NamedState{Name: ResetDeckAndShuffle}}
	gs.PossibleNextStates = []StateName{DealCards}
	return &gs
}

func (state *ResetDeckAndShuffleState) DoState(game *Game) StateName {
	// reset deck
	game.Deck.Shuffle()

	return DealCards
}

// ============================ DealCardsState ============================
type DealCardsState struct {
	NamedState
}

func NewDealCardsState() *DealCardsState {
	gs := DealCardsState{NamedState{Name: DealCards}}
	gs.PossibleNextStates = []StateName{DealCards, RevealTopCard}
	return &gs
}

func (state *DealCardsState) DoState(game *Game) StateName {
	// deal cards in standard euchre fashion
	dealerIndex := game.DealerIndex
	dealer := game.Players[dealerIndex]

	playerIndex := game.PlayerIndex
	player := game.Players[playerIndex]

	isFirstDeal := len(dealer.hand) == 0

	if isFirstDeal {
		// deal 2 cards to the player if they are the 1st or 3rd player
		if playerIndex == (dealerIndex+1)%4 || playerIndex == (dealerIndex+3)%4 {
			player.GiveCards(game.Deck.DrawCards(2))
			game.Log("%s was dealt 2 cards", player.name)
		} else {
			player.GiveCards(game.Deck.DrawCards(3))
			game.Log("%s was dealt 3 cards", player.name)
		}
	} else {
		// deal 3 cards to the player if they are the 1st or 3rd player
		if playerIndex == (dealerIndex+1)%4 || playerIndex == (dealerIndex+3)%4 {
			player.GiveCards(game.Deck.DrawCards(3))
			game.Log("%s was dealt 3 cards", player.name)
		} else {
			player.GiveCards(game.Deck.DrawCards(2))
			game.Log("%s was dealt 2 cards", player.name)
		}
	}

	// move onto next player
	game.NextPlayer()

	// if the dealer has all their cards, continue to RevealTopCardState
	if len(dealer.hand) == 5 {
		return RevealTopCard
	}
	return DealCards
}

// ============================ RevealTopCardState ============================
type RevealTopCardState struct {
	NamedState
}

func NewRevealTopCardState() *RevealTopCardState {
	gs := RevealTopCardState{NamedState{Name: RevealTopCard}}
	gs.PossibleNextStates = []StateName{TrumpSelectionOne}
	return &gs
}

func (state *RevealTopCardState) DoState(game *Game) StateName {
	game.TurnedCard = game.Deck.pop()

	// print out name of turned card
	game.Log("%s was turned", game.TurnedCard.ToString())
	return TrumpSelectionOne
}

// ============================ TrumpSelectionOneState ============================
type TrumpSelectionOneState struct {
	NamedState
}

func NewTrumpSelectionOneState() *TrumpSelectionOneState {
	gs := TrumpSelectionOneState{NamedState{Name: TrumpSelectionOne}}
	gs.PossibleNextStates = []StateName{TrumpSelectionOne, DealerPickupTrump, TrumpSelectionTwo}
	return &gs
}

func (state *TrumpSelectionOneState) DoState(game *Game) StateName {

	player := game.Players[game.PlayerIndex]

	// ask player if they want trump
	pickedUp := GetTrumpSelectionOneInput(player, *game.TurnedCard)

	// if picked up, we want to ask the dealer if they want the turned card
	if pickedUp {
		game.Log("%s ordered it up", player.name)
		game.OrderedPlayerIndex = game.PlayerIndex
		game.Trump = game.TurnedCard.suite
		return DealerPickupTrump
	}

	// if this player was the dealer, we will move on to Trump Selection Two
	if game.PlayerIndex == game.DealerIndex {
		game.NextPlayer()
		return TrumpSelectionTwo
	}

	// otherwise, move on to the next player
	game.NextPlayer()
	return TrumpSelectionOne
}

// ============================ DealerPickupTrumpState ============================
type DealerPickupTrumpState struct {
	NamedState
}

func NewDealerPickupTrumpState() *DealerPickupTrumpState {
	gs := DealerPickupTrumpState{NamedState{Name: DealerPickupTrump}}
	gs.PossibleNextStates = []StateName{StartRound}
	return &gs
}

func (state *DealerPickupTrumpState) DoState(game *Game) StateName {
	dealer := game.Players[game.DealerIndex]
	// give the dealer the turned card and let them exchange
	dealer.GiveCard(game.TurnedCard)
	burnCard := GetDealersBurnCard(dealer)
	dealer.ReturnCard(burnCard)
	game.Deck.ReturnCard(burnCard)
	game.TurnedCard = nil
	game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
	return StartRound
}

// ============================ TrumpSelectionTwoState ============================
type TrumpSelectionTwoState struct {
	NamedState
}

func NewTrumpSelectionTwoState() *TrumpSelectionTwoState {
	gs := TrumpSelectionTwoState{NamedState{Name: TrumpSelectionTwo}}
	gs.PossibleNextStates = []StateName{TrumpSelectionTwo, StartRound, ScrewDealer}
	return &gs
}

func (state *TrumpSelectionTwoState) DoState(game *Game) StateName {
	player := game.Players[game.PlayerIndex]

	// if the player is the dealer, they must select a suite
	if game.PlayerIndex == game.DealerIndex {
		game.Log("Dealer got screwed!")
		return ScrewDealer
	}

	// otherwise, let the next player pick a suite if they want
	selectedSuite := GetTrumpSelectionTwoInput(player, *game.TurnedCard)

	// if the player selected a suite, set it as trump
	if selectedSuite != NONE {
		game.Trump = selectedSuite
		game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
		game.Deck.ReturnCard(game.TurnedCard)
		game.TurnedCard = nil
		game.Log("%s picked %s as trump", player.name, selectedSuite.ToString())
		return StartRound
	}

	// move on to the next player
	game.NextPlayer()
	return TrumpSelectionTwo
}

// ============================ ScrewDealerState ============================
type ScrewDealerState struct {
	NamedState
}

func NewScrewDealerState() *ScrewDealerState {
	gs := ScrewDealerState{NamedState{Name: ScrewDealer}}
	gs.PossibleNextStates = []StateName{StartRound}
	return &gs
}

func (state *ScrewDealerState) DoState(game *Game) StateName {
	player := game.Players[game.PlayerIndex]

	selectedSuite := GetScrewTheDealerInput(player, *game.TurnedCard)

	game.Log("Dealer %s picked %s as trump", player.name, selectedSuite.ToString())

	game.Trump = selectedSuite
	game.Deck.ReturnCard(game.TurnedCard)
	game.TurnedCard = nil
	return StartRound
}

// ============================ StartRoundState ============================
type StartRoundState struct {
	NamedState
}

func NewStartRoundState() *StartRoundState {

	gs := StartRoundState{NamedState{Name: StartRound}}
	gs.PossibleNextStates = []StateName{GetPlayerCard}
	return &gs
}

func (state *StartRoundState) DoState(game *Game) StateName {
	game.PlayerIndex = (game.DealerIndex + 1) % 4
	return GetPlayerCard
}

// ============================ GetPlayerCardState ============================
type GetPlayerCardState struct {
	NamedState
}

func NewGetPlayerCardState() *GetPlayerCardState {
	gs := GetPlayerCardState{NamedState{Name: GetPlayerCard}}
	gs.PossibleNextStates = []StateName{CheckValidCard}
	return &gs
}

func (state *GetPlayerCardState) DoState(game *Game) StateName {
	player := game.Players[game.PlayerIndex]

	player.playedCard = GetCardInput(player)
	return CheckValidCard
}

// ============================ CheckValidCardState ============================
type CheckValidCardState struct {
	NamedState
}

func NewCheckValidCardState() *CheckValidCardState {
	gs := CheckValidCardState{NamedState{Name: CheckValidCard}}
	gs.PossibleNextStates = []StateName{GetPlayerCard, PlayCard}
	return &gs
}

func (state *CheckValidCardState) DoState(game *Game) StateName {
	player := game.Players[game.PlayerIndex]

	var leadCard *Card = nil

	if len(game.PlayedCards) > 0 {
		leadCard = game.PlayedCards[0]
	}

	// if the card wasn't valid, go back to GetPlayerCardState
	if !IsCardPlayable(player.playedCard, player.hand, game.Trump, leadCard) {
		game.Log("Invalid card. You must follow suite.")
		player.playedCard = nil
		return GetPlayerCard
	}

	// if card is valid, move on to play it
	return PlayCard
}

// ============================ PlayCardState ============================
type PlayCardState struct {
	NamedState
}

func NewPlayCardState() *PlayCardState {
	gs := PlayCardState{NamedState{Name: PlayCard}}
	gs.PossibleNextStates = []StateName{GetTrickWinner, GetPlayerCard}
	return &gs
}

func (state *PlayCardState) DoState(game *Game) StateName {
	player := game.Players[game.PlayerIndex]

	// remove the card from the players hand
	player.ReturnCard(player.playedCard)

	// add the card to the played cards
	game.PlayCard(player.playedCard)

	// print the card
	game.Log("%s played %s", player.name, player.playedCard.ToString())

	// if this is the last card, move on to GetTrickWinnerState
	if len(game.PlayedCards) == 4 {
		return GetTrickWinner
	}

	// move on to the next player
	game.NextPlayer()
	return GetPlayerCard
}

// ============================ GetTrickWinnerState ============================
type GetTrickWinnerState struct {
	NamedState
}

func NewGetTrickWinnerState() *GetTrickWinnerState {
	gs := GetTrickWinnerState{NamedState{Name: GetTrickWinner}}
	gs.PossibleNextStates = []StateName{GivePoints, StartRound}
	return &gs
}

func (state *GetTrickWinnerState) DoState(game *Game) StateName {

	c1 := *game.PlayedCards[0]
	c2 := *game.PlayedCards[1]
	c3 := *game.PlayedCards[2]
	c4 := *game.PlayedCards[3]

	trump := game.Trump
	lead := game.PlayedCards[0].suite

	winningCard := GetWinningCard(c1, c2, c3, c4, trump, lead)
	winningPlayer := game.Players[game.DealerIndex]

	switch winningCard {
	case c1:
		winningPlayer = game.Players[(game.PlayerIndex+1)%4]
	case c2:
		winningPlayer = game.Players[(game.PlayerIndex+2)%4]
	case c3:
		winningPlayer = game.Players[(game.PlayerIndex+3)%4]
	case c4:
		winningPlayer = game.Players[game.PlayerIndex]
	}

	// print the winner
	game.Log("%s won the trick with a %s", winningPlayer.name, winningCard.ToString())

	// give the winner the trick point
	winningPlayer.tricksTaken += 1

	// return the played cards to the deck
	game.ReturnPlayedCards()

	// if this is the last trick, move on to GivePointsState
	if len(winningPlayer.hand) == 0 {
		return GivePoints
	}

	// next player is the winner
	game.PlayerIndex = winningPlayer.index

	// start new round
	return GetPlayerCard
}

// ============================ GivePointsState ============================
type GivePointsState struct {
	NamedState
}

func NewGivePointsState() *GivePointsState {
	gs := GivePointsState{NamedState{Name: GivePoints}}
	gs.PossibleNextStates = []StateName{CheckForWinner}
	return &gs
}

func (state *GivePointsState) DoState(game *Game) StateName {
	player1 := game.Players[0]
	player2 := game.Players[1]
	player3 := game.Players[2]
	player4 := game.Players[3]

	teamOneTricks := player1.tricksTaken + player3.tricksTaken
	teamTwoTricks := player2.tricksTaken + player4.tricksTaken

	teamOneOrdered := game.OrderedPlayerIndex == 0 || game.OrderedPlayerIndex == 2

	if teamOneOrdered {
		if teamOneTricks == 5 {
			// team one gets 2 points
			player1.pointsEarned += 2
			player3.pointsEarned += 2
			game.Log("Team one won them all! They earned 2 points.")
		} else if teamOneTricks >= 3 {
			// team one gets 1 point
			player1.pointsEarned += 1
			player3.pointsEarned += 1
			game.Log("Team one won %d tricks. They earned 1 point.", teamOneTricks)
		} else {
			// team two gets 2 points
			player2.pointsEarned += 2
			player4.pointsEarned += 2
			game.Log("Team One got euchred! Team two earned 2 points.")
		}
	} else {
		if teamTwoTricks == 5 {
			// team two gets 2 points
			player2.pointsEarned += 2
			player4.pointsEarned += 2
			game.Log("Team two won them all! They earned 2 points.")
		} else if teamTwoTricks >= 3 {
			// team two gets 1 point
			player2.pointsEarned += 1
			player4.pointsEarned += 1
			game.Log("Team two won %d tricks. They earned 1 point.", teamTwoTricks)
		} else {
			// team one gets 2 points
			player1.pointsEarned += 2
			player3.pointsEarned += 2
			game.Log("Team two got euchred! Team one earned 2 points.")
		}
	}

	// reset trick count
	player1.tricksTaken = 0
	player2.tricksTaken = 0
	player3.tricksTaken = 0
	player4.tricksTaken = 0

	// check for winner
	return CheckForWinner
}

// ============================ CheckForWinnerState ============================
type CheckForWinnerState struct {
	NamedState
}

func NewCheckForWinnerState() *CheckForWinnerState {
	gs := CheckForWinnerState{NamedState{Name: CheckForWinner}}
	gs.PossibleNextStates = []StateName{ResetDeckAndShuffle, EndGame}
	return &gs
}

func (state *CheckForWinnerState) DoState(game *Game) StateName {
	player1 := game.Players[0]
	player2 := game.Players[1]

	teamOnePoints := player1.pointsEarned
	teamTwoPoints := player2.pointsEarned

	game.Log("Team One Points: %d", teamOnePoints)
	game.Log("Team Two Points: %d", teamTwoPoints)

	if teamOnePoints >= 4 {
		game.Log("Team One wins!")
		return EndGame
	} else if teamTwoPoints >= 4 {
		game.Log("Team Two wins!")
		return EndGame
	}

	// increment the dealer
	game.DealerIndex = (game.DealerIndex + 1) % 4
	game.PlayerIndex = (game.DealerIndex + 1) % 4
	game.Log("Dealer is now %s", game.Players[game.DealerIndex].name)

	game.OrderedPlayerIndex = -1
	return ResetDeckAndShuffle
}

// ============================ GameOverState ============================
type EndGameState struct {
	NamedState
}

func NewEndGameState() *EndGameState {
	return &EndGameState{NamedState{Name: EndGame}}
}

func (state *EndGameState) DoState(game *Game) StateName {
	game.Log("Game Over!")
	return EndGame
}
