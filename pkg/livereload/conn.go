package livereload

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/powerman/tr/pkg/broadcast"
)

// Conn implements server side of LiveReload connection.
type Conn struct {
	ws           *websocket.Conn
	handshake    chan struct{}
	msgc         chan interface{}
	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// NewConn creates and returns new Conn.
func NewConn(ws *websocket.Conn) *Conn {
	c := &Conn{
		ws:        ws,
		handshake: make(chan struct{}, 1),
		msgc:      make(chan interface{}, 256),
		shutdown:  make(chan struct{}),
	}
	go c.send()
	go c.recv()
	return c
}

// Close is safe to call multiple times.
func (c *Conn) Close() {
	c.shutdownOnce.Do(func() { close(c.shutdown) })
	c.ws.Close()
}

// Send implements broadcast.Subscriber interface.
func (c *Conn) Send() chan<- interface{} { return c.msgc }

// Kick implements broadcast.Subscriber interface.
func (c *Conn) Kick(_ *broadcast.Topic) { c.Close() }

// Wait until connection will be closed - either by Close, or because of
// internal handshake or I/O error.
func (c *Conn) Wait() { <-c.shutdown }

func (c *Conn) recv() {
	defer c.Close()

	for {
		var msg msgClientHello
		if c.ws.ReadJSON(&msg) != nil {
			return
		}

		if msg.Command == cmdHello {
			if !validProto(msg.Protocols) {
				return
			}
			select {
			case c.handshake <- struct{}{}:
			default:
			}
		}
	}
}

func (c *Conn) send() {
	defer c.Close()

	var msgc <-chan interface{}
	for {
		select {
		case <-c.shutdown:
			return
		case <-c.handshake:
			msgc = c.msgc
		case msg := <-msgc:
			if c.ws.WriteJSON(msg) != nil {
				return
			}
		}
	}
}
