package app

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Time        time.Time       `json:"time"`        // Time message was received.
	UUID        string          `json:"uuid"`        // UUID of the message.
	Subject     string          `json:"subject"`     // Subject of the message.
	From        string          `json:"from"`        // From is the email address of the sender.
	To          map[string]bool `json:"to"`          // To is the email addresses of the recipients.
	Text        string          `json:"text"`        // Text is the message body.
	Attachments []*Attachment   `json:"attachments"` // Attachments is a list of attachments.
}

func NewMessage(subject, from string, to map[string]bool, text string) *Message {
	return &Message{
		Time:    time.Now(),
		UUID:    uuid.New().String(),
		Subject: subject,
		From:    from,
		To:      to,
		Text:    text,
	}
}

type EndpointMessage struct {
	Text        string        `json:"text"`        // Text is the message body.
	Attachments []*Attachment `json:"attachments"` // Attachments is a list of attachments.
}
