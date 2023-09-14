package comms

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

type GameClient struct {
	playerId           string // given to us by the server
	playerName         string // set by the player
	serverHost         string
	serverPort         string
	Sock               *Sock
	lastPingSendTime   time.Time
	connectionAttempts int
	sync.Mutex
}

func NewGameClient(playerName string, serverHost string, serverPort string) *GameClient {

	gc := GameClient{}
	gc.playerName = playerName
	gc.serverHost = serverHost
	gc.serverPort = serverPort
	gc.lastPingSendTime = time.Now()
	gc.connectionAttempts = 0
	gc.Sock = nil
	gc.Mutex = sync.Mutex{}
	return &gc

}

func (gc *GameClient) logError(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgRed).Printf(format, args...)
}

func (gc *GameClient) logInfo(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgBlue).Printf(format, args...)
}

func (gc *GameClient) logSuccess(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgGreen).Printf(format, args...)
}

func (gc *GameClient) ReconnectToServer() error {
	return gc.ConnectToServer(gc.playerId)
}

func (gc *GameClient) ConnectToServer(userID string) error {
	if gc.Sock != nil {
		gc.logError("Already connected to a server")
		return nil
	}
	connStr := gc.serverHost + ":" + gc.serverPort
	gc.logInfo("Connecting to server at %s", connStr)
	server, err := net.Dial("tcp", connStr)
	if err != nil {
		switch {
		case errors.Is(err, syscall.ECONNREFUSED):
			gc.logError("Connected refused. Is server up?")
			return err
		default:
			gc.logError("Error connecting to server: %s", err)
			return err
		}
	}
	gc.Sock = NewSock(server)
	gc.logSuccess("Successfully connected to server")

	gc.logInfo("Sending hello message to server")
	err = gc.Sock.SendHelloMsg(gc.playerName, userID)
	if err != nil {
		gc.logError("Error sending hello message to server: %s", err)
		gc.DisconnectFromServer()
		return err
	}

	gc.logInfo("Waiting for player ID from server")
	ahoy, err := gc.Sock.ReadAhoyMsg()
	if err != nil {
		gc.logError("Error decoding client ID message: %s", err)
		gc.DisconnectFromServer()
		return err
	}

	if ahoy.ErrMsg != "" {
		// some failure occured
		gc.logError("Got error from server: %s", ahoy.ErrMsg)
		gc.DisconnectFromServer()
		return errors.New(ahoy.ErrMsg)
	}

	gc.playerId = ahoy.UserID
	gc.logSuccess("Successfully obtained ID: %s", ahoy.UserID)
	go gc.WatchConnection()
	return nil
}

func (gc *GameClient) DisconnectFromServer() {
	if gc.Sock == nil {
		gc.logError("Can't disconnect. No connection is present.")
		return
	}
	gc.Sock.Disconnect()
	gc.Sock = nil
}

func (gc *GameClient) LeaveServer() {
	if gc.Sock == nil {
		gc.logError("Cant't leave. No connection is present.")
		return
	}
	gc.Sock.SendLeaveMsg(gc.playerId)
	gc.DisconnectFromServer()
}

func (gc *GameClient) WatchConnection() {
	for {
		// don't want to ping if there isn't a connection
		if gc.Sock == nil {
			break
		}
		err := gc.Sock.SendPingMsg(gc.playerId)
		if err != nil {
			gc.logError("ping was unsuccessful. Disconnecting from server.")
			gc.DisconnectFromServer()
			break
		}

		_, err = gc.Sock.ReadPongMsg()
		if err != nil {
			gc.logError("pong was unsuccessful. Disconnecting from server. %s", err)
			gc.DisconnectFromServer()
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (gc *GameClient) SendPing() error {
	// send a ping to the server
	gc.logInfo("Sending ping to server")
	err := gc.Sock.SendPingMsg(gc.playerId)
	if err != nil {
		gc.logError("Error sending ping to server: %s", err)
		return err
	}
	gc.logInfo("Waiting for pong from server")
	_, err = gc.Sock.ReadPongMsg()
	if err != nil {
		gc.logError("Error waiting for pong: %s", err)
		return err
	}

	gc.logSuccess("Got pong!")

	return nil
}
