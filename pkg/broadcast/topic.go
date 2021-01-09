// Package broadcast helps sending same message to a group of subscribers.
//
// Everything in this package is safe for concurrent use by multiple goroutines.
package broadcast

import "sync"

// Subscriber provides a channel to receive broadcast messages.
// It must be safe for concurrent use by multiple goroutines.
type Subscriber interface {
	// Send should return (usually buffered) channel which must not be
	// closed while being subscribed (to any Topic).
	// It may return nil but this will lead to Kick.
	// It may be called after Kick() (by another Topic).
	Send() chan<- interface{}
	// Kick will be called on Send channel overflow or if Send
	// returns nil. It means Subscriber was unsubscribed and even in
	// case subscriber will call Subscribe again while handling Kick
	// it'll anyway miss at least one message.
	// It will not be called on manual Unsubscribe.
	// It may be called multiple times (once by each Topic).
	Kick(*Topic)
}

// SubscribeReq is used to keep run() non-blocking, even if sub.Send()
// will block.
type subscribeReq struct {
	sub  Subscriber
	send chan<- interface{}
}

// Topic broadcasts messages to subscribers.
type Topic struct {
	subs         map[Subscriber]chan<- interface{}
	subscribe    chan subscribeReq
	unsubscribe  chan Subscriber
	broadcast    chan interface{}
	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// NewTopic creates and returns new Topic.
func NewTopic() *Topic {
	h := &Topic{
		subs:        make(map[Subscriber]chan<- interface{}),
		subscribe:   make(chan subscribeReq),
		unsubscribe: make(chan Subscriber),
		broadcast:   make(chan interface{}, 64),
		shutdown:    make(chan struct{}),
	}
	go h.run()
	return h
}

// Shutdown stops broadcasting messages on the topic.
// After Shutdown any method (including Shutdown) does nothing.
func (h *Topic) Shutdown() {
	h.shutdownOnce.Do(func() { close(h.shutdown) })
}

// Subscribe sub to the topic.
func (h *Topic) Subscribe(sub Subscriber) {
	req := subscribeReq{
		sub:  sub,
		send: sub.Send(),
	}
	select {
	case <-h.shutdown:
	case h.subscribe <- req:
	}
}

// Unsubscribe sub from the topic.
func (h *Topic) Unsubscribe(sub Subscriber) {
	select {
	case <-h.shutdown:
	case h.unsubscribe <- sub:
	}
}

// Broadcast enqueue msg to be sent to all subscribers.
func (h *Topic) Broadcast(msg interface{}) {
	select {
	case <-h.shutdown:
	case h.broadcast <- msg:
	}
}

func (h *Topic) run() {
	for {
		select {
		case <-h.shutdown:
			return
		case req := <-h.subscribe:
			h.subs[req.sub] = req.send
		case sub := <-h.unsubscribe:
			delete(h.subs, sub)
		case m := <-h.broadcast:
			for sub, send := range h.subs {
				select {
				case send <- m:
				default: // Slow reader (sub.Send overflow) or sub.Send returns nil.
					delete(h.subs, sub)
					go sub.Kick(h)
				}
			}
		}
	}
}
