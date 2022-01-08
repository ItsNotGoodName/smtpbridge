package app

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Message struct {
	UUID        string       `json:"uuid"`
	Status      int          `json:"status"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	CreatedAt   string       `json:"created_at"`
	Attachments []Attachment `json:"attachments"`
}

func NewMessage(msg *core.Message, atts []core.Attachment) Message {
	var attachments []Attachment
	for _, attachment := range atts {
		attachments = append(attachments, Attachment{
			UUID: attachment.UUID,
			Name: attachment.Name,
			File: attachment.File(),
		})
	}

	var to []string
	for toAddr := range msg.To {
		to = append(to, toAddr)
	}

	subject := msg.Subject
	if msg.Subject == "" {
		subject = "No Subject"
	}

	return Message{
		UUID:        msg.UUID,
		CreatedAt:   msg.CreatedAt.Format(time.RFC822),
		From:        msg.From,
		To:          to,
		Status:      int(msg.Status),
		Subject:     subject,
		Text:        msg.Text,
		Attachments: attachments,
	}
}
