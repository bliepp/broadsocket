package main

import (
	"flag"

	"github.com/bliepp/broadsocket/application"
)

func main() {
	// adjustable host/port
	bind_addr := flag.String("b", "localhost:8888", "host and port for serving broadsocket")
	flag.Parse()

	// running the server
	a := application.New()
	if err := a.ListenAndServe(*bind_addr); err != nil {
		a.Logger.Fatal(err)
	}
}
