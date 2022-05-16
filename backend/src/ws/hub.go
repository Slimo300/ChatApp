package ws

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

type Hub struct {
	actionServerChan  <-chan *communication.Action
	messageServerChan chan<- *communication.Message
	forward           chan *communication.Message
	join              chan *client
	leave             chan *client
	clients           map[*client]bool
}

func NewHub(messageChan chan<- *communication.Message, actionChan <-chan *communication.Action) *Hub {
	return &Hub{
		actionServerChan:  actionChan,
		messageServerChan: messageChan,
		forward:           make(chan *communication.Message),
		join:              make(chan *client),
		leave:             make(chan *client),
		clients:           make(map[*client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.leave:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.forward:
			msg.SetTime()
			h.messageServerChan <- msg
			for client := range h.clients {
				for _, gr := range client.groups {
					if gr == int64(msg.Group) {
						client.send <- msg
					}
				}
			}
		case msg := <-h.actionServerChan:
			switch msg.Action {
			case "DELETE_GROUP":
				h.GroupDeleted(msg.Group)
			case "CREATE_GROUP":
				h.GroupCreated(msg.User, msg.Group)
			case "ADD_MEMBER":
				h.MemberAdded(msg.Member)
			case "DELETE_MEMBER":
				h.MemberDeleted(msg.Member)
			case "SEND_INVITE":
				h.SendGroupInvite(msg.Invite)
			}
		}
	}
}

func ServeWebSocket(w http.ResponseWriter, req *http.Request, h *Hub, groups []int64, id_user int) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	client := &client{
		id:     id_user,
		socket: socket,
		send:   make(chan communication.Sender, messageBufferSize),
		hub:    h,
		groups: groups,
	}

	h.join <- client
	defer func() { h.leave <- client }()
	go client.write()
	client.read()
}
