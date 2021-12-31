package app

import (
	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageListRequest struct {
	Page           int
	AttachmentPath string
}

func (a *App) MessageList(req *MessageListRequest) ([]dto.Message, error) {
	if req.Page < 0 {
		req.Page = 0
	}

	msgs, err := a.messageSVC.List(10, req.Page*10)
	if err != nil {
		return nil, err
	}

	var result []dto.Message
	for _, msg := range msgs {
		result = append(result, *dto.NewMessage(&msg, req.AttachmentPath))
	}

	return result, nil
}
