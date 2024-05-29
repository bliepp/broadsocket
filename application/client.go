package application

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	broadcast *Broadcast
	conn      *websocket.Conn
	send      chan []byte
}

func NewClient(b *Broadcast, conn *websocket.Conn) *Client {
	return &Client{
		broadcast: b,
		conn:      conn,
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
