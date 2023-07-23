package core

import (
	"sync"
)

type EventEnvelopeCreated struct {
	ID int64
}

type EventEnvelopeDeleted struct {
	IDS []int64
}

type EventGardenStart struct {
	Response chan<- bool
}

type Bus struct {
	Mutex           sync.Mutex
	EnvelopeCreated []func(cc *Context, evt EventEnvelopeCreated)
	EnvelopeDeleted []func(cc *Context, evt EventEnvelopeDeleted)
	GardenStart     func(cc *Context, evt EventGardenStart)
}

func NewBus() *Bus {
	return &Bus{}
}
