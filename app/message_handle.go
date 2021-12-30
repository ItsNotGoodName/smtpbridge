package app

import "github.com/ItsNotGoodName/smtpbridge/domain"

type MessageSendRequest struct {
	Message *domain.Message
}

func (a *App) MessageSend(req *MessageSendRequest) error {
	return a.endpointSVC.SendBridges(req.Message, a.bridgeSVC.GetBridges(req.Message))
}
