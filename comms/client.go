package comms

import (
	"encoding/gob"
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
	serverReader       *gob.Decoder
	serverWriter       *gob.Encoder
	server             net.Conn
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
	if gc.server != nil {
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
	gc.serverReader = gob.NewDecoder(server)
	gc.serverWriter = gob.NewEncoder(server)
	gc.server = server
	gc.logSuccess("Successfully connected to server")

	gc.logInfo("Sending hello message to server")
	hello := HelloMsg{gc.playerName, userID}
	err = gc.serverWriter.Encode(hello)
	if err != nil {
		gc.logError("Error sending hello message to server: %s", err)
		return err
	}

	gc.logInfo("Waiting for player ID from server")
	var ahoyMsg AhoyMsg
	err = gc.serverReader.Decode(&ahoyMsg)
	if err != nil {
		gc.logError("Error decoding client ID message: %s", err)
		return err
	}

	if ahoyMsg.ErrMsg != "" {
		// some failure occured
		gc.logError("Got error from server: %s", ahoyMsg.ErrMsg)
		return errors.New(ahoyMsg.ErrMsg)
	}

	gc.playerId = ahoyMsg.UserID
	gc.logSuccess("Successfully obtained ID: %s", ahoyMsg.UserID)
	go gc.WatchConnection()
	return nil
}

func (gc *GameClient) DisconnectFromServer() {
	if gc.server == nil {
		gc.logError("Can't disconnect. No connection is present.")
		return
	}
	gc.server.Close()
	gc.server = nil
	gc.serverWriter = nil
	gc.serverReader = nil
}

func (gc *GameClient) WatchConnection() {
	pingMsg := PingMsg{gc.playerId}
	var pongMsg PongMsg

	for {
		// don't want to ping if there isn't a connection
		if gc.server == nil {
			return
		}
		err := gc.serverWriter.Encode(pingMsg)
		if err != nil {
			gc.logError("ping was unsuccessful. Disconnecting from server.")
			gc.DisconnectFromServer()
			break
		}

		err = gc.serverReader.Decode(&pongMsg)
		if err != nil {
			gc.logError("pong was unsuccessful. Disconnecting from server. %s", err)
			gc.DisconnectFromServer()
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func (gc *GameClient) SendPing() error {
	// send a ping to the server
	gc.logInfo("Sending ping to server")
	pingMsg := PingMsg{gc.playerId}
	err := gc.serverWriter.Encode(pingMsg)
	if err != nil {
		gc.logError("Error sending ping to server: %s", err)
		return err
	}
	gc.logInfo("Waiting for pong from server")
	var pongMsg PongMsg
	err = gc.serverReader.Decode(&pongMsg)
	if err != nil {
		gc.logError("Error waiting for pong: %s", err)
		return err
	}

	gc.logSuccess("Got pong!")

	return nil
}
