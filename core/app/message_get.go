package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func (a *App) MessageGet(ctx context.Context, req *dto.MessageGetRequest) (*dto.Message, error) {
	msg, err := a.messageRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	atts, err := a.attachmentRepository.ListByMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	dtomsg := newMessage(msg, atts)
	return &dtomsg, nil
}
