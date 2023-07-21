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

type Bus struct {
	Mutex           sync.Mutex
	EnvelopeCreated []func(cc *Context, evt EventEnvelopeCreated)
	EnvelopeDeleted []func(cc *Context, evt EventEnvelopeDeleted)
}

func NewBus() *Bus {
	return &Bus{}
}
