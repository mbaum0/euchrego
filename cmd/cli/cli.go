package main

import (
	"fmt"
	"os"

	"github.com/mbaum0/euchrego/termui"
)

type GameStatus struct {
	isErr bool
	msg   string
}

type GameState struct {
	connected  bool
	playerTurn string
	status     GameStatus
}

func NewGameState() *GameState {
	gs := GameState{false, "Player 1", GameStatus{false, "Waiting for player..."}}
	return &gs
}

func handleInput(input string, gs *GameState) {
	switch input {
	case "connect":
		gs.connected = true
	case "disconnect":
		gs.connected = false
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

}

func updateView(ui *termui.TermUI, gs *GameState) {
	if gs.connected {
		ui.DrawText("Connected", ui.Right()-32, ui.Bottom()-1, termui.Color(termui.Green), termui.Width(15))
	} else {
		ui.DrawText("Not Connected", ui.Right()-32, ui.Bottom()-1, termui.Color(termui.Red), termui.Width(15))
	}

	if gs.status.isErr {
		ui.DrawText(gs.status.msg, ui.Left()+2, ui.Bottom()-1, termui.Color(termui.Red), termui.Width(55))
	} else {
		ui.DrawText(gs.status.msg, ui.Left()+2, ui.Bottom()-1, termui.Width(55))
	}

	ui.DrawText(gs.playerTurn, ui.Right()-1, ui.Bottom()-1, termui.Justify(termui.Right), termui.Color(termui.Blue), termui.Width(9))

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
	d.DrawRune('┬', d.Right()-42, d.Bottom()-2)
	d.DrawRune('│', d.Right()-42, d.Bottom()-1)
	d.DrawRune('┴', d.Right()-42, d.Bottom())
	d.DrawText("Status: ", d.Right()-40, d.Bottom()-1)
	d.DrawText("Turn: ", d.Right()-15, d.Bottom()-1)
	d.DrawRune('┬', d.Right()-17, d.Bottom()-2)
	d.DrawRune('│', d.Right()-17, d.Bottom()-1)
	d.DrawRune('┴', d.Right()-17, d.Bottom())
	d.DrawTitle("EuchreGo!")
	cli(d)
}
