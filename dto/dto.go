package dto

import (
	"path"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Message struct {
	UUID        string       `json:"uuid"`
	Status      string       `json:"status"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
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
		From:        msg.From,
		To:          to,
		Status:      msg.Status.String(),
		Subject:     msg.Subject,
		Text:        msg.Text,
		Attachments: attachments,
	}
}
