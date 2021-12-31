package dto

import (
	"path"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Message struct {
	UUID        string        `json:"uuid"`
	Status      domain.Status `json:"status"`
	From        string        `json:"from"`
	To          []string      `json:"to"`
	Subject     string        `json:"subject"`
	Text        string        `json:"text"`
	CreatedAt   string        `json:"created_at"`
	Attachments []Attachment  `json:"attachments"`
}

type Attachment struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewMessage(msg *domain.Message, attachmentPath string) *Message {
	var attachments []Attachment
	for _, attachment := range msg.Attachments {
		attachments = append(attachments, Attachment{
			UUID: attachment.UUID,
			Name: attachment.Name,
			Path: path.Join(attachmentPath, attachment.File()),
		})
	}

	var to []string
	for toAddr := range msg.To {
		to = append(to, toAddr)
	}

	return &Message{
		UUID:        msg.UUID,
		CreatedAt:   msg.CreatedAt.Format(time.RFC822),
		From:        msg.From,
		To:          to,
		Status:      msg.Status,
		Subject:     msg.Subject,
		Text:        msg.Text,
		Attachments: attachments,
	}
}
