package comms

import (
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
	Sock     *Sock
}

func (m *CommsManager) EstablishClient(conn net.Conn) (*Client, error) {
	newClient := Client{}
	newClient.Sock = NewSock(conn)

	// read hello message from conn
	hello, err := newClient.Sock.ReadHelloMsg()
	if err != nil {
		fmt.Println("Error decoding message: ", err)
		newClient.Sock.SendAhoyMsg("", "invalid hello msg received")
		return nil, errors.New("invalid hello msg received")
	}
	fmt.Printf("Received hello message from %s\n", hello.UserName)
	newClient.UserName = hello.UserName

	// if client didn't specify ID, establish a new client
	if hello.UserID == "" {

		// fail if we have the max no of clients
		if m.numClients() >= m.maxClients {
			newClient.Sock.SendAhoyMsg("", "max clients connected")
			return nil, errors.New("max clients connected")
		}

		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}
		newClient.ID = id
		m.clients[newClient.ID] = newClient
		fmt.Printf("Player %s is a new player. Assigned ID: %s\n", hello.UserName, id)
	} else {
		newClient.ID = hello.UserID

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
			fmt.Printf("Player %s sent a user ID that I do not recognize: %s\n", hello.UserName, newClient.ID)
			// send rejection
			newClient.Sock.SendAhoyMsg("", "invalid user ID")

			return nil, errors.New("invalid user ID")
		}
		fmt.Printf("Player %s has returned with ID: %s\n", hello.UserName, newClient.ID)
	}

	// send client AhoyMsg
	fmt.Printf("Sending client ID %s to %s\n", newClient.ID, hello.UserName)
	newClient.Sock.SendAhoyMsg(newClient.ID, "")
	return &newClient, nil
}

func (m *CommsManager) HandleClient(client *Client) {
	for {
		_, kind, err := client.Sock.ReadMsg()
		if err != nil {
			// if EOF, remove client
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", client.UserName)
				client.Sock.Disconnect()
				return
			}
			fmt.Println("Error decoding message: ", err)
			continue
		}

		switch kind {
		case PingMsgKind:
			fmt.Printf("Received ping message from %s\n", client.UserName)
			// send a PONG message back to the server
			client.Sock.SendPongMsg(client.ID)
		case LeaveMsgKind:
			fmt.Printf("Recieved leave message from %s\n", client.UserName)
			m.Lock()
			delete(m.clients, client.ID)
			return

		default:
			fmt.Printf("Received unknown message from %s\n", client.UserName)
		}
	}
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
		go m.HandleClient(newClient)
	}
}
