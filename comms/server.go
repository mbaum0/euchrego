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
	ID            string
	Reader        io.Reader
	Writer        io.Writer
	cardSelectMsg *CardSelectMsg
	boolMsg       *BoolMsg
}

func (m *CommsManager) AddClient(c Client) {
	m.clients[c.ID] = c
}

func (m *CommsManager) RemoveClient(id string) {
	delete(m.clients, id)
}

func (m *CommsManager) GetBoolMsgForClient(id string) *BoolMsg {
	return m.clients[id].boolMsg
}

func (m *CommsManager) GetCardSelectMsgForClient(id string) *CardSelectMsg {
	return m.clients[id].cardSelectMsg
}

func (m *CommsManager) Serve() {
	gob.Register(CardSelectMsg{})
	gob.Register(BoolMsg{})
	gob.Register(HelloMsg{})
	gob.Register(PingMsg{})
	gob.Register(PongMsg{})

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
					fmt.Println("Error decoding message: ", err)
					continue
				}
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
		m.Lock()
		decoder := gob.NewDecoder(client)
		var hello HelloMsg
		err = decoder.Decode(&hello)
		if err != nil {
			continue
		}
		newClient := Client{hello.UserName, client, client, nil, nil}

		// send client their ID
		encoder := gob.NewEncoder(client)

		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}
		encoder.Encode(ClientIDMsg{id})

		m.AddClient(newClient)
		m.Unlock()
	}
}
