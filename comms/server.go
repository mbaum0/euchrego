package comms

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type CommsManager struct {
	numClients int
	maxClients int
	port       string
	clients    []Client
	sync.Mutex
}

func NewCommsManager(port string) *CommsManager {
	return &CommsManager{0, 4, port, make([]Client, 0), sync.Mutex{}}
}

type Client struct {
	UserName string
	ID       string
	Reader   io.Reader
	Writer   io.Writer
}

func (m *CommsManager) EstablishClient(conn net.Conn) (*Client, error) {
	encoder := gob.NewEncoder(conn)
	if m.numClients >= m.maxClients {
		// send rejection message
		denyMsg := ServerDenyMsg{"max clients connected."}
		encoder.Encode(denyMsg)
		return nil, errors.New("max clients connected.")
	}

	// read hello message from conn
	decoder := gob.NewDecoder(conn)
	var msg HelloMsg
	err := decoder.Decode(&msg)
	if err != nil {
		fmt.Println("Error decoding message: ", err)
		denyMsg := ServerDenyMsg{"invalid hello msg received."}
		encoder.Encode(denyMsg)
	}
	fmt.Printf("Received hello message from %s\n", msg.UserName)

	var newClient Client
	// if client didn't specify ID, generate a new one
	if msg.UserID == "" {
		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}
		newClient.ID = id
	} else {
		newClient.ID = msg.UserID

		valid := false
		// reject client if ID is invalid
		for i, cli := range m.clients {
			if cli.ID == newClient.ID {
				// match! update the saved client
				m.clients[i] = newClient
				valid = true
				break
			}
		}
		if !valid {
			// send rejection
			denyMsg := ServerDenyMsg{"invalid user ID."}
			encoder.Encode(denyMsg)
			return nil, errors.New("invalid user ID.")
		}
	}

	newClient.Reader = conn
	newClient.Writer = conn
	newClient.UserName = msg.UserName

	// send client AhoyMsg
	fmt.Printf("Sending client ID %s to %s\n", newClient.ID, msg.UserName)
	ahoyMsg := AhoyMsg{newClient.ID}
	encoder.Encode(ahoyMsg)
	return &newClient, nil
}

func (m *CommsManager) Serve() {

	port := ":" + m.port
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		newClient, err := m.EstablishClient(conn)
		if err != nil {
			continue
		}
		go func() {
			for {
				decoder := gob.NewDecoder(newClient.Reader)
				var msg PingMsg
				err := decoder.Decode(&msg)
				if err != nil {
					// if EOF, remove client
					if err == io.EOF {
						fmt.Printf("Client %s disconnected\n", newClient.UserName)
						conn.Close()
						m.Lock()
						m.numClients--
						m.Unlock()
						return
					}
					fmt.Println("Error decoding message: ", err)
					continue
				}
				fmt.Printf("Received ping message from %s\n", newClient.UserName)
				// send a PONG message back to the server
				pongMsg := PongMsg(msg)
				encoder := gob.NewEncoder(newClient.Writer)
				encoder.Encode(pongMsg)
			}
		}()
	}
}
