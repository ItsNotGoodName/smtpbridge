package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func (a *App) Info(ctx context.Context) (*dto.InfoResponse, error) {
	messagesCount, err := a.messageRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	attachmentsCount, err := a.attachmentRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	attachmentsSize, err := a.dataRepository.Size(ctx)
	if err != nil {
		return nil, err
	}

	eventsCount, err := a.eventRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.InfoResponse{
		MessagesCount:    messagesCount,
		AttachmentsCount: attachmentsCount,
		AttachmentsSize:  attachmentsSize,
		EventsCount:      eventsCount,
	}, nil
}
