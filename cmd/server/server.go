package main

import (
	"github.com/mbaum0/euchrego/comms"
)

func main() {
	cm := comms.NewCommsManager("8765")
	cm.Serve()
}
