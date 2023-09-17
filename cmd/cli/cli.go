package main

import (
	"fmt"
	"os"

	"github.com/mbaum0/euchrego/godeck"
	"github.com/mbaum0/euchrego/termui"
)

type GameStatus struct {
	isErr bool
	msg   string
}

type GameState struct {
	connected   bool
	playerTurn  int
	myIndex     int
	hand        []godeck.Card
	status      GameStatus
	deck        *godeck.EuchreDeck
	turnedCard  godeck.Card
	playedCards []godeck.Card
	players     [4]string
}

func NewGameState() *GameState {
	gs := GameState{}
	gs.connected = false
	gs.playerTurn = 0
	gs.status = GameStatus{false, "Waiting for player..."}
	gs.hand = make([]godeck.Card, 0)
	gs.deck = godeck.NewEuchreDeck(godeck.PreShuffled())
	gs.players = [4]string{"Abby", "Mike", "Michael", "Dolores"}
	gs.myIndex = 0
	return &gs
}

func handleInput(input string, gs *GameState) {
	switch input {
	case "connect":
		gs.connected = true
		gs.status.isErr = false
		gs.status.msg = ""
	case "disconnect":
		gs.connected = false
		gs.status.isErr = false
		gs.status.msg = ""
	case "hit":
		if len(gs.hand) >= 5 {
			gs.status.isErr = true
			gs.status.msg = "Can only have 5 cards at a time."
			break
		}
		c, err := gs.deck.DrawCard()
		if err != nil {
			gs.status.isErr = true
			gs.status.msg = err.Error()
			break
		}
		gs.hand = append(gs.hand, c)
		gs.status.isErr = false
		gs.status.msg = fmt.Sprintf("Drew the %s", c.String())

	case "play":
		if len(gs.hand) <= 0 {
			gs.status.isErr = true
			gs.status.msg = "Out of cards to play!"
			break
		}
		c := gs.hand[len(gs.hand)-1]
		gs.hand = gs.hand[:len(gs.hand)-1]
		gs.deck.ReturnCard(c)
		gs.status.isErr = false
		gs.status.msg = fmt.Sprintf("Played the %s", c.String())
	case "trump":
		c, err := gs.deck.DrawCard()
		if err != nil {
			gs.status.isErr = true
			gs.status.msg = err.Error()
			break
		}
		gs.turnedCard = c
		gs.status.isErr = false
		gs.status.msg = fmt.Sprintf("%s was turned!", c.String())
	case "next":
		gs.playerTurn += 1
		gs.playerTurn %= 4
	default:
		gs.status.isErr = true
		gs.status.msg = "invalid input"
	}
}

func cli(d *termui.TermUI) {
	gs := NewGameState()
	for {
		updateView(d, gs)
		input := d.PollForInput()
		handleInput(input, gs)

	}
}

func updateStatusBarView(ui *termui.TermUI, gs *GameState) {
	if gs.connected {
		ui.DrawText("Connected", ui.Right()-28, ui.Bottom()-1, termui.Color(termui.Green), termui.Width(16))
	} else {
		ui.DrawText("Not Connected", ui.Right()-28, ui.Bottom()-1, termui.Color(termui.Red), termui.Width(16))
	}

	if gs.status.isErr {
		ui.DrawText(gs.status.msg, ui.Left()+2, ui.Bottom()-1, termui.Color(termui.Red), termui.Width(70))
	} else {
		ui.DrawText(gs.status.msg, ui.Left()+2, ui.Bottom()-1, termui.Color(termui.Yellow), termui.Width(70))
	}

	ui.DrawText(gs.players[gs.myIndex], ui.Right()-5, ui.Bottom()-1, termui.Width(9), termui.Color(termui.Blue), termui.Justify(termui.Center))

}

func updateTrumpView(ui *termui.TermUI, gs *GameState) {
	ui.DrawCard(ui.Right()-20, ui.Top()+3, gs.turnedCard)
	ui.DrawText("turned card", ui.Right()-20, ui.Top()+12)
}

func updateHandView(ui *termui.TermUI, gs *GameState) {
	ui.DrawHand(ui.Left()+4, ui.Bottom()-12, gs.hand, true, false)
}

func updatePlayedCardView(ui *termui.TermUI, gs *GameState) {
	// ui.DrawText("Played Cards", ui.Left()+4, ui.Top()+2, termui.Color(termui.Red))
	for i := 0; i < 4; i++ {
		ui.DrawCard(ui.Left()+4+(12*i), ui.Top()+3, godeck.EmptyCard())
		pColor := termui.White
		if gs.playerTurn == i {
			pColor = termui.Blue
		}
		ui.DrawText(gs.players[i], ui.Left()+9+(12*i), ui.Top()+12, termui.Justify(termui.Center), termui.Color(pColor))
	}
	ui.DrawHorizontalLine(ui.Left(), ui.Top()+14, ui.Width())
	ui.DrawRune('├', ui.Left(), ui.Top()+14)
	ui.DrawRune('┤', ui.Right(), ui.Top()+14)
}

func updateView(ui *termui.TermUI, gs *GameState) {
	updateStatusBarView(ui, gs)
	updateHandView(ui, gs)
	updateTrumpView(ui, gs)
	updatePlayedCardView(ui, gs)
	ui.Render()
}

func main() {
	d, err := termui.NewTermUI(termui.Size(100, 30), termui.EnableInput())
	if err != nil {
		fmt.Printf("Error occured %s", err)
		os.Exit(1)
	}

	d.DrawRect(d.Left(), d.Top(), d.Width(), d.Height())
	d.DrawHorizontalLine(d.Left(), d.Bottom()-2, d.Width())
	d.DrawRune('├', d.Left(), d.Bottom()-2)
	d.DrawRune('┤', d.Right(), d.Bottom()-2)
	d.DrawRune('┬', d.Right()-38, d.Bottom()-2)
	d.DrawRune('│', d.Right()-38, d.Bottom()-1)
	d.DrawRune('┴', d.Right()-38, d.Bottom())
	d.DrawText("Status: ", d.Right()-36, d.Bottom()-1)
	// d.DrawText("You: ", d.Right()-15, d.Bottom()-1)
	d.DrawRune('┬', d.Right()-13, d.Bottom()-2)
	d.DrawRune('│', d.Right()-13, d.Bottom()-1)
	d.DrawRune('┴', d.Right()-13, d.Bottom())
	d.DrawTitle("EuchreGo!")
	cli(d)
}
