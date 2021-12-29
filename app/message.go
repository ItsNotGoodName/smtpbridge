package app

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	CreatedAt   time.Time       `json:"created_at"`      // Time message was received.
	UUID        string          `json:"uuid" storm:"id"` // UUID of the message.
	Subject     string          `json:"subject"`         // Subject of the message.
	From        string          `json:"from"`            // From is the email address of the sender.
	To          map[string]bool `json:"to"`              // To is the email addresses of the recipients.
	Text        string          `json:"text"`            // Text is the message body.
	Attachments []Attachment    `json:"-"`               // Attachment is the attachments of the message.
	Status      Status          `json:"status"`          // Status is the status of the message.
}

type Status uint8

const (
	StatusCreated Status = iota
	StatusSent
	StatusFailed
)

func NewMessage(subject, from string, to map[string]bool, text string) *Message {
	return &Message{
		CreatedAt: time.Now(),
		UUID:      uuid.New().String(),
		Subject:   subject,
		From:      from,
		To:        to,
		Text:      text,
		Status:    StatusCreated,
	}
}

type EndpointMessage struct {
	Text        string               // Text is the message body.
	Attachments []EndpointAttachment // Attachments is a list of attachments.
}

func (em *EndpointMessage) IsEmpty() bool {
	return em.Text == "" && len(em.Attachments) == 0
}
