package application

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
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
		ServeWebsocketClient(topic, w, r)
	})

	return a
}

func (a *Application) ListenAndServe(addr string) error {
	// url form of bind addr
	webuiAddr := url.URL{Scheme: "http", Host: addr, Path: "/"}
	if webuiAddr.Hostname() == "" {
		webuiAddr.Host = "0.0.0.0" + webuiAddr.Host
	}
	wsAddr := webuiAddr
	wsAddr.Path = ".ws/"

	// welcome message(s)
	boldPrintf("Welcome to %s\n\n", color.RedString("Broadsocket"))
	bluePrintf("\tWeb UI:    %s\n", webuiAddr.String())
	bluePrintf("\tWebsocket: %s\n\n", wsAddr.String())

	// try writing to an additional log file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
