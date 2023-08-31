package comms

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/mbaum0/euchrego/fsm"
)

type GameClient struct {
	playerId         string // given to us by the server
	playerName       string // set by the player
	serverHost       string
	serverPort       string
	serverReader     io.Reader
	serverWriter     io.Writer
	pingInterval     time.Duration
	lastPingSendTime time.Time
}

func NewGameClient(playerName string, serverHost string, serverPort string, pingInterval int) *GameClient {
	return &GameClient{"", playerName, serverHost, serverPort, nil, nil, time.Duration(pingInterval) * time.Second, time.Now()}
}

func (gc *GameClient) log(format string, args ...interface{}) {
	format += "\n"
	fmt.Printf(format, args...)
}

// get information from the user
func (gc *GameClient) StartState() (fsm.StateFunc, error) {
	return gc.ConnectToServerState, nil
}

func (gc *GameClient) ConnectToServerState() (fsm.StateFunc, error) {
	connStr := gc.serverHost + ":" + gc.serverPort
	gc.log("Connecting to server at %s", connStr)
	server, err := net.Dial("tcp", connStr)
	if err != nil {
		gc.log("Error connecting to server: %s", err)
		// sleep for 5 seconds and try again
		time.Sleep(5 * time.Second)
		return gc.ConnectToServerState, nil
	}
	gc.log("Connected to server at %s", connStr)
	gc.serverReader = server
	gc.serverWriter = server
	return gc.HelloState, nil
}

func (gc *GameClient) HelloState() (fsm.StateFunc, error) {
	gc.log("Sending hello message to server")
	encoder := gob.NewEncoder(gc.serverWriter)
	hello := HelloMsg{gc.playerName}
	err := encoder.Encode(hello)
	if err != nil {
		gc.log("Error sending hello message to server: %s", err)
		return nil, err
	}
	return gc.Wait4PlayerIdState, nil
}

func (gc *GameClient) Wait4PlayerIdState() (fsm.StateFunc, error) {
	gc.log("Waiting for player ID from server")
	var clientIdMsg ClientIdMsg
	decoder := gob.NewDecoder(gc.serverReader)
	err := decoder.Decode(&clientIdMsg)
	if err != nil {
		// wait for a 1 seconds and try again
		time.Sleep(1 * time.Second)
		return gc.Wait4PlayerIdState, nil
	}
	gc.playerId = clientIdMsg.UserID
	return gc.SendPingState, nil
}

func (gc *GameClient) SendPingState() (fsm.StateFunc, error) {
	// send a ping to the server
	gc.log("Sending ping to server")
	pingMsg := PingMsg{gc.playerId}
	encoder := gob.NewEncoder(gc.serverWriter)
	err := encoder.Encode(pingMsg)
	if err != nil {
		gc.log("Error sending ping to server: %s", err)
		return nil, err
	}
	return gc.Wait4PongState, nil
}

func (gc *GameClient) Wait4PongState() (fsm.StateFunc, error) {
	// wait for a pong from the server
	gc.log("Waiting for pong from server")
	var pongMsg PongMsg
	decoder := gob.NewDecoder(gc.serverReader)
	err := decoder.Decode(&pongMsg)
	if err != nil {
		// wait for a 1 seconds and try again
		time.Sleep(1 * time.Second)
		return gc.Wait4PongState, nil
	}
	pingTime := time.Since(gc.lastPingSendTime)
	gc.log("Got pong from server after %d ms", pingTime.Milliseconds())
	gc.lastPingSendTime = time.Now()
	time.Sleep(gc.pingInterval)

	return gc.SendPingState, nil
}
