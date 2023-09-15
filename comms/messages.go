package comms

import (
	"errors"
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

// The base message type
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

// LeaveMsg is sent when a user wishes to disconnect from the server
type LeaveMsg struct {
	UserID string
}

func (h *ConnMgr) SendHelloMsg(UserName string, UserID string) error {
	msg := Msg{Kind: HelloMsgKind, Data: HelloMsg{UserName, UserID}}
	return h.Writer.Encode(msg)
}

func (h *ConnMgr) SendAhoyMsg(UserID string, ErrMsg string) error {
	msg := Msg{Kind: AhoyMsgKind, Data: AhoyMsg{UserID, ErrMsg}}
	return h.Writer.Encode(msg)
}

func (h *ConnMgr) SendPingMsg(UserID string) error {
	msg := Msg{Kind: PingMsgKind, Data: PingMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *ConnMgr) SendPongMsg(UserID string) error {
	msg := Msg{Kind: PongMsgKind, Data: PongMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *ConnMgr) SendLeaveMsg(UserID string) error {
	msg := Msg{Kind: LeaveMsgKind, Data: LeaveMsg{UserID}}
	return h.Writer.Encode(msg)
}

func (h *ConnMgr) ReadAhoyMsg() (AhoyMsg, error) {
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

func (h *ConnMgr) ReadHelloMsg() (HelloMsg, error) {
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

func (h *ConnMgr) ReadPingMsg() (PingMsg, error) {
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

func (h *ConnMgr) ReadPongMsg() (PongMsg, error) {
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

func (h *ConnMgr) ReadLeaveMsg() (LeaveMsg, error) {
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

func (h *ConnMgr) ReadMsg() (interface{}, MsgKind, error) {
	timeoutDuration := 5 * time.Second
	h.conn.SetDeadline(time.Now().Add(timeoutDuration))
	var msg Msg
	err := h.Reader.Decode(&msg)
	return msg.Data, msg.Kind, err
}
