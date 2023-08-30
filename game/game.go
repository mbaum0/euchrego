package game

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbaum0/euchrego/fsm"
)

// func logToFile(format string, args ...interface{}) {
// 	file, err := os.OpenFile("log.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		fmt.Println("Error opening log file: ", err)
// 		return
// 	}
// 	defer file.Close()
// 	format += "\n"
// 	fmt.Fprintf(file, format, args...)
// }

// func DeleteLogFile() {
// 	os.Remove("log.out")
// }

func Run() {
	gameBoard := NewGameBoard()
	gameUpdated := make(chan bool, 1)
	inputRequest := make(chan string, 1)
	inputValue := make(chan string, 1)
	gameMachine := GameMachine{&gameBoard, inputRequest, inputValue}
	display := NewTextDisplay()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	runner := fsm.New("Euchre FSM", gameMachine.InitGameState, fsm.Notifier(gameUpdated))

	// start game
	go runner.Run()

	// if the game has been updated, display the game
	go func() {
		for v := range gameUpdated {
			if v {
				display.DrawBoard(&gameBoard)
			}
		}
	}()

	// if the game has requested input, get input from the user
	go func() {
		for v := range inputRequest {
			display.DrawBoard(&gameBoard)
			fmt.Printf("%s: ", v)
			var input string
			fmt.Scanln(&input)
			gameMachine.Input <- input
		}
	}()

	sig := <-terminate
	ClearTerminal()
	fmt.Printf("Received %s, exiting...\n", sig)

}
