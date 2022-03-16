package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

func newMessage(msg *message.Message, atts []attachment.Attachment) dto.Message {
	var attachments []dto.Attachment = make([]dto.Attachment, 0, len(atts))
	for _, att := range atts {
		attachments = append(attachments, newAttachment(&att))
	}

	var to []string
	for toAddr := range msg.To {
		to = append(to, toAddr)
	}

	return dto.Message{
		ID:          msg.ID,
		CreatedAt:   msg.CreatedAt,
		From:        msg.From,
		To:          to,
		Subject:     msg.Subject,
		Text:        msg.Text,
		Attachments: attachments,
	}
}
