package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func (a *App) MessageDelete(ctx context.Context, req *dto.MessageDeleteRequest) error {
	msg, err := a.messageRepository.Get(ctx, req.ID)
	if err != nil {
		return err
	}

	return a.messageRepository.Delete(ctx, msg)
}
