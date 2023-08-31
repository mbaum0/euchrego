package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mbaum0/euchrego/comms"
	"github.com/mbaum0/euchrego/fsm"
)

func main() {
	// get port and host from command line args

	if len(os.Args) != 5 {
		fmt.Println("Usage: euchrego <name> <host> <port> <ping interval>")
		return
	}

	name := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]
	pingInterval, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("Error converting ping interval to int: ", err)
		return
	}

	gameClient := comms.NewGameClient(name, host, port, pingInterval)
	runner := fsm.New("Client FSM", gameClient.StartState)
	runner.Run()
}
