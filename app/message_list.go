package app

import "github.com/ItsNotGoodName/smtpbridge/domain"

type MessageListRequest struct {
	Limit  int
	Offset int
}

func (a *App) MessageList(req *MessageListRequest) ([]domain.Message, error) {
	return a.messageSVC.List(req.Limit, req.Offset)
}
