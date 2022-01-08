package app

import (
	"math"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type MessageListRequest struct {
	Page   int
	Status int
}

type MessageListResponse struct {
	Messages []Message
	Page     int
	PageMin  int
	PageMax  int
}

func (a *App) MessageList(req *MessageListRequest) (*MessageListResponse, error) {
	limit := 10
	pageMin := 1

	searchParam := &core.MessageParam{
		Limit:   limit,
		Offset:  (req.Page - pageMin) * limit,
		Reverse: true,
		Status:  core.Status(req.Status),
	}

	count, err := a.messageREPO.Count(searchParam)
	if err != nil {
		return nil, err
	}

	pageMax := int(math.Ceil(float64(count) / float64(limit)))

	if req.Page < pageMin || req.Page > pageMax {
		req.Page = pageMin
	}

	msgs, err := a.messageREPO.List(searchParam)
	if err != nil {
		return nil, err
	}

	var result []Message
	for _, msg := range msgs {
		atts, err := a.attachmentREPO.ListByMessage(&msg)
		if err != nil {
			return nil, err
		}
		result = append(result, NewMessage(&msg, atts))
	}

	return &MessageListResponse{Messages: result, PageMin: pageMin, PageMax: pageMax, Page: req.Page}, nil
}
