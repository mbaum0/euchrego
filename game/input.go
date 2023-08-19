package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func isValidSuite(invalidSuite Suite, input string) bool {
	switch input {
	case "h":
		return invalidSuite != HEART

	case "d":
		return invalidSuite != DIAMOND

	case "c":
		return invalidSuite != CLUB

	case "s":
		return invalidSuite != SPADE
	default:
		return false
	}
}

func GetTrumpSelectionOneInput(player *Player, card Card) bool {
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
// not be that of the turned up card.
func GetTrumpSelectionTwoInput(player *Player, card Card) Suite {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s: Do you want to pick a suite? (y/n): ", player.name))
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if input == "y" {
			return GetSuiteInput(player, card.suite)
		} else if input == "n" {
			return NONE
		}
		showInvalid = true
	}
}

// GetScrewTheDealerInput is the same as GetTrumpSelectionTwoInput, expect they must select a suite.
func GetScrewTheDealerInput(player *Player, turnedCard Card) Suite {
	var builder strings.Builder

	// write a prompt string that doesn't include the invalid suite
	switch turnedCard.suite {
	case HEART:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (d/c/s): ", player.name))
	case DIAMOND:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/c/s): ", player.name))
	case CLUB:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/s): ", player.name))
	case SPADE:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c): ", player.name))
	}
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if isValidSuite(turnedCard.suite, input) {
			switch input {
			case "h":
				return HEART

			case "d":
				return DIAMOND

			case "c":
				return CLUB

			case "s":
				return SPADE
			}
		}
		showInvalid = true
	}
}

// GetDealersBurnCard prompts the dealer to select a card to discard. The input
// will be the index of the card in their hand
func GetDealersBurnCard(dealer *Player) *Card {
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
func GetSuiteInput(player *Player, invalidSuite Suite) Suite {
	var builder strings.Builder

	// write a prompt string that doesn't include the invalid suite
	switch invalidSuite {
	case HEART:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (d/c/s): ", player.name))
	case DIAMOND:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/c/s): ", player.name))
	case CLUB:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/s): ", player.name))
	case SPADE:
		builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c): ", player.name))
	}
	prompt := builder.String()
	showInvalid := false
	for {
		input := promptUser(prompt, showInvalid)
		if isValidSuite(invalidSuite, input) {
			switch input {
			case "h":
				return HEART

			case "d":
				return DIAMOND

			case "c":
				return CLUB

			case "s":
				return SPADE
			}
		}
		showInvalid = true
	}
}

// Prompt the player to select a card from their hand. The input will be the index of the card in their hand
func GetCardInput(player *Player) *Card {
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
