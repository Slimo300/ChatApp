package ws

import (
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/gorilla/websocket"
)

type client struct {
	id     int
	socket *websocket.Conn
	send   chan *database.Message
	hub    *Hub
	groups []int64
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg database.Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
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
