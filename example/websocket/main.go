package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/dev-cloverlab/go-message-broadcaster"
	"golang.org/x/net/websocket"
)

type Conn struct {
	ws *websocket.Conn
}

func (c Conn) SendMessage(msg *broadcaster.ResponseMessage) error {
	return websocket.JSON.Send(c.ws, msg)
}

func (c Conn) SendEvent(e *broadcaster.EventMessage) error {
	return websocket.JSON.Send(c.ws, e)
}

func (c Conn) Receive() (*broadcaster.RequestMessage, error) {
	msg := &broadcaster.RequestMessage{}
	err := websocket.JSON.Receive(c.ws, msg)
	return msg, err
}

func (c *Conn) Close() error {
	return c.ws.Close()
}

func main() {
	port := flag.Int("p", 9218, "websocket listen port")
	endpoint := flag.String("e", "/", "websocket application endpoint path")
	flag.Parse()

	// Create message handlers for each message
	handlers := broadcaster.MessageHandlers{
		1: func(msg *broadcaster.RequestMessage, c context.Context) (*broadcaster.ResponseMessage, error) {
			// Create broadcasting message object (This is simple echo handler).
			return broadcaster.NewResponseMessage(broadcaster.All, msg.Body), nil
		},
	}

	// Initialize the server and listen all clients
	sv := broadcaster.NewServer(context.Background(), handlers)
	go sv.Listen()

	http.Handle(*endpoint, websocket.Handler(func(ws *websocket.Conn) {
		// Add new client to server and listen ws connection
		sv.NewClient(&Conn{ws: ws}).Listen()
	}))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
