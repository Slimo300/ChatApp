package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	forward chan *Message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func NewHub() *Hub {
	return &Hub{
		forward: make(chan *Message),
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
			for client := range h.clients {
				for gr := range client.groups {
					if gr == int(msg.Group) {
						client.send <- msg
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

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func ServeWebSocket(w http.ResponseWriter, req *http.Request, h *Hub, groups []int64) {

	// Upgrading connection to WebSockets
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// Creating a client object to represent a user in a room
	client := &client{
		socket: socket,
		send:   make(chan *Message, messageBufferSize),
		hub:    h,
		groups: groups,
	}
	h.join <- client
	defer func() { h.leave <- client }()
	go client.write()
	go client.read()
}
