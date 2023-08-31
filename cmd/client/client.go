package main

import (
	"fmt"
	"os"

	"github.com/mbaum0/euchrego/comms"
)

func main() {
	// get port and host from command line args

	if len(os.Args) != 3 {
		fmt.Println("Usage: euchrego <host> <port>")
		return
	}

	host := os.Args[1]
	port := os.Args[2]
	comms.Connect(host, port)
}
