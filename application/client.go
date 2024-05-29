package application

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	broadcast *Broadcast
	conn      *websocket.Conn
	send      chan []byte
}

func NewClient(b *Broadcast, c *websocket.Conn) *Client {
	return &Client{
		broadcast: b,
		conn:      c,
		send:      make(chan []byte, 256),
	}
}

func (c *Client) Register() {
	c.broadcast.register <- c
}

func (c *Client) Unregister() {
	c.broadcast.unregister <- c
	c.conn.Close()
}
