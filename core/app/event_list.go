package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/cursor"
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
)

func (a *App) eventList(ctx context.Context, req *dto.EventListRequest, exec func(context.Context, *event.ListParam) error) (*dto.EventListResponse, error) {
	listParam := event.ListParam{
		Cursor: cursor.New(req.Ascending, req.Limit, req.Cursor),
	}
	err := exec(ctx, &listParam)
	if err != nil {
		return nil, err
	}

	res := dto.EventListResponse{
		Events:     make([]dto.Event, 0, len(listParam.Events)),
		NextCursor: listParam.Cursor.NextCursor,
		HasMore:    listParam.Cursor.HasMore,
	}
	for _, ev := range listParam.Events {
		res.Events = append(res.Events, *newEvent(&ev))
	}

	return &res, nil

}

func (a *App) MessageEventList(ctx context.Context, req *dto.EventListRequest) (*dto.EventListResponse, error) {
	return a.eventList(ctx, req, func(ctx context.Context, param *event.ListParam) error {
		return a.eventRepository.ListByEntityID(ctx, param, entity.Message, req.EntityID)
	})
}

func (a *App) EventList(ctx context.Context, req *dto.EventListRequest) (*dto.EventListResponse, error) {
	return a.eventList(ctx, req, func(ctx context.Context, param *event.ListParam) error {
		return a.eventRepository.List(ctx, param)
	})
}
