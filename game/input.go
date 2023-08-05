package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func stringToSuite(input string) Suite {
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
	panic("Got unexpected string value")
}

func isValidSuit(suites []Suite, input string) bool {
	var parsedInput Suite
	switch input {
	case "h":
		parsedInput = HEART

	case "d":
		parsedInput = DIAMOND

	case "c":
		parsedInput = CLUB

	case "s":
		parsedInput = SPADE
	default:
		return false
	}
	for _, s := range suites {
		if parsedInput == s {
			return true
		}
	}
	return false
}

func GetSuiteInput(suites ...Suite) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString("Enter a suite: ")
	for _, s := range suites {
		switch s {
		case HEART:
			builder.WriteString("(h)earts, ")
		case CLUB:
			builder.WriteString("(c)lubs, ")
		case DIAMOND:
			builder.WriteString("(d)iamonds, ")
		case SPADE:
			builder.WriteString("(s)pades, ")
		case NONE:
			builder.WriteString("or (p)ass")
		}
	}
	builder.WriteString(": ")

	prompt := builder.String()
	for {

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if isValidSuit(suites, input) {
			return stringToSuite(input)
		}
		fmt.Println("Invalid input. ")
	}
}
