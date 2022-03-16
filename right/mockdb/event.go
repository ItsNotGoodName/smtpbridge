package mockdb

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type Event struct{}

func NewEvent() *Event {
	return &Event{}
}

func (Event) Create(ev *event.Event) error {
	return nil
}

func (Event) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (Event) Get(ctx context.Context, id int64) (*event.Event, error) {
	return nil, message.ErrNotFound
}

func (Event) List(ctx context.Context, param *event.ListParam) error {
	return nil
}

func (Event) ListByEntity(ctx context.Context, param *event.ListParam, entity entity.Entity) error {
	return nil
}

func (Event) ListByEntityID(ctx context.Context, param *event.ListParam, entity entity.Entity, entityID int64) error {
	return nil
}
