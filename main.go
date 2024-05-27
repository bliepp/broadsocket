package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"broadsocket/broadcast"
)

func main() {
	// adjustable port
	bind_addr := flag.String("b", "localhost:8888", "host and port for serving broadsocket")
	flag.Parse()

	// UI page
	http.HandleFunc("GET /{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		broadcast.UITemplate.Execute(w, topic)
	})

	// actual websocket
	http.HandleFunc("GET /.ws/{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		broadcast.ServeWebsocketClient(topic, w, r)
	})

	// running the server
	log.Printf("Starting server at %s", *bind_addr)
	if err := http.ListenAndServe(*bind_addr, nil); err != nil {
		log.Fatal(err)
	}
}
