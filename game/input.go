package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Order it up or pass? (o/p): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "o" {
			return true
		} else if input == "p" {
			return false
		}
		fmt.Println("Invalid input. ")
	}
}

// GetTrumpSelectionTwoInput asks the player if they want to select a suite for trump. The suite can
// not be that of the turned up card.
func GetTrumpSelectionTwoInput(player *Player, card Card) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Do you want to pick a suite? (y/n): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" {
			return GetSuiteInput(player, card.suite)
		} else if input == "n" {
			return NONE
		}
		fmt.Println("Invalid input. ")
	}
}

// GetScrewTheDealerInput is the same as GetTrumpSelectionTwoInput, expect they must select a suite.
func GetScrewTheDealerInput(player *Player, turnedCard Card) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c/s): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
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
		fmt.Println("Invalid input. ")
	}
}

// GetDealersBurnCard prompts the dealer to select a card to discard. The input
// will be the index of the card in their hand
func GetDealersBurnCard(dealer *Player) *Card {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a card to discard: ", dealer.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. ")
			continue
		}

		if index < 0 || index >= len(dealer.hand) {
			fmt.Println("Invalid input. ")
			continue
		}

		return dealer.hand[index]

	}
}

// GetSuiteInput prompts the player to select a suite that isn't the invalidSuite
func GetSuiteInput(player *Player, invalidSuite Suite) Suite {
	reader := bufio.NewReader(os.Stdin)

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
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
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
		fmt.Println("Invalid input. ")
	}
}

// Prompt the player to select a card from their hand. The input will be the index of the card in their hand
func GetCardInput(player *Player) *Card {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a card: ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. ")
			continue
		}

		if index < 0 || index >= len(player.hand) {
			fmt.Println("Invalid input. ")
			continue
		}

		return player.hand[index]

	}
}
