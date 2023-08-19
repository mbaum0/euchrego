package game

import (
	"fmt"

	"github.com/fatih/color"
)

type TextDisplay struct {
	width  int
	height int
	grid   [][]string
}

func NewTextDisplay(width, height int) *TextDisplay {
	t := TextDisplay{}
	t.width = width
	t.height = height
	t.grid = make([][]string, height)
	for i := 0; i < height; i++ {
		t.grid[i] = make([]string, width)
	}
	return &t
}

func (t *TextDisplay) Render() {
	moveCursorHome()
	for _, row := range t.grid {
		// combine the row into a string
		rowString := ""
		for _, cell := range row {
			rowString += cell
		}
		fmt.Print(rowString)

		fmt.Print("\n")
	}
}

func (t *TextDisplay) ClearDisplay() {
	for i := 0; i < t.height; i++ {
		for j := 0; j < t.width; j++ {
			t.grid[i][j] = " "
		}
	}
}

func (t *TextDisplay) DrawVerticalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.grid[y+i][x] = "│"
	}
}

func (t *TextDisplay) DrawHorizontalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.grid[y][x+i] = "-"
	}
}

func (t *TextDisplay) DrawCard(x, y int, card Card) {
	cardArt := getCardArt(card)

	colorWay := color.New(color.FgRed).SprintFunc()
	if card.suite == SPADE || card.suite == CLUB {
		colorWay = color.New(color.FgGreen).SprintFunc()
	}
	for i, row := range cardArt {
		for j, cell := range row {
			t.grid[y+i][x+j] = colorWay(string(cell))
		}
	}
}

func getCardArt(c Card) [][]rune {
	cardRows := 9
	cardCols := 11
	cardArt := make([][]rune, cardRows)
	for i := 0; i < cardRows; i++ {
		cardArt[i] = make([]rune, cardCols)
	}

	var suitSymbol string
	switch c.suite {
	case HEART:
		suitSymbol = "♥ ♥ ♥"
	case DIAMOND:
		suitSymbol = "♦ ♦ ♦"
	case CLUB:
		suitSymbol = "♣ ♣ ♣"
	case SPADE:
		suitSymbol = "♠ ♠ ♠"
	case NONE:
		suitSymbol = "   "
	}

	rank := c.GetRank()
	rankChar := rank.ToChar()

	cardArt[0] = []rune("┌─────────┐")
	cardArt[1] = []rune(fmt.Sprintf("│  %s      │", rankChar))
	cardArt[2] = []rune("│         │")
	cardArt[3] = []rune("│         │")
	cardArt[4] = []rune(fmt.Sprintf("│  %s  │", suitSymbol))
	cardArt[5] = []rune("│         │")
	cardArt[6] = []rune("│         │")
	cardArt[7] = []rune(fmt.Sprintf("│       %s │", rankChar))
	cardArt[8] = []rune("└─────────┘")

	// 10 is a special case because it has two characters
	if rank == TEN {
		cardArt[1] = []rune("│  10     │")
		cardArt[7] = []rune("│     10  │")
	}

	return cardArt
}

func (t *TextDisplay) DrawText(x, y int, text string) {
	for i, c := range text {
		t.grid[y][x+i] = string(c)
	}
}

func (t *TextDisplay) DrawPlayerHand(x, y int, player Player, enumerate bool) {
	cards := player.hand
	for i, card := range cards {
		t.DrawCard(x+i*12, y, *card)
	}

	// draw the index of the card beneath each card
	if enumerate {
		for i := range cards {
			t.DrawText(x+4+i*12, y+9, fmt.Sprintf("(%d)", i))
		}
	}
}

func (t *TextDisplay) DrawPlayerHands(game *Game) {
	player1 := *game.Players[0]
	player2 := *game.Players[1]
	player3 := *game.Players[2]
	player4 := *game.Players[3]

	t.DrawText(1, 2, player1.name)
	t.DrawPlayerHand(0, 3, player1, true)

	t.DrawText(1, 14, player2.name)
	t.DrawPlayerHand(0, 15, player2, true)

	t.DrawText(1, 26, player3.name)
	t.DrawPlayerHand(0, 27, player3, true)

	t.DrawText(1, 38, player4.name)
	t.DrawPlayerHand(0, 39, player4, true)
}

func (t *TextDisplay) DrawDealerArrow(game *Game) {
	dealerIndex := game.DealerIndex

	y := 5 + 12*dealerIndex
	x := 60
	t.DrawText(x, y, "<-- Dealer")
}

func (t *TextDisplay) DrawTurnArrow(game *Game) {
	playerIndex := game.PlayerIndex

	y := 6 + 12*playerIndex
	x := 60
	t.DrawText(x, y, "<-- Turn")
}

func (t *TextDisplay) DrawPlayedCards(game *Game) {
	t.DrawText(80, 2, "Played Cards")
	cards := game.PlayedCards
	if len(cards) > 4 {
		return
	}

	for i, card := range cards {
		t.DrawCard(80, 5+(10*i), *card)
	}
}

func (t *TextDisplay) DrawTurnedCard(game *Game) {
	t.DrawText(100, 2, "Turned Card")
	card := game.TurnedCard
	if card == nil {
		return
	}
	t.DrawCard(100, 5, *card)
}

func (t *TextDisplay) DrawLogs(game *Game) {
	for i, log := range game.logs {
		t.DrawText(120, 2+i, log)
	}
}

func (t *TextDisplay) DrawStats(game *Game) {
	t.DrawText(170, 2, "Stats")
	t.DrawText(170, 3, "-----")
	t.DrawText(170, 4, fmt.Sprintf("Trump:         \t%s", game.Trump.ToString()))
	t.DrawText(170, 5, fmt.Sprintf("Ordered Up:    \t%s", game.Players[game.OrderedPlayerIndex].name))
	t.DrawText(170, 6, fmt.Sprintf("Dealer:        \t%s", game.Players[game.DealerIndex].name))
	t.DrawText(170, 7, fmt.Sprintf("Turn:          \t%s", game.Players[game.PlayerIndex].name))
	turnedCardString := ""
	if game.TurnedCard != nil {
		turnedCardString = game.TurnedCard.ToString()
	}
	t.DrawText(170, 8, fmt.Sprintf("Turned Card:   \t%s", turnedCardString))
	t.DrawText(170, 9, fmt.Sprintf("Played Cards:  \t%d", len(game.PlayedCards)))
	t.DrawText(170, 10, fmt.Sprintf("State:        \t%s", game.State.GetName()))
	t.DrawText(170, 11, fmt.Sprintf("Cards in Deck:\t%d", len(game.Deck.cards)))
	t.DrawText(170, 12, fmt.Sprintf("Team 1 Tricks:\t%d", game.Players[0].tricksTaken+game.Players[2].tricksTaken))
	t.DrawText(170, 13, fmt.Sprintf("Team 2 Tricks:\t%d", game.Players[1].tricksTaken+game.Players[3].tricksTaken))
	t.DrawText(170, 14, fmt.Sprintf("Team 1 Points:\t%d", game.Players[0].pointsEarned))
	t.DrawText(170, 15, fmt.Sprintf("Team 2 Points:\t%d", game.Players[1].pointsEarned))
}

func (t *TextDisplay) DrawBoard(game *Game) {
	t.ClearDisplay()
	if game.State.GetName() == InitGame {
		return
	}
	t.DrawPlayerHands(game)
	t.DrawDealerArrow(game)
	t.DrawTurnArrow(game)
	t.DrawVerticalLine(75, 0, 50)
	t.DrawPlayedCards(game)
	t.DrawVerticalLine(95, 0, 50)
	t.DrawTurnedCard(game)
	t.DrawVerticalLine(115, 0, 50)
	t.DrawLogs(game)
	t.DrawVerticalLine(165, 0, 50)
	t.DrawStats(game)
	t.Render()
}

func moveCursorHome() {
	fmt.Print("\033[H")
}
