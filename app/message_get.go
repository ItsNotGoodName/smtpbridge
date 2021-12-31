package app

import (
	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageGetRequest struct {
	UUID string
}

func (a *App) MessageGet(req *MessageGetRequest) (*dto.Message, error) {
	msg, err := a.dao.Message.GetMessage(req.UUID)
	if err != nil {
		return nil, err
	}

	return dto.NewMessage(msg), nil
}
