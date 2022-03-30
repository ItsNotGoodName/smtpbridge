package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

func (a *App) MessageList(ctx context.Context, req *dto.MessageListRequest) (*dto.MessageListResponse, error) {
	listParam := message.ListParam{
		Cursor: paginate.NewCursor(req.Ascending, req.Limit, req.Cursor),
	}
	err := a.messageRepository.List(ctx, &listParam)
	if err != nil {
		return nil, err
	}

	res := dto.MessageListResponse{
		Messages:   []dto.Message{},
		NextCursor: listParam.Cursor.NextCursor,
		HasMore:    listParam.Cursor.HasMore,
	}
	for _, msg := range listParam.Messages {
		atts, err := a.attachmentRepository.ListByMessage(ctx, &msg)
		if err != nil {
			return nil, err
		}

		res.Messages = append(res.Messages, newMessage(&msg, atts))
	}

	return &res, nil
}
