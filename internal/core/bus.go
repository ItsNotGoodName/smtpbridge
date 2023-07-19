package core

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type EventEnvelopeCreated struct {
	ID int64
}

type EventEnvelopeDeleted struct {
	ID int64
}

type EventStorageRead struct {
	Storage models.Storage
}

type Bus struct {
	Mutex           sync.Mutex
	EnvelopeCreated []func(cc *Context, evt EventEnvelopeCreated)
	EnvelopeDeleted []func(cc *Context, evt EventEnvelopeDeleted)
	StorageRead     []func(cc *Context, evt EventStorageRead)
}

func NewBus() *Bus {
	return &Bus{}
}
