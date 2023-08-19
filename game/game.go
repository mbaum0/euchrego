package game

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Game struct {
	StateMachine       StateMachine
	Deck               Deck
	Players            [4]*Player
	DealerIndex        int
	PlayerIndex        int
	TurnedCard         *Card
	PlayedCards        []*Card
	Trump              Suite
	OrderedPlayerIndex int // the player who ordered it up
	logs               []string
}

func NewGame() Game {
	game := Game{}
	game.StateMachine = NewStateMachine()
	game.PlayedCards = nil
	game.logs = make([]string, 0)
	game.OrderedPlayerIndex = -1
	game.DealerIndex = -1
	game.PlayerIndex = -1
	return game
}

func (g *Game) Log(format string, args ...interface{}) {
	g.logs = append(g.logs, fmt.Sprintf(format, args...))
	logToFile(format, args...)
}

func logToFile(format string, args ...interface{}) {
	file, err := os.OpenFile("log.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file: ", err)
		return
	}
	defer file.Close()
	format += "\n"
	fmt.Fprintf(file, format, args...)
}

func DeleteLogFile() {
	os.Remove("log.out")
}

func (g *Game) PlayCard(card *Card) {
	g.PlayedCards = append(g.PlayedCards, card)
}

func (g *Game) ReturnPlayedCards() {
	g.Deck.ReturnCards(&g.PlayedCards)
}

func (g *Game) NextPlayer() {
	g.PlayerIndex = (g.PlayerIndex + 1) % 4
}

func Run() {
	game := NewGame()
	display := NewTextDisplay()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start game
	go func() {
		for {
			game.StateMachine.Step(&game)
			display.DrawBoard(&game)
			// delay for .5 seconds for animation
			time.Sleep(500 * time.Millisecond)

			if game.StateMachine.CurrentState.GetName() == EndGame {
				break
			}
		}
	}()

	sig := <-terminate
	ClearTerminal()
	fmt.Printf("Received %s, exiting...\n", sig)

}
