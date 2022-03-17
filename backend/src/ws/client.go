package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan *Message
	hub    *Hub
	groups []int64
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		c.hub.forward <- &msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
