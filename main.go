package main

import (
	"log"
	"net/http"
	"strings"

	"broadsocket/broadcast"
)

func main() {
	// UI page
	http.HandleFunc("GET /{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		broadcast.UITemplate.Execute(w, topic)
	})

	// actual websocket
	http.HandleFunc("GET /.ws/{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		broadcast.ServeWebsocket(topic, w, r)
	})

	// running the server
	log.Println("Starting server at localhost:8888")
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		log.Fatal(err)
	}
}
