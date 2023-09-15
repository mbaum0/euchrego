package comms

import (
	"encoding/gob"
	"net"
)

type ConnMgr struct {
	Reader *gob.Decoder
	Writer *gob.Encoder
	conn   net.Conn
}

func NewConnMgr(conn net.Conn) *ConnMgr {
	gob.Register(HelloMsg{})
	gob.Register(AhoyMsg{})
	gob.Register(PingMsg{})
	gob.Register(PongMsg{})
	gob.Register(LeaveMsg{})
	mgr := ConnMgr{}
	mgr.Reader = gob.NewDecoder(conn)
	mgr.Writer = gob.NewEncoder(conn)
	mgr.conn = conn
	return &mgr
}

func (s *ConnMgr) IsConnected() bool {
	return s.conn != nil
}

func (s *ConnMgr) Disconnect() {
	if s.conn != nil {
		s.conn.Close()
	}
	s.conn = nil
	s.Reader = nil
	s.Writer = nil
}
