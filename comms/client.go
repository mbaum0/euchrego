package comms

import (
	"errors"
	"net"
	"sync"
	"syscall"
	"time"
)

type GameClient struct {
	playerId           string // given to us by the server
	playerName         string // set by the player
	serverHost         string
	serverPort         string
	lastPingSendTime   time.Time
	connectionAttempts int
	*ConnMgr
	sync.Mutex
}

func NewGameClient(playerName string, serverHost string, serverPort string) *GameClient {

	gc := GameClient{}
	gc.playerName = playerName
	gc.serverHost = serverHost
	gc.serverPort = serverPort
	gc.lastPingSendTime = time.Now()
	gc.connectionAttempts = 0
	gc.Mutex = sync.Mutex{}
	gc.ConnMgr = &ConnMgr{}
	return &gc

}

func (gc *GameClient) ReconnectToServer() error {
	return gc.ConnectToServer(gc.playerId)
}

func (gc *GameClient) ConnectToServer(userID string) error {
	if gc.IsConnected() {
		logError("Already connected to a server")
		return nil
	}
	connStr := gc.serverHost + ":" + gc.serverPort
	logInfo("Connecting to server at %s", connStr)
	server, err := net.Dial("tcp", connStr)
	if err != nil {
		switch {
		case errors.Is(err, syscall.ECONNREFUSED):
			logError("Connected refused. Is server up?")
			return err
		default:
			logError("Error connecting to server: %s", err)
			return err
		}
	}
	gc.ConnMgr = NewConnMgr(server)
	logSuccess("Successfully connected to server")

	logInfo("Sending hello message to server")
	err = gc.SendHelloMsg(gc.playerName, userID)
	if err != nil {
		logError("Error sending hello message to server: %s", err)
		gc.DisconnectFromServer()
		return err
	}

	logInfo("Waiting for player ID from server")
	ahoy, err := gc.ReadAhoyMsg()
	if err != nil {
		logError("Error decoding client ID message: %s", err)
		gc.DisconnectFromServer()
		return err
	}

	if ahoy.ErrMsg != "" {
		// some failure occured
		logError("Got error from server: %s", ahoy.ErrMsg)
		gc.DisconnectFromServer()
		return errors.New(ahoy.ErrMsg)
	}

	gc.playerId = ahoy.UserID
	logSuccess("Successfully obtained ID: %s", ahoy.UserID)
	go gc.WatchConnection()
	return nil
}

func (gc *GameClient) DisconnectFromServer() {
	if !gc.IsConnected() {
		logError("Can't disconnect. No connection is present.")
		return
	}
	gc.Disconnect()
}

func (gc *GameClient) LeaveServer() {
	if !gc.IsConnected() {
		logError("Cant't leave. No connection is present.")
		return
	}
	gc.SendLeaveMsg(gc.playerId)
	gc.DisconnectFromServer()
}

func (gc *GameClient) WatchConnection() {
	for {
		// don't want to ping if there isn't a connection
		if !gc.IsConnected() {
			break
		}
		err := gc.SendPingMsg(gc.playerId)
		if err != nil {
			logError("ping was unsuccessful. Disconnecting from server.")
			gc.DisconnectFromServer()
			break
		}

		_, err = gc.ReadPongMsg()
		if err != nil {
			logError("pong was unsuccessful. Disconnecting from server. %s", err)
			gc.DisconnectFromServer()
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (gc *GameClient) SendPing() error {
	// send a ping to the server
	logInfo("Sending ping to server")
	err := gc.SendPingMsg(gc.playerId)
	if err != nil {
		logError("Error sending ping to server: %s", err)
		return err
	}
	logInfo("Waiting for pong from server")
	_, err = gc.ReadPongMsg()
	if err != nil {
		logError("Error waiting for pong: %s", err)
		return err
	}

	logSuccess("Got pong!")

	return nil
}
