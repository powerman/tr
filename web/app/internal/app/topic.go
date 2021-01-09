//go:generate gobin -m -run github.com/cheekybits/genny -in=$GOFILE -out=gen.$GOFILE gen "Initialized=RecordAdded"
//go:generate sed -i -e "\\,^//go:generate,d" gen.$GOFILE

package app

import (
	"sync"

	"github.com/vugu/vugu"
)

type (
	// TopicInitialized manages subscription on EventInitialized.
	// While filtering can be implemented inside OnInitialized method,
	// using filters may increase performance because filters will be
	// called without locked EventEnv. Subscriber will get an event if
	// there are no filters provided or ANY of given filters match.
	TopicInitialized interface {
		SubscribeInitialized(c OnInitialized, ee EventEnver, filters ...filterInitialized)
		UnsubscribeInitialized(c OnInitialized)
	}
	// OnInitialized must be implemented by EventInitialized subscribers.
	// It'll be called with locked EventEnv.
	OnInitialized interface {
		OnInitialized(EventInitialized) bool // Return true to request re-render.
	}
	// FilterInitialized is not exported to make it easier to ensure
	// all filters are safe to call without locked EventEnv by
	// implementing all filters in app package.
	// TODO Add custom linter to ensure filters won't get ref in args?
	filterInitialized func(EventInitialized) bool
	// topicInitialized provides a way to broadcast EventInitialized to any
	// amount of subscribers.
	//
	// It subscribe/unsubscribe/emit methods are guaranteed to not
	// block and thus safe to call everywhere.
	//
	// Ref to zero value is ready to use topic.
	topicInitialized struct {
		mu          sync.Mutex
		queue       []EventInitialized
		subs        map[OnInitialized]*subscriptionInitialized
		broadcastMu sync.Mutex
	}
	subscriptionInitialized struct {
		c       OnInitialized
		ee      vugu.EventEnv
		filters []filterInitialized
	}
)

// SubscribeInitialized adds or replaces c subscription to the topic using
// given filters. Provided ee is locked before calling c.OnInitialized and
// unlocked (with render, if it'll return true) when it returns.
func (t *topicInitialized) SubscribeInitialized(c OnInitialized, ee EventEnver, filters ...filterInitialized) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.subs == nil {
		t.subs = make(map[OnInitialized]*subscriptionInitialized)
	}
	t.subs[c] = &subscriptionInitialized{
		c:       c,
		ee:      ee.EventEnv(),
		filters: filters,
	}
}

// UnsubscribeInitialized unsubscribes c from the topic.
func (t *topicInitialized) UnsubscribeInitialized(c OnInitialized) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.subs, c)
}

// EmitInitialized sends ev to all subscribers whose filters accepts it.
// It isn't exported because only app should be able to emit such events.
func (t *topicInitialized) emitInitialized(ev EventInitialized) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.queue = append(t.queue, ev)
	go t.broadcastOne()
}

func (t *topicInitialized) broadcastOne() {
	t.broadcastMu.Lock()
	defer t.broadcastMu.Unlock()

	t.mu.Lock()
	var ev EventInitialized
	ev, t.queue = t.queue[0], t.queue[1:]
	subs := make([]*subscriptionInitialized, 0, len(t.subs))
	for c := range t.subs {
		subs = append(subs, t.subs[c])
	}
	t.mu.Unlock()

	for _, s := range subs {
		pass := len(s.filters) == 0
		for i := 0; !pass && i < len(s.filters); i++ {
			pass = s.filters[i](ev)
		}
		if pass {
			s.ee.Lock()
			if s.c.OnInitialized(ev) {
				s.ee.UnlockRender()
			} else {
				s.ee.UnlockOnly()
			}
		}
	}
}

// FilterInitializedAnd returns true if all filters returns true or no filters
// given.
func FilterInitializedAnd(filters ...filterInitialized) filterInitialized {
	return func(ev EventInitialized) bool {
		for _, filter := range filters {
			if !filter(ev) {
				return false
			}
		}
		return true
	}
}
