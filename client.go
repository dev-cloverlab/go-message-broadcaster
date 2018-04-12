package broadcaster

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"
)

type ClientID string

type Client struct {
	ID      ClientID
	conn    Conn
	sv      *Server
	message chan *ResponseMessage
	event   chan *EventMessage
	delete  chan struct{}
	close   chan struct{}
}

func NewClient(conn Conn, sv *Server) *Client {
	return &Client{
		ID:      ClientID(mustUUID()),
		conn:    conn,
		sv:      sv,
		message: make(chan *ResponseMessage),
		event:   make(chan *EventMessage),
		delete:  make(chan struct{}),
		close:   make(chan struct{}),
	}
}

func (m *Client) OnSend(msg *ResponseMessage) {
	m.message <- msg
}

func (m *Client) OnDelete() {
	if err := m.conn.Close(); err != nil {
		log.Printf("connection close failed: client id %s", m.ID)
	}
	m.delete <- struct{}{}
	close(m.close)
}

func (m *Client) Listen() {
	go m.listenWrite()
	m.listenRead()
	<-m.close
}

func (m *Client) listenWrite() {
	for {
		select {
		case msg := <-m.message:
			err := m.conn.Send(msg)
			if err != nil {
				m.sv.OnError(err)
			}
		case <-m.delete:
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (m *Client) listenRead() {
	for {
		msg, err := m.conn.Receive()
		if err == io.EOF {
			m.sv.OnDelClient(m)
			return
		} else if err != nil {
			m.sv.OnError(err)
		} else {
			msg.SenderID = m.ID
			m.sv.OnEnqueueMessage(msg)
		}
	}
}

func mustUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("err: rand.Read failed for reason %s", err.Error()))
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
