package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
)

func newEvent(e *event.Event) *dto.Event {
	return &dto.Event{
		ID:          e.ID,
		Name:        string(e.Name),
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
	}
}
