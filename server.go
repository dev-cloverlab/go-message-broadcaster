package broadcaster

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Server struct {
	addClient    chan *Client
	delClient    chan *Client
	err          chan error
	broadcast    chan *ResponseMessage
	messageQueue chan *RequestMessage
	done         chan struct{}
}

func NewServer(ctx context.Context, messageHandlers MessageHandlers) *Server {
	sv := &Server{
		addClient:    make(chan *Client),
		delClient:    make(chan *Client),
		err:          make(chan error),
		broadcast:    make(chan *ResponseMessage),
		messageQueue: make(chan *RequestMessage),
		done:         make(chan struct{}),
	}

	go func() {
		for msg := range sv.messageQueue {
			if cmd, ok := messageHandlers[msg.HandlerID]; ok {
				if res, err := cmd(msg, ctx); err != nil {
					sv.OnError(err)
				} else {
					res.SenderID = msg.SenderID
					res.HandlerID = msg.HandlerID
					sv.OnBroadCast(res)
				}
			} else {
				sv.OnError(fmt.Errorf("undefined message handler specified `%d'", msg.HandlerID))
			}
		}
	}()

	return sv
}

func (m *Server) NewClient(conn Conn) *Client {
	c := NewClient(conn, m)
	m.OnAddClient(c)
	return c
}

func (m *Server) OnAddClient(c *Client) {
	m.addClient <- c
}

func (m *Server) OnDelClient(c *Client) {
	m.delClient <- c
}

func (m *Server) OnError(err error) {
	m.err <- err
}

func (m *Server) OnBroadCast(msg *ResponseMessage) {
	m.broadcast <- msg
}

func (m *Server) OnEnqueueMessage(msg *RequestMessage) {
	m.messageQueue <- msg
}

func (m *Server) OnDone() {
	m.done <- struct{}{}
}

func (m *Server) Listen() {
	clients := []*Client{}
	for {
		select {

		case c := <-m.addClient:
			clients = append(clients, c)
			emit(clients, NewEventMessage(OnAddClient, c.ID))

		case c := <-m.delClient:
			for i, v := range clients {
				if v.ID != c.ID {
					continue
				}
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
			emit(clients, NewEventMessage(OnDelClient, c.ID))
			c.OnDelete()

		case err := <-m.err:
			log.Println(err)

		case msg := <-m.broadcast:
			broadcast(clients, msg)

		case <-m.done:
			l := len(clients)
			for i := 0; i < l; i++ {
				m.delClient <- clients[i]
			}
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func emit(clients []*Client, ev *EventMessage) {
	for _, c := range clients {
		go c.OnEvent(ev)
	}
}

func broadcast(clients []*Client, msg *ResponseMessage) {
	castFor := []*Client{}
	switch msg.CastType {
	case Self:
		for _, c := range clients {
			if c.ID == msg.SenderID {
				castFor = append(castFor, c)
			}
		}
	case Exclusive:
		for _, c := range clients {
			for _, ex := range msg.CastFor {
				if c.ID == ex {
					castFor = append(castFor, c)
				}
			}
		}
	default:
		castFor = clients
	}

	for _, c := range castFor {
		go c.OnSend(msg)
	}
}
