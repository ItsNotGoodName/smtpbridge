package app

import (
	"math"
)

type MessageListRequest struct {
	Page int
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

	count, err := a.messageREPO.Count()
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

	var result []Message
	for _, msg := range msgs {
		result = append(result, NewMessage(&msg))
	}

	return &MessageListResponse{Messages: result, PageMin: pageMin, PageMax: pageMax, Page: req.Page}, nil
}
