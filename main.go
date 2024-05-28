package main

import (
	"flag"
	"log"

	"broadsocket/broadcast"
)

func main() {
	// adjustable host/port
	bind_addr := flag.String("b", "localhost:8888", "host and port for serving broadsocket")
	flag.Parse()

	// running the server
	a := broadcast.New()
	if err := a.ListenAndServe(*bind_addr); err != nil {
		log.Fatal(err)
	}
}
