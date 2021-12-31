package app

import (
	"math"

	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageListRequest struct {
	Page int
}

type MessageListResponse struct {
	Messages []dto.Message
	Page     int
	PageMin  int
	PageMax  int
}

func (a *App) MessageList(req *MessageListRequest) (*MessageListResponse, error) {
	// TODO: move this logic into message service
	limit := 10
	pageMin := 1

	count, err := a.dao.Message.CountMessages()
	if err != nil {
		return nil, err
	}

	pageMax := int(math.Ceil(float64(count) / float64(limit)))

	if req.Page < pageMin || req.Page > pageMax {
		req.Page = pageMin
	}

	msgs, err := a.messageSVC.List(limit, (req.Page-pageMin)*10)
	if err != nil {
		return nil, err
	}

	var result []dto.Message
	for _, msg := range msgs {
		result = append(result, *dto.NewMessage(&msg))
	}

	return &MessageListResponse{Messages: result, PageMin: pageMin, PageMax: pageMax, Page: req.Page}, nil
}
