package broadcast

import (
	"bytes"
	"log"
	"net/http"
	"time"

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

func ServeWebsocketClient(topic string, w http.ResponseWriter, r *http.Request) {
	b := NewBroadcast(topic)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error establishing websocket connection: %v", err)
		return
	}

	client := b.NewClient(conn)
	client.Register()

	// getting messages from client and pushing them to the broadcast
	go func() {
		defer client.Unregister()

		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error: %v", err)
				}
				return
			}

			client.broadcast.messages <- bytes.TrimSpace(bytes.Replace(msg, []byte("\n"), []byte(" "), -1))
		}
	}()

	// getting messages from the broadcast and pushing it to the current client
	go func() {
		defer client.Unregister()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-client.send:
				if !ok {
					client.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := client.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(msg)

				// Add queued chat messages to the current websocket message.
				n := len(client.send)
				for i := 0; i < n; i++ {
					w.Write([]byte("\n"))
					w.Write(<-client.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			case <-ticker.C:
				if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

}
