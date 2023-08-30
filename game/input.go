package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mbaum0/euchrego/godeck"
)

// promptUser should prompt the user for input with the given string.
// the prompt should overwrite the previous prompt and the user's input
func promptUser(prompt string, showInvalid bool) string {
	fmt.Print("\033[1A")  // moves the cursor up 1 line
	fmt.Print("\r\033[K") // erases the current line
	fmt.Print("  > ")
	if showInvalid {
		fmt.Print("Received invalid input! ")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func isValidSuite(invalidSuite godeck.Suit, input string) bool {
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

func GetTrumpSelectionOneInput(player *Player, card godeck.Card) bool {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Order it up or pass? (o/p): ", player.name))
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if input == "o" {
			return true
		} else if input == "p" {
			return false
		}
		showInvalid = true
	}
}

// GetTrumpSelectionTwoInput asks the player if they want to select a suite for trump. The suite can
// not be that of the turned up godeck.
func GetTrumpSelectionTwoInput(player *Player, c godeck.Card) godeck.Suit {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Do you want to pick a suite? (y/n): ", player.name))
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if input == "y" {
			return GetSuiteInput(player, c.Suit())
		} else if input == "n" {
			return godeck.None
		}
		showInvalid = true
	}
}

// GetScrewTheDealerInput is the same as GetTrumpSelectionTwoInput, expect they must select a suite.
func GetScrewTheDealerInput(player *Player, turnedCard godeck.Card) godeck.Suit {
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
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if isValidSuite(turnedCard.Suit(), input) {
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
		showInvalid = true
	}
}

// GetDealersBurnCard prompts the dealer to select a card to disgodeck. The input
// will be the index of the card in their hand
func GetDealersBurnCard(dealer *Player) godeck.Card {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Pick a card to discard: ", dealer.name))
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		index, err := strconv.Atoi(input)
		if err != nil {
			showInvalid = true
			continue
		}

		if index < 0 || index >= len(dealer.hand) {
			showInvalid = true
			continue
		}

		return dealer.hand[index]

	}
}

// GetSuiteInput prompts the player to select a suite that isn't the invalidSuite
func GetSuiteInput(player *Player, invalidSuite godeck.Suit) godeck.Suit {
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
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if isValidSuite(invalidSuite, input) {
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
		showInvalid = true
	}
}

// Prompt the player to select a card from their hand. The input will be the index of the card in their hand
func GetCardInput(player *Player) godeck.Card {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Pick a card: ", player.name))
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		index, err := strconv.Atoi(input)
		if err != nil {
			showInvalid = true
			continue
		}

		if index < 0 || index >= len(player.hand) {
			showInvalid = true
			continue
		}

		return player.hand[index]

	}
}
