package event

import (
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/entity"
)

var (
	ErrNotFound = fmt.Errorf("event not found")
)

const (
	MessageCreated   = Name("message.created")
	MessageProcessed = Name("message.processed")
	EndpointError    = Name("endpoint.error")
	EndpointSuccess  = Name("endpoint.success")
)

func New(name Name) Creator {
	return Creator{&Event{
		Name:      name,
		Entity:    entity.None,
		CreatedAt: time.Now(),
	}}
}

func (c Creator) Done() *Event {
	return c.e
}

func (c Creator) WithName(name Name) Creator {
	c.e.Name = name
	return c
}

func (c Creator) WithDescription(description string) Creator {
	c.e.Description = description
	return c
}

func (c Creator) WithEntity(e entity.Entity, eid int64) Creator {
	c.e.Entity = e
	c.e.EntityID = eid
	return c
}
