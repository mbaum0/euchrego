package comms

import (
	"errors"
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
	*ConnMgr
}

func (m *CommsManager) EstablishClient(conn net.Conn) (*Client, error) {
	newClient := Client{}
	newClient.ConnMgr = NewConnMgr(conn)

	// read hello message from conn
	hello, err := newClient.ReadHelloMsg()
	if err != nil {
		logError("Error decoding message: ", err)
		newClient.SendAhoyMsg("", "invalid hello msg received")
		return nil, errors.New("invalid hello msg received")
	}
	logInfo("Received hello message from %s", hello.UserName)
	newClient.UserName = hello.UserName

	// if client didn't specify ID, establish a new client
	if hello.UserID == "" {

		// fail if we have the max no of clients
		if m.numClients() >= m.maxClients {
			newClient.SendAhoyMsg("", "max clients connected")
			return nil, errors.New("max clients connected")
		}

		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}
		newClient.ID = id
		m.clients[newClient.ID] = newClient
		logInfo("Player %s is a new player. Assigned ID: %s", hello.UserName, id)
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
			logWarn("Player %s sent a user ID that I do not recognize: %s", hello.UserName, newClient.ID)
			// send rejection
			newClient.SendAhoyMsg("", "invalid user ID")

			return nil, errors.New("invalid user ID")
		}
		logInfo("Player %s has returned with ID: %s", hello.UserName, newClient.ID)
	}

	// send client AhoyMsg
	logSuccess("Sending client ID %s to %s", newClient.ID, hello.UserName)
	newClient.SendAhoyMsg(newClient.ID, "")
	return &newClient, nil
}

func (m *CommsManager) HandleClient(client *Client) {
	for {
		_, kind, err := client.ReadMsg()
		if err != nil {
			// if EOF, remove client
			if err == io.EOF {
				logWarn("Client %s disconnected", client.UserName)
				client.Disconnect()
				return
			}
			logError("Error decoding message: ", err)
			continue
		}

		switch kind {
		case PingMsgKind:
			logInfo("Received ping message from %s", client.UserName)
			// send a PONG message back to the server
			client.SendPongMsg(client.ID)
		case LeaveMsgKind:
			logWarn("Recieved leave message from %s", client.UserName)
			m.Lock()
			delete(m.clients, client.ID)
			return

		default:
			logError("Received unknown message from %s", client.UserName)
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
		logSuccess("Accepted connection from %s", conn.RemoteAddr())

		newClient, err := m.EstablishClient(conn)
		if err != nil {
			conn.Close()
			continue
		}
		go m.HandleClient(newClient)
	}
}
