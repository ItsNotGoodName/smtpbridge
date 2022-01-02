package app

import (
	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type MessageCreateRequest struct {
	Subject     string
	From        string
	To          map[string]bool
	Text        string
	attachments []attachmentCreateRequest
}

type attachmentCreateRequest struct {
	name string
	data []byte
}

func (c *MessageCreateRequest) AddAttachment(name string, data []byte) {
	c.attachments = append(c.attachments, attachmentCreateRequest{name, data})
}

func (a *App) MessageCreate(req *MessageCreateRequest) (*domain.Message, error) {
	msg, err := a.messageSVC.Create(req.Subject, req.From, req.To, req.Text)
	if err != nil {
		return nil, err
	}

	for _, attachment := range req.attachments {
		_, err = a.messageSVC.CreateAttachment(msg, attachment.name, attachment.data)
		if err != nil {
			a.messageREPO.Delete(msg)
			return nil, err
		}
	}

	return msg, nil
}
