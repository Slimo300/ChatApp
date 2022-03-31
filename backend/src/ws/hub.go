package ws

import (
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/gorilla/websocket"
)

type Hub struct {
	db      database.DBlayer
	dbconn  <-chan *communication.Action
	forward chan *communication.Message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func NewHub(db database.DBlayer, ch <-chan *communication.Action) *Hub {
	return &Hub{
		db:      db,
		dbconn:  ch,
		forward: make(chan *communication.Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
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
			message, err := h.db.AddMessage(*msg)
			if err != nil {
				panic(err)
			}
			message.Nick = msg.Nick
			for client := range h.clients {
				for _, gr := range client.groups {
					if gr == int64(msg.Group) {
						client.send <- &message
					}
				}
			}
		case msg := <-h.dbconn:
			switch msg.Action {
			case "pop":
				for client := range h.clients {
					if client.id == msg.User || msg.User == 0 {
						for i, gr := range client.groups {
							if gr == int64(msg.Group) {
								client.groups = append(client.groups[:i], client.groups[:i+1]...)
							}
						}
					}
				}
			case "insert":
				for client := range h.clients {
					if client.id == msg.User || msg.User == 0 {
						client.groups = append(client.groups, int64(msg.Group))
					}
				}
			case "add":
				for client := range h.clients {
					for _, gr := range client.groups {
						if gr == int64(msg.Group) {
							client.send <- &communication.Message{}
						}
					}
				}
			}
		}
	}
}

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
