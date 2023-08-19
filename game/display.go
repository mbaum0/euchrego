package game

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

const DISPLAY_WIDTH = 166
const DISPLAY_HEIGHT = 51

type TextDisplay struct {
	width  int
	height int
	grid   [][]string
}

func NewTextDisplay() *TextDisplay {
	ClearTerminal()
	t := TextDisplay{}
	t.width = DISPLAY_WIDTH
	t.height = DISPLAY_HEIGHT
	t.grid = make([][]string, DISPLAY_HEIGHT)
	for i := 0; i < DISPLAY_HEIGHT; i++ {
		t.grid[i] = make([]string, DISPLAY_WIDTH)
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

func (t *TextDisplay) DrawRune(x, y int, r rune) {
	t.grid[y][x] = string(r)
}

func (t *TextDisplay) DrawVerticalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.grid[y+i][x] = "│"
	}
}

func (t *TextDisplay) DrawHorizontalLine(x, y, length int) {
	for i := 0; i < length; i++ {
		t.grid[y][x+i] = "─"
	}
}

func (t *TextDisplay) DrawCard(x, y int, card Card) {
	cardArt := getCardArt(card)

	colorWay := color.New(color.FgWhite).SprintFunc()

	switch card.suite {
	case HEART:
		colorWay = color.New(color.FgRed).SprintFunc()
	case DIAMOND:
		colorWay = color.New(color.FgMagenta).SprintFunc()
	case CLUB:
		colorWay = color.New(color.FgYellow).SprintFunc()
	case SPADE:
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

	t.DrawText(3, 2, player1.name)
	t.DrawPlayerHand(2, 3, player1, true)

	t.DrawText(3, 14, player2.name)
	t.DrawPlayerHand(2, 15, player2, true)

	t.DrawText(3, 26, player3.name)
	t.DrawPlayerHand(2, 27, player3, true)

	t.DrawText(3, 38, player4.name)
	t.DrawPlayerHand(2, 39, player4, true)
}

func (t *TextDisplay) DrawDealerArrow(game *Game) {
	dealerIndex := game.DealerIndex

	y := 5 + 12*dealerIndex
	x := 61
	t.DrawText(x, y, "<-- Dealer")
}

func (t *TextDisplay) DrawTurnArrow(game *Game) {
	playerIndex := game.PlayerIndex

	y := 6 + 12*playerIndex
	x := 61
	t.DrawText(x, y, "<-- Turn")
}

func (t *TextDisplay) DrawPlayedCards(game *Game) {
	t.DrawText(80, 2, "Played Cards")
	cards := game.PlayedCards

	if game.StateMachine.CurrentState.GetName() == DrawForDealer {
		if len(cards) > 0 {
			lastIndex := len(cards) - 1
			t.DrawCard(80, 5, *cards[lastIndex])
		}
	} else {
		for i, card := range cards {
			t.DrawCard(80, 5+(10*i), *card)
		}
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
	t.DrawText(120, 2, "Stats")
	t.DrawText(120, 3, "-----")
	t.DrawText(120, 4, fmt.Sprintf("Trump:          %s", game.Trump.ToString()))
	orderedPlayer := ""
	if game.OrderedPlayerIndex != -1 {
		orderedPlayer = game.Players[game.OrderedPlayerIndex].name
	}
	t.DrawText(120, 5, fmt.Sprintf("Ordered Up:     %s", orderedPlayer))
	t.DrawText(120, 6, fmt.Sprintf("Dealer:         %s", game.Players[game.DealerIndex].name))
	t.DrawText(120, 7, fmt.Sprintf("Turn:           %s", game.Players[game.PlayerIndex].name))
	turnedCardString := ""
	if game.TurnedCard != nil {
		turnedCardString = game.TurnedCard.ToString()
	}
	t.DrawText(120, 8, fmt.Sprintf("Turned Card:    %s", turnedCardString))
	t.DrawText(120, 9, fmt.Sprintf("Played Cards:   %d", len(game.PlayedCards)))
	t.DrawText(120, 10, fmt.Sprintf("State:          %s", game.StateMachine.CurrentState.GetName()))
	t.DrawText(120, 11, fmt.Sprintf("Cards in Deck:  %d", len(game.Deck.cards)))
	t.DrawText(120, 12, fmt.Sprintf("Team 1 Tricks:  %d", game.Players[0].tricksTaken+game.Players[2].tricksTaken))
	t.DrawText(120, 13, fmt.Sprintf("Team 2 Tricks:  %d", game.Players[1].tricksTaken+game.Players[3].tricksTaken))
	t.DrawText(120, 14, fmt.Sprintf("Team 1 Points:  %d", game.Players[0].pointsEarned))
	t.DrawText(120, 15, fmt.Sprintf("Team 2 Points:  %d", game.Players[1].pointsEarned))
}

func (t *TextDisplay) DrawBounds() {
	t.DrawVerticalLine(0, 0, 50)
	t.DrawHorizontalLine(0, 0, 165)
	t.DrawVerticalLine(75, 1, 50)
	t.DrawVerticalLine(95, 1, 50)
	t.DrawVerticalLine(115, 1, 50)
	t.DrawVerticalLine(165, 1, 50)
	t.DrawHorizontalLine(0, 50, 165)
	t.DrawRune(0, 0, '┌')
	t.DrawRune(0, 50, '└')
	t.DrawRune(165, 0, '┐')
	t.DrawRune(165, 50, '┘')
}

func (t *TextDisplay) DrawBoard(game *Game) {
	t.ClearDisplay()
	if game.StateMachine.CurrentState.GetName() == InitGame {
		return
	}
	t.DrawBounds()
	t.DrawPlayerHands(game)
	t.DrawDealerArrow(game)
	t.DrawTurnArrow(game)
	t.DrawPlayedCards(game)
	t.DrawTurnedCard(game)
	t.DrawStats(game)
	t.Render()
}

func ClearTerminal() {
	// Clear command based on the operating system
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	// Set the correct output device for the command
	cmd.Stdout = os.Stdout

	// Run the command to clear the screen
	cmd.Run()
}

func moveCursorHome() {
	fmt.Print("\033[H")
}
