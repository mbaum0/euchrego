package comms

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
)

type CommsManager struct {
	clients map[string]Client
	port    string
	sync.Mutex
}

func NewCommsManager(port string) *CommsManager {
	return &CommsManager{make(map[string]Client), port, sync.Mutex{}}
}

type Client struct {
	UserName string
	ID       string
	Reader   io.Reader
	Writer   io.Writer
}

func (m *CommsManager) AddClient(c Client) {
	m.clients[c.ID] = c
}

func (m *CommsManager) RemoveClient(id string) {
	delete(m.clients, id)
}

func (m *CommsManager) Serve() {

	port := ":" + m.port
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// goroutine for reading from clients
	go func() {
		for {
			for _, c := range m.clients {
				decoder := gob.NewDecoder(c.Reader)
				var msg PingMsg
				err := decoder.Decode(&msg)
				if err != nil {
					// if EOF, remove client
					if err == io.EOF {
						fmt.Printf("Client %s disconnected\n", c.UserName)
						m.RemoveClient(c.ID)
						continue
					}
					fmt.Println("Error decoding message: ", err)
					continue
				}
				fmt.Printf("Received ping message from %s\n", c.UserName)
				// send a PONG message back to the server
				pongMsg := PongMsg(msg)
				encoder := gob.NewEncoder(c.Writer)
				encoder.Encode(pongMsg)
			}
		}
	}()

	for {
		client, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Accepted connection from %s\n", client.RemoteAddr())
		m.Lock()
		// read hello message from client
		decoder := gob.NewDecoder(client)
		var msg HelloMsg
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			continue
		}
		fmt.Printf("Received hello message from %s\n", msg.UserName)

		// send client their ID
		encoder := gob.NewEncoder(client)

		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Sending client ID %s to %s\n", id, msg.UserName)
		newClient := Client{msg.UserName, id, client, client}
		clientIdMsg := ClientIdMsg{newClient.ID}
		encoder.Encode(clientIdMsg)

		m.AddClient(newClient)
		m.Unlock()
	}
}
