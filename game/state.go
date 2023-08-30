package game

import (
	"github.com/mbaum0/euchrego/fsm"
	"github.com/mbaum0/euchrego/godeck"
)

type GameMachine struct {
	*GameBoard
	RequestInput chan string
	Input        chan string
}

func (gm *GameMachine) InitGameState() (fsm.StateFunc, error) {
	gm.Deck = godeck.NewEuchreDeck(godeck.RandomShuffleSeed())
	gm.Deck.Shuffle()
	return gm.DrawForDealerState, nil
}

func (gm *GameMachine) DrawForDealerState() (fsm.StateFunc, error) {
	c, _ := gm.Deck.DrawCard()
	gm.PlayedCards = append(gm.PlayedCards, c)
	lastIndex := len(gm.PlayedCards) - 1

	// print drawn card
	gm.Log("%s was drawn", gm.PlayedCards[lastIndex])

	// if the card is a jack, the player who drew it is the dealer
	if c.Rank() == godeck.Jack {
		gm.DealerIndex = gm.PlayerIndex

		// first player is the player to the left of the dealer
		gm.PlayerIndex = (gm.DealerIndex + 1) % 4

		dealer := gm.Players[gm.DealerIndex]
		gm.Log("%s is the dealer", dealer.GetName())
		gm.ReturnPlayedCards()
		return gm.ResetDeckAndShuffleState, nil
	}

	// otherwise, the next player draws
	gm.NextPlayer()
	return gm.DrawForDealerState, nil
}

func (gm *GameMachine) ResetDeckAndShuffleState() (fsm.StateFunc, error) {
	gm.Deck.Shuffle()
	return gm.DealCardsState, nil
}

func (gm *GameMachine) DealCardsState() (fsm.StateFunc, error) {
	// deal cards in standard euchre fashion
	dealerIndex := gm.DealerIndex
	dealer := gm.Players[dealerIndex]

	playerIndex := gm.PlayerIndex
	player := gm.Players[playerIndex]

	isFirstDeal := len(dealer.hand) == 0

	if isFirstDeal {
		// deal 2 cards to the player if they are the 1st or 3rd player
		if playerIndex == (dealerIndex+1)%4 || playerIndex == (dealerIndex+3)%4 {
			drawnCards, _ := gm.Deck.DrawCards(2)
			player.GiveCards(drawnCards)
			gm.Log("%s was dealt 2 cards", player.name)
		} else {
			drawnCards, _ := gm.Deck.DrawCards(3)
			player.GiveCards(drawnCards)
			gm.Log("%s was dealt 3 cards", player.name)
		}
	} else {
		// deal 3 cards to the player if they are the 1st or 3rd player
		if playerIndex == (dealerIndex+1)%4 || playerIndex == (dealerIndex+3)%4 {
			drawnCards, _ := gm.Deck.DrawCards(3)
			player.GiveCards(drawnCards)
			gm.Log("%s was dealt 3 cards", player.name)
		} else {
			drawnCards, _ := gm.Deck.DrawCards(2)
			player.GiveCards(drawnCards)
			gm.Log("%s was dealt 2 cards", player.name)
		}
	}

	// move onto next player
	gm.NextPlayer()

	// if the dealer has all their cards, continue to RevealTopCardState
	if len(dealer.hand) == 5 {
		return gm.RevealTopCardState, nil
	}
	return gm.DealCardsState, nil
}

func (gm *GameMachine) RevealTopCardState() (fsm.StateFunc, error) {
	c, _ := gm.Deck.DrawCard()
	gm.TurnedCard = c

	// print out name of turned card
	gm.Log("%s was turned", gm.TurnedCard)
	return gm.TrumpSelectionOneState, nil
}

func (gm *GameMachine) TrumpSelectionOneState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	// ask player if they want trump
	pickedUp := GetTrumpSelectionOneInput(player, gm.TurnedCard)

	// if picked up, we want to ask the dealer if they want the turned card
	if pickedUp {
		gm.Log("%s ordered it up", player.name)
		gm.OrderedPlayerIndex = gm.PlayerIndex
		gm.Trump = gm.TurnedCard.Suit()
		return gm.DealerPickupTrumpState, nil
	}

	// if this player was the dealer, we will move on to Trump Selection Two
	if gm.PlayerIndex == gm.DealerIndex {
		gm.NextPlayer()
		return gm.TrumpSelectionTwoState, nil
	}

	// otherwise, move on to the next player
	gm.NextPlayer()
	return gm.TrumpSelectionOneState, nil
}

func (gm *GameMachine) DealerPickupTrumpState() (fsm.StateFunc, error) {
	dealer := gm.Players[gm.DealerIndex]
	// give the dealer the turned card and let them exchange
	dealer.GiveCard(gm.TurnedCard)
	burnCard := GetDealersBurnCard(dealer)
	dealer.ReturnCard(burnCard)
	gm.Deck.ReturnCard(burnCard)
	gm.TurnedCard = godeck.EmptyCard()
	gm.PlayerIndex = (gm.DealerIndex + 1) % 4 // first player is next to dealer
	return gm.StartRoundState, nil
}

func (gm *GameMachine) TrumpSelectionTwoState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	// if the player is the dealer, they must select a suite
	if gm.PlayerIndex == gm.DealerIndex {
		gm.Log("Dealer got screwed!")
		return gm.ScrewDealerState, nil
	}

	// otherwise, let the next player pick a suite if they want
	selectedSuite := GetTrumpSelectionTwoInput(player, gm.TurnedCard)

	// if the player selected a suite, set it as trump
	if selectedSuite != godeck.None {
		gm.Trump = selectedSuite
		gm.PlayerIndex = (gm.DealerIndex + 1) % 4 // first player is next to dealer
		gm.Deck.ReturnCard(gm.TurnedCard)
		gm.TurnedCard = godeck.EmptyCard()
		gm.Log("%s picked %s as trump", player.name, selectedSuite)
		return gm.StartRoundState, nil
	}

	// move on to the next player
	gm.NextPlayer()
	return gm.TrumpSelectionTwoState, nil
}

func (gm *GameMachine) ScrewDealerState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	selectedSuite := GetScrewTheDealerInput(player, gm.TurnedCard)

	gm.Log("Dealer %s picked %s as trump", player.name, selectedSuite)

	gm.Trump = selectedSuite
	gm.Deck.ReturnCard(gm.TurnedCard)
	gm.TurnedCard = godeck.EmptyCard()
	return gm.StartRoundState, nil
}

func (gm *GameMachine) StartRoundState() (fsm.StateFunc, error) {
	gm.PlayerIndex = (gm.DealerIndex + 1) % 4
	return gm.GetPlayerCardState, nil
}

func (gm *GameMachine) GetPlayerCardState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	player.playedCard = GetCardInput(player)
	return gm.CheckValidCardState, nil
}

func (gm *GameMachine) CheckValidCardState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	var leadCard godeck.Card

	if len(gm.PlayedCards) > 0 {
		leadCard = gm.PlayedCards[0]
	}

	// if the card wasn't valid, go back to GetPlayerCardState
	if gm.Deck.IsCardPlayable(player.playedCard, player.hand, gm.Trump, leadCard) {
		gm.Log("Invalid card. You must follow suite.")
		player.playedCard = godeck.EmptyCard()
		return gm.GetPlayerCardState, nil
	}

	// if card is valid, move on to play it
	return gm.PlayCardState, nil
}

func (gm *GameMachine) PlayCardState() (fsm.StateFunc, error) {
	player := gm.Players[gm.PlayerIndex]

	// remove the card from the players hand
	player.ReturnCard(player.playedCard)

	// add the card to the played cards
	gm.PlayCard(player.playedCard)

	// print the card
	gm.Log("%s played %s", player.name, player.playedCard)

	// if this is the last card, move on to GetTrickWinnerState
	if len(gm.PlayedCards) == 4 {
		return gm.GetTrickWinnerState, nil
	}

	// move on to the next player
	gm.NextPlayer()
	return gm.GetPlayerCardState, nil
}

func (gm *GameMachine) GetTrickWinnerState() (fsm.StateFunc, error) {
	c1 := gm.PlayedCards[0]
	c2 := gm.PlayedCards[1]
	c3 := gm.PlayedCards[2]
	c4 := gm.PlayedCards[3]

	trump := gm.Trump
	lead := gm.PlayedCards[0].Suit()

	winningCard := gm.Deck.GetWinningCard(c1, c2, c3, c4, trump, lead)
	winningPlayer := gm.Players[gm.DealerIndex]

	switch winningCard {
	case c1:
		winningPlayer = gm.Players[(gm.PlayerIndex+1)%4]
	case c2:
		winningPlayer = gm.Players[(gm.PlayerIndex+2)%4]
	case c3:
		winningPlayer = gm.Players[(gm.PlayerIndex+3)%4]
	case c4:
		winningPlayer = gm.Players[gm.PlayerIndex]
	}

	// print the winner
	gm.Log("%s won the trick with a %s", winningPlayer.name, winningCard)

	// give the winner the trick point
	winningPlayer.tricksTaken += 1

	// return the played cards to the deck
	gm.ReturnPlayedCards()

	// if this is the last trick, move on to GivePointsState
	if len(winningPlayer.hand) == 0 {
		return gm.GivePointsState, nil
	}

	// next player is the winner
	gm.PlayerIndex = winningPlayer.index

	// get card for next trick
	return gm.GetPlayerCardState, nil
}

func (gm *GameMachine) GivePointsState() (fsm.StateFunc, error) {
	player1 := gm.Players[0]
	player2 := gm.Players[1]
	player3 := gm.Players[2]
	player4 := gm.Players[3]

	teamOneTricks := player1.tricksTaken + player3.tricksTaken
	teamTwoTricks := player2.tricksTaken + player4.tricksTaken

	teamOneOrdered := gm.OrderedPlayerIndex == 0 || gm.OrderedPlayerIndex == 2

	if teamOneOrdered {
		if teamOneTricks == 5 {
			// team one gets 2 points
			player1.pointsEarned += 2
			player3.pointsEarned += 2
			gm.Log("Team one won them all! They earned 2 points.")
		} else if teamOneTricks >= 3 {
			// team one gets 1 point
			player1.pointsEarned += 1
			player3.pointsEarned += 1
			gm.Log("Team one won %d tricks. They earned 1 point.", teamOneTricks)
		} else {
			// team two gets 2 points
			player2.pointsEarned += 2
			player4.pointsEarned += 2
			gm.Log("Team One got euchred! Team two earned 2 points.")
		}
	} else {
		if teamTwoTricks == 5 {
			// team two gets 2 points
			player2.pointsEarned += 2
			player4.pointsEarned += 2
			gm.Log("Team two won them all! They earned 2 points.")
		} else if teamTwoTricks >= 3 {
			// team two gets 1 point
			player2.pointsEarned += 1
			player4.pointsEarned += 1
			gm.Log("Team two won %d tricks. They earned 1 point.", teamTwoTricks)
		} else {
			// team one gets 2 points
			player1.pointsEarned += 2
			player3.pointsEarned += 2
			gm.Log("Team two got euchred! Team one earned 2 points.")
		}
	}

	// reset trick count
	player1.tricksTaken = 0
	player2.tricksTaken = 0
	player3.tricksTaken = 0
	player4.tricksTaken = 0

	// check for winner
	return gm.CheckForWinnerState, nil
}

func (gm *GameMachine) CheckForWinnerState() (fsm.StateFunc, error) {
	player1 := gm.Players[0]
	player2 := gm.Players[1]

	teamOnePoints := player1.pointsEarned
	teamTwoPoints := player2.pointsEarned

	gm.Log("Team One Points: %d", teamOnePoints)
	gm.Log("Team Two Points: %d", teamTwoPoints)

	if teamOnePoints >= 4 {
		gm.Log("Team One wins!")
		return gm.EndGameState, nil
	} else if teamTwoPoints >= 4 {
		gm.Log("Team Two wins!")
		return gm.EndGameState, nil
	}

	// increment the dealer
	gm.DealerIndex = (gm.DealerIndex + 1) % 4
	gm.PlayerIndex = (gm.DealerIndex + 1) % 4
	gm.Log("Dealer is now %s", gm.Players[gm.DealerIndex].name)

	gm.OrderedPlayerIndex = -1
	return gm.ResetDeckAndShuffleState, nil
}

func (gm *GameMachine) EndGameState() (fsm.StateFunc, error) {
	gm.Log("Game Over!")
	return nil, nil
}
