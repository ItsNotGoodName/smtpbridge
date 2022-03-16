package event

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/entity"
)

type EndpointService struct {
	eventService    Service
	endpointService endpoint.Service
}

func NewEndpointService(eventService Service, endpointSerice endpoint.Service) *EndpointService {
	return &EndpointService{
		eventService:    eventService,
		endpointService: endpointSerice,
	}
}

func (es *EndpointService) Send(ctx context.Context, req []endpoint.SendRequest) <-chan endpoint.SendResponse {
	res := es.endpointService.Send(ctx, req)
	resC := make(chan endpoint.SendResponse, len(req))

	go (func() {
		for s := range res {
			var creator Creator
			if s.Error != nil {
				creator = New(EndpointError).WithDescription(s.Error.Error())
			} else {
				creator = New(EndpointSuccess).WithDescription(fmt.Sprintf("sent to '%s'", s.Facade.Name))
			}
			es.eventService.Create(creator.WithEntity(entity.Message, s.Envelope.Message.ID).Done())
			resC <- s
		}

		close(resC)
	})()

	return resC
}
