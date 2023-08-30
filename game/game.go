package game

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	inputDevice := NewInputDevice(inputRequest, inputValue)
	gameMachine := GameMachine{gameBoard, inputDevice}
	display := NewTextDisplay()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	runner := fsm.New("Euchre FSM", gameMachine.InitGameState, fsm.Notifier(gameUpdated))

	// start game
	go runner.Run()

	// if the game has been updated and its been .5 seconds since the last update, redraw the board
	go func() {
		for v := range gameUpdated {
			if v {
				// sleep for .1 seconds
				time.Sleep(100 * time.Millisecond)
				display.DrawBoard(gameBoard)
			}
		}
	}()

	// if the game has requested input, get input from the user
	go func() {
		for v := range inputRequest {
			scanner := bufio.NewScanner(os.Stdin)
			// prompt the user for input
			// put the cursor at beginning of line under the board
			fmt.Printf("\033[%d;%dH", DISPLAY_HEIGHT+1, 1)
			// erase the current line
			fmt.Print("\r\033[K")
			fmt.Printf("%s", v)
			scanner.Scan()
			fmt.Print("\r\033[K")
			input := scanner.Text()
			gameMachine.inputChan <- input
		}
	}()

	sig := <-terminate
	ClearTerminal()
	fmt.Printf("Received %s, exiting...\n", sig)
}
