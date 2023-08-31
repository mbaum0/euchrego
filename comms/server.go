package comms

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
)

type CommsManager struct {
	numClients int
	port       string
	sync.Mutex
}

func NewCommsManager(port string) *CommsManager {
	return &CommsManager{0, port, sync.Mutex{}}
}

type Client struct {
	UserName string
	ID       string
	Reader   io.Reader
	Writer   io.Writer
}

func (m *CommsManager) Serve() {

	port := ":" + m.port
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// goroutine for reading from clients

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		// max number of clients is 4
		if m.numClients == 4 {
			fmt.Println("Max number of clients reached")
			conn.Close()
			continue
		}
		m.numClients++

		// read hello message from conn
		decoder := gob.NewDecoder(conn)
		var msg HelloMsg
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message: ", err)
			continue
		}
		fmt.Printf("Received hello message from %s\n", msg.UserName)

		// send client their ID
		encoder := gob.NewEncoder(conn)

		// generate a random ID for the client
		id, err := generateRandomID(32)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Sending client ID %s to %s\n", id, msg.UserName)
		newClient := Client{msg.UserName, id, conn, conn}
		clientIdMsg := ClientIdMsg{newClient.ID}
		encoder.Encode(clientIdMsg)
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
