package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isValidSuite(invalidSuite, input string) bool {
	switch input {
	case "h":
		return invalidSuite != "h"

	case "d":
		return invalidSuite != "d"

	case "c":
		return invalidSuite != "c"

	case "s":
		return invalidSuite != "s"
	default:
		return false
	}
}

func GetTrumpSelectionOneInput(player Player, card Card) bool {
	reader := bufio.NewReader(os.Stdin)

	player.PrintHand()

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Do you want to pick it up or pass? (p/u): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "p" {
			return true
		} else if input == "u" {
			return false
		}
		fmt.Println("Invalid input. ")
	}
}

// GetTrumpSelectionTwoInput asks the player if they want to select a suite for trump. The suite can
// not be that of the turned up card.
func GetTrumpSelectionTwoInput(player Player, card Card) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Do you want to pick a suite? (y/n): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" {
			return GetSuiteInput(&player, card.suite)
		} else if input == "n" {
			return NONE
		}
		fmt.Println("Invalid input. ")
	}
}

// GetScrewTheDealerInput is the same as GetTrumpSelectionTwoInput, expect they must select a suite.
func GetScrewTheDealerInput(player Player, turnedCard Card) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c/s): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if isValidSuite(turnedCard.suite.ToString(), input) {
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

func GetDealerWantsToPickItUp(dealer Player, card Card) bool {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Do you want to exchange it up or pass? (p/u): ", dealer.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "e" {
			return true
		} else if input == "u" {
			return false
		}
		fmt.Println("Invalid input. ")
	}
}

func GetDealersBurnCard(dealer Player) *Card {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a card to discard: ", dealer.name))

	// print out the cards in the dealers hand
	for i, c := range dealer.hand {
		if i == len(dealer.hand)-1 {
			builder.WriteString("or ")
		}
		builder.WriteString(c.ToString())
		if i != len(dealer.hand)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(": ")

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		for _, c := range dealer.hand {
			if c.ToString() == input {
				return c
			}
		}
		fmt.Println("Invalid input. ")
	}
}

// GetSuiteInput prompts the player to select a suite that isn't the invalidSuite
func GetSuiteInput(player *Player, invalidSuite Suite) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s: Pick a suite (h/d/c/s): ", player.name))

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if isValidSuite(invalidSuite.ToString(), input) {
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
