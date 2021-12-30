package app

import (
	"path"

	"github.com/ItsNotGoodName/smtpbridge/dto"
)

type MessageListRequest struct {
	Page           int
	AttachmentPath string
}

func (a *App) MessageList(req *MessageListRequest) ([]dto.Message, error) {
	if req.Page < 0 {
		req.Page = 0
	}

	msgs, err := a.messageSVC.List(10, req.Page*10)
	if err != nil {
		return nil, err
	}

	var result []dto.Message
	for _, msg := range msgs {
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

		result = append(result, dto.Message{
			UUID:        msg.UUID,
			From:        msg.From,
			To:          to,
			Status:      msg.Status.String(),
			Subject:     msg.Subject,
			Text:        msg.Text,
			Attachments: attachments,
		})
	}

	return result, nil
}
