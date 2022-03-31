package ws

import (
	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/gorilla/websocket"
)

type client struct {
	id     int
	socket *websocket.Conn
	send   chan communication.Sender
	hub    *Hub
	groups []int64
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		// socket can read only communication message
		var msg communication.Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		c.hub.forward <- &msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := msg.Send(c.socket); err != nil {
			break
		}
	}
}
