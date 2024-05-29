package application

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var boldPrintf = color.New(color.Bold).PrintfFunc()
var bluePrintf = color.New(color.FgBlue).PrintfFunc()

type Application struct {
	router *http.ServeMux
	Logger *log.Logger
}

func New() *Application {
	// app instance
	a := &Application{
		router: http.NewServeMux(),
		Logger: log.New(os.Stderr, "| ", log.LstdFlags|log.Lmsgprefix),
	}

	// UI page
	a.router.HandleFunc("GET /{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		UITemplate.Execute(w, topic)
	})

	// actual websocket
	a.router.HandleFunc("GET /.ws/{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		a.handleWebsocket(topic, w, r)
	})

	return a
}

func (a *Application) handleWebsocket(topic string, w http.ResponseWriter, r *http.Request) {
	b := NewBroadcast(topic)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.Logger.Printf("Error establishing websocket connection: %v", err)
		return
	}

	client := b.NewClient(conn)

	// getting messages from client and pushing them to the broadcast
	go func() {
		defer client.Unregister()

		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					a.Logger.Printf("Error: %v", err)
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

func (a *Application) ListenAndServe(addr string) error {
	// url form of bind addr
	webAddr := addr
	if strings.HasPrefix(webAddr, ":") {
		webAddr = "0.0.0.0" + webAddr
	}

	// welcome message(s)
	boldPrintf("Welcome to %s\n\n", color.RedString("Broadsocket"))
	bluePrintf("\tWeb UI:    http://%s/\n", webAddr)
	bluePrintf("\tWebsocket: http://%s/.ws/\n\n", webAddr)

	// try writing to an additional log file
	logFile, err := os.OpenFile("broadsocket.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		defer logFile.Close()
		a.Logger.SetOutput(io.MultiWriter(a.Logger.Writer(), logFile))
	} else {
		a.Logger.Printf("Could not open log file, err \"%v\"", err)
	}

	// ready message
	a.Logger.Printf("Ready on %s", addr)

	// running the server
	return http.ListenAndServe(addr, a.router)
}
