package event

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

type (
	Event struct {
		ID          int64         // ID of event, should only increment.
		Name        Name          // Name of event.
		Description string        // Description of event.
		Entity      entity.Entity // Entity is the type of entity that the event belongs to.
		EntityID    int64         // EntityID is the entity's ID.
		CreatedAt   time.Time     // CreatedAt is the time this event was created.
	}

	Creator struct {
		e *Event
	}

	ListParam struct {
		Page   paginate.Page
		Events []Event
	}

	Service interface {
		// Create a new event.
		Create(event *Event) error
	}

	Repository interface {
		// Create a new event.
		Create(event *Event) error
		// Count events.
		Count(ctx context.Context) (int, error)
		// Get an event by id.
		Get(ctx context.Context, id int64) (*Event, error)
		// List events.
		List(ctx context.Context, param *ListParam) error
		// List events by entity.
		ListByEntity(ctx context.Context, param *ListParam, entity entity.Entity) error
		// List events by entity id.
		ListByEntityID(ctx context.Context, param *ListParam, entity entity.Entity, entityID int64) error
	}

	Name string
)
