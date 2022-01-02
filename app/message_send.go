package app

import "github.com/ItsNotGoodName/smtpbridge/domain"

type MessageSendRequest struct {
	Message *domain.Message
}

func (a *App) MessageSend(req *MessageSendRequest) error {
	err := a.endpointSVC.SendBridges(req.Message, a.bridgeSVC.ListByMessage(req.Message))
	if err != nil {
		a.messageSVC.UpdateStatus(req.Message, domain.StatusFailed)
		return err
	}

	return a.messageSVC.UpdateStatus(req.Message, domain.StatusSent)
}
