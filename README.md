# go-message-broadcaster

[![GitHub release](https://img.shields.io/github/release/dev-cloverlab/go-message-broadcaster.svg?style=flat-square)](https://github.com/dev-cloverlab/go-message-broadcaster)
[![license](https://img.shields.io/github/license/dev-cloverlab/go-message-broadcaster.svg?style=flat-square)](https://github.com/dev-cloverlab/go-message-broadcaster)

go-message-broadcaster is bidirectional message communication middleware that composed on server-client architecture.  
This can be used on any connection interfaces as you like.  

# WebSocket example

This is message broadcasting example using WebSocket connection.  
Following code is a quote from example/websocket/main.go:

```go
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
```

Start the server:

```
% go run example/websocket/main.go
```

Connect to the server from two or more terminal client using `wscat`:

```
% wscat -c ws://localhost:9218/ -o localhost:9218
connected (press CTRL+C to quit)
< {"EventType":1,"ClientID":"e68e76f0-7725-3f76-d9fb-bb4a72ca2121"}
> 
```

Send messages from each others:

```
> {"HandlerID": 1, "Body":"SGVsbG8sIHdvcmxkIQ=="}
< {"SenderID":"98f1088e-1e28-3cec-c73a-ab7e4c678149","HandlerID":1,"CastType":1,"CastFor":null,"Body":"SGVsbG8sIHdvcmxkIQ=="}
```

# Connection interface

For using go-message-broadcaster, you need to implement the following Conn interface:

```go
type Conn interface {
	SendMessage(*ResponseMessage) error
	SendEvent(*EventMessage) error
	Receive() (*RequestMessage, error)
	Close() error
}
```

See the example: example/websocket/main.go

# Message and Event

go-message-broadcaster has three message types `RequestMessage`, `ResponseMessage` and `EventMessage`.  

## RequestMessage 

```go
type RequestMessage struct {
	SenderID  ClientID
	HandlerID MessageHandlerID
	Body      []byte
}
```

RequestMessage is the type when sending from client to server.  

- `SenderID` is the client ID that was granted by the server.  
- `HandlerID` is message handler ID of the process that corresponding for each message as you defined.  
- `Body` is the request body

You don't need to set the `SenderID`, this is granted by the system automatically.   

## ResponseMessage

```go
type ResponseMessage struct {
	SenderID  ClientID
	HandlerID MessageHandlerID
	CastType  CastType
	CastFor   []ClientID
	Body      []byte
}
```

ResponseMessage is the type when sending from server to client.  

- `SenderID` is the sender client ID of the request message.  
- `HanderID` is the executed handler ID.   
- `Body` is the response binary that was formatted as you like.

`CastType` is message broadcasting type:  

- All - The message will be sent for all connected clients.
- Self - The message will be sent for only sender client.
- Exclusive - The message will be sent for only clients that specified `CastFor` slice.

You don't need to set the `SenderID` and `HandlerID`, these are granted by the system automatically.  

## EventMessage

```go
type EventMessage struct {
	EventType EventType
	ClientID  ClientID
}
```

EventMessage is the type when sending from server to client.  
There are following types of events:

- OnAddClient - emit when new client adding to the server
- OnDelClient - emit when client deleting from server

# Handlers 

go-message-broadcaster can define several handlers for each message and event as following types:

```go
type MessageHandlerID int
type MessageHandler func(msg *RequestMessage, ctx context.Context) (*ResponseMessage, error)
type MessageHandlers map[MessageHandlerID]MessageHandler

type EventType int
type EventHandler func(msg *EventMessage, ctx context.Context) (*ResponseMessage, error)
type EventHandlers map[EventType]EventHandler

```

These handlers are set when the server initialized.  
These message handlers are called when the server received the request from clients.  
These event handlers are called when the client adding or deleting.  

# Contribution

1. Fork ([https://github.com/dev-cloverlab/carpenter/cmd/carpenter/fork](https://github.com/dev-cloverlab/carpenter/cmd/carpenter/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

# Author

[@hatajoe](https://twitter.com/hatajoe)

# Licence

MIT
