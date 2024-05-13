package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

//go:embed index.html
var index string

var upgrader = websocket.Upgrader{}

func main() {
	homeTemplate := template.Must(template.New("").Parse(index))

	// UI page
	http.HandleFunc("GET /{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		homeTemplate.Execute(w, topic)
	})

	// actual websocket
	http.HandleFunc("GET /.ws/{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		serveWebsocket(topic, w, r)
	})

	// running the server
	log.Println("Starting server at localhost:8888")
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		log.Fatal(err)
	}
}
