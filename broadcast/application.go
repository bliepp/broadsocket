package broadcast

import (
	"log"
	"net/http"
	"strings"
)

type Application struct {
	router *http.ServeMux
}

func New() *Application {
	a := &Application{
		router: http.NewServeMux(),
	}

	// UI page
	a.router.HandleFunc("GET /{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		UITemplate.Execute(w, topic)
	})

	// actual websocket
	a.router.HandleFunc("GET /.ws/{topic...}", func(w http.ResponseWriter, r *http.Request) {
		topic := "/" + strings.TrimRight(r.PathValue("topic"), "/")
		ServeWebsocketClient(topic, w, r)
	})

	return a
}

func (a *Application) ListenAndServe(addr string) error {
	// welcome message(s)
	log.Printf("Starting broadsocket on %s", addr)

	// running the server
	return http.ListenAndServe(addr, a.router)
}
