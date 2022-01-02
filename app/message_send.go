package app

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type messageSendRequest struct {
	Message *domain.Message
}

func (a *App) messageSend(req *messageSendRequest) error {
	status, err := a.endpointSVC.SendBridges(req.Message, a.bridgeSVC.ListByMessage(req.Message))
	err2 := a.messageSVC.UpdateStatus(req.Message, status)
	if err == nil {
		return err2
	}
	if err2 == nil {
		return err
	}
	return fmt.Errorf("%s, %s", err, err2)
}
