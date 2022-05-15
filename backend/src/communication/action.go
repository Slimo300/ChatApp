package communication

import (
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gorilla/websocket"
)

const TIME_FORMAT = "2006-02-01 15:04:05"

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
	Invite models.Invite `json:"invite"` // invite
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

func ShortenMessages(messages []models.Message) (shortMessages []Message) {
	for _, msg := range messages {
		shortMessages = append(shortMessages, Message{
			Group:   uint64(msg.Member.GroupID),
			Member:  uint64(msg.MemberID),
			Nick:    msg.Member.Nick,
			Message: msg.Text,
			When:    msg.Posted.Format(TIME_FORMAT),
		})
	}
	return
}
