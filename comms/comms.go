package comms

import (
	"encoding/gob"
	"io"

	"github.com/mbaum0/euchrego/godeck"
)

type CardSelectMsg struct {
	UserID string
	Suit   godeck.Suit
	Rank   godeck.Rank
}

type BoolMsg struct {
	UserID string
	Value  bool
}

type Client struct {
	ID            string
	Reader        io.Reader
	Writer        io.Writer
	cardSelectMsg *CardSelectMsg
	boolMsg       *BoolMsg
}

type CommsManager struct {
	clients map[string]Client
}

func (m *CommsManager) AddClient(c Client) {
	m.clients[c.ID] = c
}

func (m *CommsManager) RemoveClient(id string) {
	delete(m.clients, id)
}

func (m *CommsManager) GetBoolMsgForClient(id string) *BoolMsg {
	return m.clients[id].boolMsg
}

func (m *CommsManager) GetCardSelectMsgForClient(id string) *CardSelectMsg {
	return m.clients[id].cardSelectMsg
}

// PollForInput reads over each client's Reader, parses the input and sets the appropriate message type
func (m *CommsManager) PollForInput() {
	gob.Register(CardSelectMsg{})
	gob.Register(BoolMsg{})

	for _, c := range m.clients {
		var msg interface{}
		dec := gob.NewDecoder(c.Reader)
		err := dec.Decode(&msg)
		if err != nil {
			continue
		}

		switch msg.(type) {
		case CardSelectMsg:
			c.cardSelectMsg = msg.(*CardSelectMsg)
		case BoolMsg:
			c.boolMsg = msg.(*BoolMsg)
		}
	}
}
