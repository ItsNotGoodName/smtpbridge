package app

import (
	"path"

	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageGetRequest struct {
	UUID           string
	AttachmentPath string
}

func (a *App) MessageGet(req *MessageGetRequest) (*dto.Message, error) {
	msg, err := a.dao.Message.GetMessage(req.UUID)
	if err != nil {
		return nil, err
	}

	var attachments []dto.Attachment
	for _, attachment := range msg.Attachments {
		attachments = append(attachments, dto.Attachment{
			UUID: attachment.UUID,
			Name: attachment.Name,
			Path: path.Join(req.AttachmentPath, a.dao.Attachment.GetAttachmentFile(&attachment)),
		})
	}

	var to []string
	for toAddr := range msg.To {
		to = append(to, toAddr)
	}

	return &dto.Message{
		UUID:        msg.UUID,
		From:        msg.From,
		To:          to,
		Status:      msg.Status.String(),
		Subject:     msg.Subject,
		Text:        msg.Text,
		Attachments: attachments,
	}, nil
}
