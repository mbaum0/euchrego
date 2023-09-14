package comms

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"net"
	"time"
)

type MsgKind int

const (
	HelloMsgKind MsgKind = iota
	AhoyMsgKind
	PingMsgKind
	PongMsgKind
	LeaveMsgKind
)

type Msg struct {
	Kind MsgKind
	Data interface{}
}

// Used when a new user connects to the server
type HelloMsg struct {
	UserName string
	UserID   string // may be empty if user doesn't have one yet
}

// AhoyMsg is sent from the server to the client in response to a HelloMsg. It provides
// the user with a unique ID.
// If the server can't accept the client, success is set to false.
type AhoyMsg struct {
	UserID string
	ErrMsg string
}

type PingMsg struct {
	UserID string
}

type PongMsg struct {
	UserID string
}

type LeaveMsg struct {
	UserID string
}

func generateRandomID(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes using base64
	randomID := base64.RawURLEncoding.EncodeToString(randomBytes)

	return randomID, nil
}

type Sock struct {
	Reader *gob.Decoder
	Writer *gob.Encoder
	Conn   net.Conn
}

func NewSock(conn net.Conn) *Sock {
	gob.Register(HelloMsg{})
	gob.Register(AhoyMsg{})
	gob.Register(PingMsg{})
	gob.Register(PongMsg{})
	gob.Register(LeaveMsg{})
	handler := Sock{}
	handler.Reader = gob.NewDecoder(conn)
	handler.Writer = gob.NewEncoder(conn)
	handler.Conn = conn
	return &handler
}

func (s *Sock) Disconnect() {
	if s.Conn != nil {
		s.Conn.Close()
	}
	s.Conn = nil
	s.Reader = nil
	s.Writer = nil
}

func (h *Sock) SendHelloMsg(UserName string, UserID string) error {
	msg := Msg{Kind: HelloMsgKind, Data: HelloMsg{UserName, UserID}}
	return h.Writer.Encode(msg)
}

func (h *Sock) SendAhoyMsg(UserID string, ErrMsg string) error {
	msg := Msg{Kind: AhoyMsgKind, Data: AhoyMsg{UserID, ErrMsg}}
	return h.Writer.Encode(msg)
}

func (h *Sock) SendPingMsg(UserID string) error {
	msg := Msg{Kind: PingMsgKind, Data: PingMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *Sock) SendPongMsg(UserID string) error {
	msg := Msg{Kind: PongMsgKind, Data: PongMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *Sock) SendLeaveMsg(UserID string) error {
	msg := Msg{Kind: LeaveMsgKind, Data: LeaveMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *Sock) ReadAhoyMsg() (AhoyMsg, error) {
	var ahoyMsg AhoyMsg
	data, kind, err := h.ReadMsg()

	if err != nil {
		return ahoyMsg, errors.New("failed to parse ahoy msg")
	}

	if kind != AhoyMsgKind {
		return ahoyMsg, errors.New("msg is incorrect type")
	}

	return data.(AhoyMsg), nil
}

func (h *Sock) ReadHelloMsg() (HelloMsg, error) {
	var helloMsg HelloMsg
	data, kind, err := h.ReadMsg()

	if err != nil {
		return helloMsg, errors.New("failed to parse hello msg")
	}

	if kind != HelloMsgKind {
		return helloMsg, errors.New("msg is incorrect type")
	}

	return data.(HelloMsg), nil
}

func (h *Sock) ReadPingMsg() (PingMsg, error) {
	var pingMsg PingMsg
	data, kind, err := h.ReadMsg()

	if err != nil {
		return pingMsg, errors.New("failed to parse ping msg")
	}

	if kind != PingMsgKind {
		return pingMsg, errors.New("msg is incorrect type")
	}

	return data.(PingMsg), nil
}

func (h *Sock) ReadPongMsg() (PongMsg, error) {
	var pongMsg PongMsg
	data, kind, err := h.ReadMsg()

	if err != nil {
		return pongMsg, errors.New("failed to parse pong msg")
	}

	if kind != PongMsgKind {
		return pongMsg, errors.New("msg is incorrect type")
	}

	return data.(PongMsg), nil
}

func (h *Sock) ReadLeaveMsg() (LeaveMsg, error) {
	var leaveMsg LeaveMsg
	data, kind, err := h.ReadMsg()

	if err != nil {
		return leaveMsg, errors.New("failed to parse disconnect msg")
	}

	if kind != LeaveMsgKind {
		return leaveMsg, errors.New("msg is incorrect type")
	}

	return data.(LeaveMsg), nil
}

func (h *Sock) ReadMsg() (interface{}, MsgKind, error) {
	timeoutDuration := 5 * time.Second
	h.Conn.SetDeadline(time.Now().Add(timeoutDuration))
	var msg Msg
	err := h.Reader.Decode(&msg)
	return msg.Data, msg.Kind, err
}
