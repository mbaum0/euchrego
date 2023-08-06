package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isValidSuit(allowedSuites []Suite, input string) bool {
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
	case "p":
		parsedInput = NONE
	default:
		return false
	}
	for _, s := range allowedSuites {
		if parsedInput == s {
			return true
		}
	}
	return false
}

func GetSuiteInput(player *Player, suites ...Suite) Suite {
	reader := bufio.NewReader(os.Stdin)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s pick a trump suite: ", player.name))

	for i, s := range suites {
		if i == len(suites)-1 {
			builder.WriteString("or ")
		}

		suiteStr := s.ToString()
		suiteStr = fmt.Sprintf("(%s)%s", string(suiteStr[0]), suiteStr[1:])
		builder.WriteString(suiteStr)
		if i != len(suites)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(": ")

	prompt := builder.String()
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if isValidSuit(suites, input) {
			return SuiteFromChar(input)
		}
		fmt.Println("Invalid input. ")
	}
}
