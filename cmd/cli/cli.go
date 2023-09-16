package main

import (
	"fmt"
	"os"

	"github.com/mbaum0/euchrego/termui"
)

func cli(d *termui.TermUI) {
	for {
		d.Render()
		_ = d.PollForInput()
	}
}

func main() {
	d, err := termui.NewTermUI(termui.Size(100, 30), termui.EnableInput())
	if err != nil {
		fmt.Printf("Error occured %s", err)
		os.Exit(1)
	}

	d.DrawRect(d.Left(), d.Top(), d.Width(), d.Height())
	d.DrawHorizontalLine(d.Left(), d.Bottom()-2, d.Width())
	d.DrawRune(d.Left(), d.Bottom()-2, '├')
	d.DrawRune(d.Right(), d.Bottom()-2, '┤')
	d.DrawText(d.Left()+2, d.Bottom()-1, "Status: Connected")
	d.DrawTextRightAligned(d.Right()-1, d.Bottom()-1, "Turn: Player 1")
	d.DrawTitle("EuchreGo!")
	cli(d)
}
