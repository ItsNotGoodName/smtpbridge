package event

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type MessageService struct {
	eventService   Service
	messageService message.Service
}

func NewMessageService(eventService Service, messageService message.Service) *MessageService {
	return &MessageService{
		eventService:   eventService,
		messageService: messageService,
	}
}

func (ms *MessageService) Create(ctx context.Context, param *message.Param) (*message.Message, error) {
	msg, err := ms.messageService.Create(ctx, param)
	if err == nil {
		err := ms.eventService.Create(New(MessageCreated).WithEntity(entity.Message, msg.ID).Done())
		if err != nil {
			return nil, err
		}
	}

	return msg, err
}

func (ms *MessageService) Processed(ctx context.Context, msg *message.Message) error {
	err := ms.messageService.Processed(ctx, msg)
	if err == nil {
		err := ms.eventService.Create(New(MessageProcessed).WithEntity(entity.Message, msg.ID).Done())
		if err != nil {
			return err
		}
	}

	return err
}
