package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mbaum0/euchrego/comms"
)

func main() {

	// channel for getting ctrl c
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		client := comms.NewGameClient("player0", "localhost", "8765")

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			userInput, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Something weird happend: %s\n", err)
				continue
			}
			userInput = strings.TrimSpace(userInput)

			switch userInput {
			case "connect":
				client.ConnectToServer("")
			case "reconnect":
				client.ReconnectToServer()
			case "disconnect":
				client.DisconnectFromServer()
			case "ping":
				client.SendPing()
			case "leave":
				client.LeaveServer()
			default:
				fmt.Println("Invalid command")
			}
		}

	}()

	// wait for ctrl c
	<-terminate
}
