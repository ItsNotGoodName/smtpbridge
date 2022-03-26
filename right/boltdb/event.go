package boltdb

import (
	"context"
	"math"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type eventModel struct {
	ID          int64     `json:"id" storm:"id,increment"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Entity      int       `json:"entity" storm:"index"`
	EntityID    int64     `json:"entity_id" storm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

func convertEventD(evD *event.Event) *eventModel {
	return &eventModel{
		ID:          evD.ID,
		Name:        string(evD.Name),
		Description: evD.Description,
		Entity:      int(evD.Entity),
		EntityID:    evD.EntityID,
		CreatedAt:   evD.CreatedAt,
	}
}

func convertEventM(evM *eventModel) *event.Event {
	return &event.Event{
		ID:          evM.ID,
		Name:        event.Name(evM.Name),
		Description: evM.Description,
		Entity:      entity.Entity(evM.Entity),
		EntityID:    evM.EntityID,
		CreatedAt:   evM.CreatedAt,
	}
}

type Event struct {
	db *storm.DB
}

func NewEvent(database *Database) *Event {
	return &Event{
		db: database.db,
	}
}

func (e *Event) Create(ev *event.Event) error {
	evD := convertEventD(ev)
	if err := e.db.Save(evD); err != nil {
		return err
	}

	ev.ID = evD.ID
	return nil
}

func (e *Event) Count(ctx context.Context) (int, error) {
	count, err := e.db.Select().Count(&eventModel{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, err
}

func (e *Event) Get(ctx context.Context, id int64) (*event.Event, error) {
	var evM *eventModel
	err := e.db.One("ID", id, evM)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, event.ErrNotFound
		}
		return nil, err
	}

	return convertEventM(evM), nil
}

func (e *Event) list(ctx context.Context, param *event.ListParam, filters ...q.Matcher) error {
	var query storm.Query
	cursorOffset := 0
	if param.Cursor.Ascending {
		filters = append(filters, q.Gt("ID", param.Cursor.Cursor))
		query = e.db.Select(filters...).OrderBy("ID")
	} else {
		if param.Cursor.Cursor == 0 {
			param.Cursor.Cursor = math.MaxInt64
		}
		filters = append(filters, q.Lt("ID", param.Cursor.Cursor))
		query = e.db.Select(filters...).OrderBy("ID").Reverse()
		cursorOffset = 1
	}

	query.Limit(param.Cursor.Limit + 1)

	var eventsM []eventModel
	if err := query.Find(&eventsM); err != nil && err != storm.ErrNotFound {
		return err
	}

	if len(eventsM) == param.Cursor.Limit+1 {
		param.Cursor.SetHasMore(true)
		param.Cursor.SetNextCursor(eventsM[param.Cursor.Limit-cursorOffset].ID)
		eventsM = eventsM[:param.Cursor.Limit]
	}

	var events []event.Event
	for _, evM := range eventsM {
		events = append(events, *convertEventM(&evM))
	}
	param.Events = events

	return nil
}

func (e *Event) List(ctx context.Context, param *event.ListParam) error {
	return e.list(ctx, param)
}

func (e *Event) ListByEntity(ctx context.Context, param *event.ListParam, entity entity.Entity) error {
	return e.list(ctx, param, q.Eq("Entity", int(entity)))
}

func (e *Event) ListByEntityID(ctx context.Context, param *event.ListParam, entity entity.Entity, entityID int64) error {
	return e.list(ctx, param, q.Eq("Entity", int(entity)), q.Eq("EntityID", entityID))
}
