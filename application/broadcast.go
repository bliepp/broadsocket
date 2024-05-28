package application

import (
	"github.com/gorilla/websocket"
)

type Broadcast struct {
	topic      string
	clients    map[*Client]bool
	messages   chan []byte
	register   chan *Client
	unregister chan *Client
}

// map of all broadcasts (key is the topic)
var broadcasts map[string]*Broadcast = make(map[string]*Broadcast)
var upgrader = websocket.Upgrader{}

func NewBroadcast(topic string) *Broadcast {
	// reuse existing broadcasts on the same topic
	b, ok := broadcasts[topic]
	if ok {
		return b
	}

	// if there is no matching boradcast for given topic,
	// create a new one and start the goroutine
	b = &Broadcast{
		topic:      topic,
		clients:    make(map[*Client]bool),
		messages:   make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	broadcasts[topic] = b
	go b.Run()

	return b
}

// create client that is attached to that broadcast
func (b *Broadcast) NewClient(c *websocket.Conn) *Client {
	return NewClient(b, c)
}

func (b *Broadcast) Run() {
	done := make(chan bool)

	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
		case client := <-b.unregister:
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.send)
			}

			if len(b.clients) == 0 {
				delete(broadcasts, b.topic)
				done <- true // break the goroutine
			}
		case msg := <-b.messages:
			for client := range b.clients {
				select {
				case client.send <- msg:
				default:
					delete(b.clients, client)
					close(client.send)
				}
			}
		case <-done:
			return
		}
	}
}
