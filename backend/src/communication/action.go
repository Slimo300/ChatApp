package communication

import (
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gorilla/websocket"
)

// Sender interface decides which struct can be sent via websocket connection
type Sender interface {
	Send(*websocket.Conn) error
}

// Action represents type for signalizing changes to hub and further to frontend
type Action struct {
	Action string        `json:"action"` // pop or insert
	Group  int           `json:"group"`  // id_group
	User   int           `json:"-"`      // id_user
	Member models.Member `json:"member"` // member info for updates
}

// Send sends itself through websocket connection
func (a *Action) Send(ws *websocket.Conn) error {
	if err := ws.WriteJSON(a); err != nil {
		return err
	}
	return nil
}

// Message is a plain message in chat app
type Message struct {
	Group   uint64 `json:"group"`
	Member  uint64 `json:"member"`
	Message string `json:"text"`
	Nick    string `json:"nick"`
	When    string `json:"created"`
}

// Send sends itself through websocket connection
func (m *Message) Send(ws *websocket.Conn) error {
	if err := ws.WriteJSON(m); err != nil {
		return err
	}
	return nil
}
