package comms

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func Connect(host string, port string) {
	// connect to server and send hello message
	connStr := host + ":" + port
	server, err := net.Dial("tcp", connStr)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	// send hello message
	encoder := gob.NewEncoder(server)
	hello := HelloMsg{"mikeb"}
	err = encoder.Encode(hello)
	if err != nil {
		panic(err)
	}

	// receive client ID
	decoder := gob.NewDecoder(server)
	var clientID ClientIDMsg
	err = decoder.Decode(&clientID)
	if err != nil {
		panic(err)
	}

	// goroutine that sends pings every 5 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			pingMsg := Message(PingMsg(clientID))
			encoder := gob.NewEncoder(server)
			encoder.Encode(pingMsg)
		}
	}()

	// goroutine that listens for messages from the server
	for {
		decoder := gob.NewDecoder(server)
		var msg interface{}
		err := decoder.Decode(&msg)
		if err != nil {
			panic(err)
		}
		switch msg := msg.(type) {
		case PingMsg:
			pingMsg := msg
			// send a PONG message back to the server
			pongMsg := PongMsg(pingMsg)
			encoder := gob.NewEncoder(server)
			encoder.Encode(pongMsg)
		case PongMsg:
			pongMsg := msg
			// print the pongMsg
			fmt.Printf("PONG from: %s\n", pongMsg.UserID)
		}
	}
}
