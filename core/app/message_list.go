package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

func (a *App) MessageList(ctx context.Context, req *dto.MessageListRequest) (*dto.MessageListResponse, error) {
	listParam := message.ListParam{
		Cursor: paginate.NewCursor(req.Cursor, req.Limit, req.Ascending),
	}
	err := a.messageRepository.List(ctx, &listParam)
	if err != nil {
		return nil, err
	}

	res := dto.MessageListResponse{
		Messages:   []dto.Message{},
		HasBack:    listParam.Cursor.HasBack(),
		BackCursor: listParam.Cursor.BackCursor,
		NextCursor: listParam.Cursor.NextCursor,
		HasNext:    listParam.Cursor.HasNext(),
	}
	for _, msg := range listParam.Messages {
		atts, err := a.attachmentRepository.ListByMessage(ctx, &msg)
		if err != nil {
			return nil, err
		}

		res.Messages = append(res.Messages, a.newMessage(&msg, atts))
	}

	return &res, nil
}
