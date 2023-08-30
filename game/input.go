package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbaum0/euchrego/godeck"
)

type InputDevice struct {
	requestChan chan string
	inputChan   chan string
}

func NewInputDevice(requestChan chan string, inputChan chan string) *InputDevice {
	return &InputDevice{requestChan, inputChan}
}

func (i *InputDevice) pollForInput(prompt string) string {
	i.requestChan <- prompt
	input := <-i.inputChan
	return strings.TrimSpace(strings.ToLower(input))
}

func (i *InputDevice) isValidSuite(invalidSuite godeck.Suit, input string) bool {
	switch input {
	case "h":
		return invalidSuite != godeck.Hearts

	case "d":
		return invalidSuite != godeck.Diamonds

	case "c":
		return invalidSuite != godeck.Clubs

	case "s":
		return invalidSuite != godeck.Spades
	default:
		return false
	}
}

func (i *InputDevice) GetTrumpSelectionOneInput(player *Player, card godeck.Card) bool {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Order it up or pass? (o/p): ", player.name))
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		if input == "o" {
			return true
		} else if input == "p" {
			return false
		}
	}
}

// GetTrumpSelectionTwoInput asks the player if they want to select a suite for trump. The suite can
// not be that of the turned up godeck.
func (i *InputDevice) GetTrumpSelectionTwoInput(player *Player, c godeck.Card) godeck.Suit {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Do you want to pick a suite? (y/n): ", player.name))
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		if input == "y" {
			return i.GetSuiteInput(player, c.Suit())
		} else if input == "n" {
			return godeck.None
		}
	}
}

// GetScrewTheDealerInput is the same as GetTrumpSelectionTwoInput, expect they must select a suite.
func (i *InputDevice) GetScrewTheDealerInput(player *Player, turnedCard godeck.Card) godeck.Suit {
	var builder strings.Builder

	// write a prompt string that doesn't include the invalid suite
	switch turnedCard.Suit() {
	case godeck.Hearts:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (d/c/s): ", player.name))
	case godeck.Diamonds:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/c/s): ", player.name))
	case godeck.Clubs:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/s): ", player.name))
	case godeck.Spades:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c): ", player.name))
	}
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		if i.isValidSuite(turnedCard.Suit(), input) {
			switch input {
			case "h":
				return godeck.Hearts

			case "d":
				return godeck.Diamonds

			case "c":
				return godeck.Clubs

			case "s":
				return godeck.Spades
			}
		}
	}
}

// GetDealersBurnCard prompts the dealer to select a card to disgodeck. The input
// will be the index of the card in their hand
func (i *InputDevice) GetDealersBurnCard(dealer *Player) godeck.Card {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Pick a card to discard: ", dealer.name))
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		index, err := strconv.Atoi(input)
		if err != nil {
			continue
		}

		if index < 0 || index >= len(dealer.hand) {
			continue
		}

		return dealer.hand[index]

	}
}

// GetSuiteInput prompts the player to select a suite that isn't the invalidSuite
func (i *InputDevice) GetSuiteInput(player *Player, invalidSuite godeck.Suit) godeck.Suit {
	var builder strings.Builder

	// write a prompt string that doesn't include the invalid suite
	switch invalidSuite {
	case godeck.Hearts:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (d/c/s): ", player.name))
	case godeck.Diamonds:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/c/s): ", player.name))
	case godeck.Clubs:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/s): ", player.name))
	case godeck.Spades:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c): ", player.name))
	}
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		if i.isValidSuite(invalidSuite, input) {
			switch input {
			case "h":
				return godeck.Hearts

			case "d":
				return godeck.Diamonds

			case "c":
				return godeck.Clubs

			case "s":
				return godeck.Spades
			}
		}
	}
}

// Prompt the player to select a card from their hand. The input will be the index of the card in their hand
func (i *InputDevice) GetCardInput(player *Player) godeck.Card {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Pick a card: ", player.name))
	prompt := builder.String()
	for {
		input := i.pollForInput(prompt)
		index, err := strconv.Atoi(input)
		if err != nil {
			continue
		}

		if index < 0 || index >= len(player.hand) {
			continue
		}

		return player.hand[index]

	}
}
