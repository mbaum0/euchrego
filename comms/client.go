package comms

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"

	"github.com/mbaum0/euchrego/fsm"
)

type GameClient struct {
	playerId           string // given to us by the server
	playerName         string // set by the player
	serverHost         string
	serverPort         string
	serverReader       io.Reader
	serverWriter       io.Writer
	pingInterval       time.Duration
	lastPingSendTime   time.Time
	connectionAttempts int
}

func NewGameClient(playerName string, serverHost string, serverPort string, pingInterval int) *GameClient {
	gc := GameClient{}
	gc.playerName = playerName
	gc.serverHost = serverHost
	gc.serverPort = serverPort
	gc.pingInterval = time.Duration(pingInterval) * time.Millisecond
	gc.lastPingSendTime = time.Now()
	gc.connectionAttempts = 0
	return &gc

}

func (gc *GameClient) log(format string, args ...interface{}) {
	format = fmt.Sprintf("client: %s -\t%s\n", gc.playerName, format)
	fmt.Printf(format, args...)
}

// get information from the user
func (gc *GameClient) StartState() (fsm.StateFunc, error) {
	return gc.ConnectToServerState, nil
}

func (gc *GameClient) ResetServerConnectionState() (fsm.StateFunc, error) {
	gc.serverReader = nil
	gc.serverWriter = nil
	return gc.ConnectToServerState, nil
}

func (gc *GameClient) ConnectToServerState() (fsm.StateFunc, error) {
	connStr := gc.serverHost + ":" + gc.serverPort
	gc.log("Connecting to server at %s", connStr)
	server, err := net.Dial("tcp", connStr)
	if err != nil {
		switch {
		case errors.Is(err, syscall.ECONNREFUSED):
			gc.log("Connection refused. Is server up?")
			time.Sleep(5 * time.Second)
			return gc.ConnectToServerState, nil
		default:
			gc.log("Error connecting to server: %s", err)
			// sleep for 5 seconds and try again
			time.Sleep(5 * time.Second)
			return gc.ConnectToServerState, nil
		}
	}

	gc.log("Connected to server at %s", connStr)
	gc.serverReader = server
	gc.serverWriter = server

	if gc.playerId == "" {
		// this is the first time we are joining
		return gc.HelloState, nil
	}
	return gc.HelloAgainState, nil

}

func (gc *GameClient) HelloState() (fsm.StateFunc, error) {
	gc.log("Sending hello message to server")
	encoder := gob.NewEncoder(gc.serverWriter)
	hello := HelloMsg{gc.playerName, ""}
	err := encoder.Encode(hello)
	if err != nil {
		// if broken pipe, reconnect to server
		if err == io.ErrClosedPipe {
			// if we've tried to connect 5 times, give up
			if gc.connectionAttempts == 5 {
				gc.log("Max connection attempts reached. Giving up.")
				return nil, err
			}
			gc.connectionAttempts++
			// sleep for 5 seconds and try again
			time.Sleep(5 * time.Second)
			return gc.ConnectToServerState, nil
		}
		gc.log("Error sending hello message to server: %s", err)
		return nil, err
	}
	return gc.Wait4PlayerIdState, nil
}

func (gc *GameClient) HelloAgainState() (fsm.StateFunc, error) {
	gc.log("Sending hello message to server")
	encoder := gob.NewEncoder(gc.serverWriter)
	hello := HelloMsg{gc.playerName, gc.playerId}
	err := encoder.Encode(hello)
	if err != nil {
		// if broken pipe, reconnect to server
		if err == io.ErrClosedPipe {
			// if we've tried to connect 5 times, give up
			if gc.connectionAttempts == 5 {
				gc.log("Max connection attempts reached. Giving up.")
				return nil, err
			}
			gc.connectionAttempts++
			// sleep for 5 seconds and try again
			time.Sleep(5 * time.Second)
			return gc.ConnectToServerState, nil
		}
		gc.log("Error sending hello message to server: %s", err)
		return nil, err
	}
	return gc.Wait4PlayerIdState, nil
}

func (gc *GameClient) Wait4PlayerIdState() (fsm.StateFunc, error) {
	gc.log("Waiting for player ID from server")
	var ahoyMsg AhoyMsg
	// if there are no bytes available, wait for 1 second and try again

	decoder := gob.NewDecoder(gc.serverReader)
	err := decoder.Decode(&ahoyMsg)
	if err != nil {
		gc.log("Error decoding client ID message: %s", err)
		return nil, err
	}

	if ahoyMsg.ErrMsg != "" {
		// some failure occured
		gc.log("Got error from server: %s", ahoyMsg.ErrMsg)
		return nil, errors.New(ahoyMsg.ErrMsg)
	}

	gc.playerId = ahoyMsg.UserID
	return gc.SendPingState, nil
}

func (gc *GameClient) SendPingState() (fsm.StateFunc, error) {
	// send a ping to the server
	gc.log("Sending ping to server")
	pingMsg := PingMsg{gc.playerId}
	encoder := gob.NewEncoder(gc.serverWriter)
	err := encoder.Encode(pingMsg)
	if err != nil {
		// if broken pipe, reconnect to server
		if err == io.ErrClosedPipe {
			return gc.ConnectToServerState, nil
		}
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
		switch {
		case errors.Is(err, syscall.ECONNRESET):
			gc.log("Connection was reset. Trying to reconnect...")
			return gc.ResetServerConnectionState, nil
		}
		gc.log("Error waiting for pong: %s", err)
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
