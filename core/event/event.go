package event

import "sync"

type (
	Event struct {
		Topic Topic
		Data  interface{}
	}

	Topic string
)

const (
	TopicEnvelopeCreated Topic = "envelope.created" // *envelope.Envelope
)

type Pub struct {
	subsMu sync.Mutex
	subs   map[Topic][]chan<- Event
}

func NewPub() *Pub {
	return &Pub{
		subs: make(map[Topic][]chan<- Event),
	}
}

func (ps *Pub) Publish(ev Event) {
	ps.subsMu.Lock()
	chs, ok := ps.subs[ev.Topic]
	if !ok {
		ps.subsMu.Unlock()
		return
	}

	for _, ch := range chs {
		ch <- ev
	}
	ps.subsMu.Unlock()
}

func (ps *Pub) Subscribe(topic Topic, ch chan<- Event) {
	ps.subsMu.Lock()
	if chs, ok := ps.subs[topic]; ok {
		ps.subs[topic] = append(chs, ch)
	} else {
		ps.subs[topic] = []chan<- Event{ch}
	}
	ps.subsMu.Unlock()
}
