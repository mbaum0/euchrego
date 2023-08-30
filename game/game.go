package game

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbaum0/euchrego/godeck"
)

type Game struct {
	StateMachine       StateMachine
	Deck               godeck.EuchreDeck
	Players            [4]*Player
	DealerIndex        int
	PlayerIndex        int
	TurnedCard         godeck.Card
	PlayedCards        []godeck.Card
	Trump              godeck.Suit
	OrderedPlayerIndex int // the player who ordered it up
	logs               []string
	RandSeed           int64
}

func NewGame() Game {
	game := Game{}
	game.StateMachine = NewStateMachine()
	game.PlayedCards = nil
	game.logs = make([]string, 0)
	game.OrderedPlayerIndex = -1
	game.DealerIndex = 0
	game.PlayerIndex = 0
	game.RandSeed = int64(1)
	game.Players[0] = InitPlayer("Player 1", 0)
	game.Players[1] = InitPlayer("Player 2", 1)
	game.Players[2] = InitPlayer("Player 3", 2)
	game.Players[3] = InitPlayer("Player 4", 3)
	game.TurnedCard = godeck.EmptyCard()
	game.Trump = godeck.None
	game.PlayedCards = make([]godeck.Card, 0)
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

func (g *Game) PlayCard(card godeck.Card) {
	g.PlayedCards = append(g.PlayedCards, card)
}

func (g *Game) ReturnPlayedCards() {
	g.Deck.ReturnCards(g.PlayedCards)
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
			time.Sleep(100 * time.Millisecond)

			if game.StateMachine.CurrentState.GetName() == EndGame {
				break
			}
		}
	}()

	sig := <-terminate
	ClearTerminal()
	fmt.Printf("Received %s, exiting...\n", sig)

}
