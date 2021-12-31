package app

import (
	"math"

	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageListRequest struct {
	Page           int
	AttachmentPath string
}

type MessageListResponse struct {
	Messages []dto.Message
	Page     int
	PageMax  int
}

func (a *App) MessageList(req *MessageListRequest) (*MessageListResponse, error) {
	limit := 10

	count, err := a.dao.Message.CountMessages()
	if err != nil {
		return nil, err
	}

	pageMax := int(math.Ceil(float64(count) / float64(limit)))

	if req.Page < 1 || req.Page > pageMax {
		req.Page = 1
	}

	msgs, err := a.messageSVC.List(limit, (req.Page-1)*10)
	if err != nil {
		return nil, err
	}

	var result []dto.Message
	for _, msg := range msgs {
		result = append(result, *dto.NewMessage(&msg, req.AttachmentPath))
	}

	return &MessageListResponse{Messages: result, PageMax: pageMax, Page: req.Page}, nil
}
