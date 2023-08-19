package game

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
	//fmt.Printf("\n<-- Entered %s State -->\n", state.Name)
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
	game.Log("%s was drawn", game.PlayedCards[lastIndex].ToString())

	if game.PlayedCards[lastIndex].rank == JACK {
		// got trump. Set dealer and continue
		game.DealerIndex = game.PlayerIndex
		game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
		dealer := game.Players[game.DealerIndex]
		game.Log("%s is dealer", dealer.name)

		game.TransitionState(NewResetDeckAndShuffleState())
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
		game.TransitionState(NewRevealTopCardState())
	}
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
	game.Log("%s was turned", game.TurnedCard.ToString())
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

	player := game.Players[game.PlayerIndex]

	// ask player if they want trump
	pickedUp := GetTrumpSelectionOneInput(player, *game.TurnedCard)

	// if picked up, we want to ask the dealer if they want the turned card
	if pickedUp {
		game.Log("%s ordered it up", player.name)
		game.OrderedPlayerIndex = game.PlayerIndex
		game.Trump = game.TurnedCard.suite
		game.TransitionState(NewDealerPickupTrumpState())
		return
	}

	// if this player was the dealer, we will move on to Trump Selection Two
	if game.PlayerIndex == game.DealerIndex {
		game.NextPlayer()
		game.TransitionState(NewTrumpSelectionTwoState())
		return
	}

	// otherwise, move on to the next player
	game.NextPlayer()
}

// ============================ DealerPickupTrumpState ============================
type DealerPickupTrumpState struct {
	NamedState
}

func NewDealerPickupTrumpState() *DealerPickupTrumpState {
	return &DealerPickupTrumpState{NamedState{Name: DealerPickupTrump}}
}

func (state *DealerPickupTrumpState) DoState(game *Game) {
	dealer := game.Players[game.DealerIndex]
	// give the dealer the turned card and let them exchange
	dealer.GiveCard(game.TurnedCard)
	burnCard := GetDealersBurnCard(dealer)
	dealer.ReturnCard(burnCard)
	game.Deck.ReturnCard(burnCard)
	game.TurnedCard = nil
	game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
	game.TransitionState(NewStartRoundState())
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
		game.Log("Dealer got screwed!")
		game.TransitionState(NewScrewDealerState())
		return
	}

	// otherwise, let the next player pick a suite if they want
	selectedSuite := GetTrumpSelectionTwoInput(player, *game.TurnedCard)

	// if the player selected a suite, set it as trump
	if selectedSuite != NONE {
		game.Trump = selectedSuite
		game.PlayerIndex = (game.DealerIndex + 1) % 4 // first player is next to dealer
		game.Deck.ReturnCard(game.TurnedCard)
		game.TurnedCard = nil
		game.TransitionState(NewStartRoundState())
		game.Log("%s picked %s as trump", player.name, selectedSuite.ToString())
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

	game.Log("Dealer %s picked %s as trump", player.name, selectedSuite.ToString())

	game.Trump = selectedSuite
	game.Deck.ReturnCard(game.TurnedCard)
	game.TurnedCard = nil
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
	game.PlayerIndex = (game.DealerIndex + 1) % 4
	game.TransitionState(NewGetPlayerCardState())
}

// ============================ GetPlayerCardState ============================
type GetPlayerCardState struct {
	NamedState
}

func NewGetPlayerCardState() *GetPlayerCardState {
	return &GetPlayerCardState{NamedState{Name: GetPlayerCard}}
}

func (state *GetPlayerCardState) DoState(game *Game) {
	player := game.Players[game.PlayerIndex]

	player.playedCard = GetCardInput(player)
	game.TransitionState(NewCheckValidCardState())
}

// ============================ CheckValidCardState ============================
type CheckValidCardState struct {
	NamedState
}

func NewCheckValidCardState() *CheckValidCardState {
	return &CheckValidCardState{NamedState{Name: CheckValidCard}}
}

func (state *CheckValidCardState) DoState(game *Game) {
	player := game.Players[game.PlayerIndex]

	var leadCard *Card = nil

	if len(game.PlayedCards) > 0 {
		leadCard = game.PlayedCards[0]
	}

	// if the card wasn't valid, go back to GetPlayerCardState
	if !IsCardPlayable(player.playedCard, player.hand, game.Trump, leadCard) {
		game.Log("Invalid card. You must follow suite.")
		player.playedCard = nil
		game.TransitionState(NewGetPlayerCardState())
		return
	}

	// if card is valid, move on to play it
	game.TransitionState(NewPlayCardState())
}

// ============================ PlayCardState ============================
type PlayCardState struct {
	NamedState
}

func NewPlayCardState() *PlayCardState {
	return &PlayCardState{NamedState{Name: PlayCard}}
}

func (state *PlayCardState) DoState(game *Game) {
	player := game.Players[game.PlayerIndex]

	// remove the card from the players hand
	player.ReturnCard(player.playedCard)

	// add the card to the played cards
	game.PlayCard(player.playedCard)

	// print the card
	game.Log("%s played %s", player.name, player.playedCard.ToString())

	// if this is the last card, move on to GetTrickWinnerState
	if len(game.PlayedCards) == 4 {
		game.TransitionState(NewGetTrickWinnerState())
		return
	}

	// move on to the next player
	game.NextPlayer()
	game.TransitionState(NewGetPlayerCardState())
}

// ============================ GetTrickWinnerState ============================
type GetTrickWinnerState struct {
	NamedState
}

func NewGetTrickWinnerState() *GetTrickWinnerState {
	return &GetTrickWinnerState{NamedState{Name: GetTrickWinner}}
}

func (state *GetTrickWinnerState) DoState(game *Game) {

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
		game.TransitionState(NewGivePointsState())
		return
	}

	// next player is the winner
	game.PlayerIndex = winningPlayer.index

	// start new round
	game.TransitionState(NewGetPlayerCardState())
}

// ============================ GivePointsState ============================
type GivePointsState struct {
	NamedState
}

func NewGivePointsState() *GivePointsState {
	return &GivePointsState{NamedState{Name: GivePoints}}
}

func (state *GivePointsState) DoState(game *Game) {
	player1 := game.Players[0]
	player2 := game.Players[1]
	player3 := game.Players[2]
	player4 := game.Players[3]

	teamOneTricks := player1.pointsEarned + player3.pointsEarned
	teamTwoTricks := player2.pointsEarned + player4.pointsEarned

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
	game.TransitionState(NewCheckForWinnerState())
}

// ============================ CheckForWinnerState ============================
type CheckForWinnerState struct {
	NamedState
}

func NewCheckForWinnerState() *CheckForWinnerState {
	return &CheckForWinnerState{NamedState{Name: CheckForWinner}}
}

func (state *CheckForWinnerState) DoState(game *Game) {
	player1 := game.Players[0]
	player2 := game.Players[1]

	teamOnePoints := player1.pointsEarned
	teamTwoPoints := player2.pointsEarned

	game.Log("Team One Points: %d", teamOnePoints)
	game.Log("Team Two Points: %d", teamTwoPoints)

	if teamOnePoints >= 4 {
		game.Log("Team One wins!")
		game.TransitionState(NewGameOverState())
	} else if teamTwoPoints >= 4 {
		game.Log("Team Two wins!")
		game.TransitionState(NewGameOverState())
	} else {
		game.TransitionState(NewResetDeckAndShuffleState())
	}
}

// ============================ GameOverState ============================
type GameOverState struct {
	NamedState
}

func NewGameOverState() *GameOverState {
	return &GameOverState{NamedState{Name: EndGame}}
}

func (state *GameOverState) DoState(game *Game) {
	game.Log("Game Over!")
}
