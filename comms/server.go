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
	maxClients int
	port       string
	clients    map[string]Client
	sync.Mutex
}

func (m *CommsManager) numClients() int {
	return len(m.clients)
}

func NewCommsManager(port string) *CommsManager {
	return &CommsManager{2, port, make(map[string]Client), sync.Mutex{}}
}

type Client struct {
	UserName string
	ID       string
	Reader   io.Reader
	Writer   io.Writer
}

func (m *CommsManager) EstablishClient(conn net.Conn) (*Client, error) {
	encoder := gob.NewEncoder(conn)
	var ahoyMsg AhoyMsg
	if m.numClients() >= m.maxClients {
		ahoyMsg = AhoyMsg{"", "max clients connected."}
		encoder.Encode(ahoyMsg)
		return nil, errors.New("max clients connected.")
	}

	// read hello message from conn
	decoder := gob.NewDecoder(conn)
	var msg HelloMsg
	err := decoder.Decode(&msg)
	if err != nil {
		fmt.Println("Error decoding message: ", err)
		ahoyMsg = AhoyMsg{"", "invalid hello msg received."}
		encoder.Encode(ahoyMsg)
		return nil, errors.New("invalid hello msg received.")
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
		m.clients[newClient.ID] = newClient
		fmt.Printf("Player %s is a new player. Assigned ID: %s\n", msg.UserName, id)
	} else {
		newClient.ID = msg.UserID

		valid := false
		// reject client if ID is invalid
		for clientId := range m.clients {
			if clientId == newClient.ID {
				// match! update the saved client
				m.clients[clientId] = newClient
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Player %s sent a user ID that I do not recognize: %s\n", msg.UserName, newClient.ID)
			// send rejection
			ahoyMsg = AhoyMsg{"", "invalid user ID."}
			encoder.Encode(ahoyMsg)
			return nil, errors.New("invalid user ID.")
		}
		fmt.Printf("Player %s has returned with ID: %s\n", msg.UserName, newClient.ID)
	}

	newClient.Reader = conn
	newClient.Writer = conn
	newClient.UserName = msg.UserName

	// send client AhoyMsg
	fmt.Printf("Sending client ID %s to %s\n", newClient.ID, msg.UserName)
	ahoyMsg = AhoyMsg{newClient.ID, ""}
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
			conn.Close()
			continue
		}
		go func() {
			decoder := gob.NewDecoder(newClient.Reader)
			encoder := gob.NewEncoder(newClient.Writer)
			for {
				var msg PingMsg
				err := decoder.Decode(&msg)
				if err != nil {
					// if EOF, remove client
					if err == io.EOF {
						fmt.Printf("Client %s disconnected\n", newClient.UserName)
						conn.Close()
						//m.Lock()
						//delete(m.clients, newClient.ID)
						//m.Unlock()
						return
					}
					fmt.Println("Error decoding message: ", err)
					continue
				}
				fmt.Printf("Received ping message from %s\n", newClient.UserName)
				// send a PONG message back to the server
				pongMsg := PongMsg(msg)
				encoder.Encode(pongMsg)
			}
		}()
	}
}
