package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mbaum0/euchrego/comms"
	"github.com/mbaum0/euchrego/fsm"
)

func main() {
	// get port and host from command line args

	// if len(os.Args) != 5 {
	// 	fmt.Println("Usage: euchrego <name> <host> <port> <ping interval>")
	// 	return
	// }

	// name := os.Args[1]
	// host := os.Args[2]
	// port := os.Args[3]
	// pingInterval, err := strconv.Atoi(os.Args[4])
	// if err != nil {
	// 	fmt.Println("Error converting ping interval to int: ", err)
	// 	return
	// }

	// gameClient := comms.NewGameClient(name, host, port, pingInterval)
	// runner := fsm.New("Client FSM", gameClient.StartState)
	// runner.Run()

	// get number of clients to start
	if len(os.Args) != 2 {
		fmt.Println("Usage: euchrego <num clients>")
		return
	}

	numClients, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error converting num clients to int: ", err)
		return
	}

	// channel for getting ctrl c
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start clients
	for i := 0; i < numClients; i++ {
		// generate random integer between 500 and 10000
		pingInterval := generateRandomInterval(100, 10000)
		gameClient := comms.NewGameClient("Player "+strconv.Itoa(i), "localhost", "8765", pingInterval)
		runner := fsm.New("Client FSM", gameClient.StartState)
		go runner.Run()
	}

	// wait for ctrl c
	<-terminate
}

func generateRandomInterval(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer within the range [min, max)
	randomNumber := r.Intn(max-min) + min
	return randomNumber
}
