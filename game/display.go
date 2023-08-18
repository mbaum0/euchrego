package game

// func getHandArt(cards []*Card, enumerate bool) string {
// 	if len(cards) == 0 {
// 		return ""
// 	}

// 	var cardArts = make([]string, 0)
// 	for _, c := range cards {
// 		cardArts = append(cardArts, getCardArt(*c))
// 	}

// 	var builder strings.Builder
// 	rows := len(strings.Split(cardArts[0], "\n")) // cards have the same number of rows
// 	for row := 0; row < rows; row++ {
// 		for _, cardArt := range cardArts {
// 			lines := strings.Split(cardArt, "\n")
// 			builder.WriteString(lines[row] + " ")
// 		}
// 		if row != rows-1 {
// 			builder.WriteString("\n")
// 		} else {
// 			builder.WriteString("\r")
// 		}

// 	}

// 	if enumerate {
// 		for i := range cards {
// 			startSpaces := strings.Repeat(" ", 4)
// 			endSpaces := strings.Repeat(" ", 5)
// 			builder.WriteString(fmt.Sprintf("%s(%d)%s", startSpaces, i, endSpaces))
// 		}
// 	}
// 	return builder.String()
// }
