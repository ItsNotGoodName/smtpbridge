package domain

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMessageNotFound      = fmt.Errorf("message not found")
	ErrMessageAlreadyExists = fmt.Errorf("message already exists")
)

const (
	StatusCreated Status = iota
	StatusSent
	StatusFailed
	StatusSkipped
)

type (
	Message struct {
		UUID        string          `json:"uuid" storm:"id"` // UUID of the message.
		From        string          `json:"from"`            // From is the email address of the sender.
		To          map[string]bool `json:"to"`              // To is the email addresses of the recipients.
		Subject     string          `json:"subject"`         // Subject of the message.
		Text        string          `json:"text"`            // Text is the message body.
		Attachments []Attachment    `json:"-"`               // Attachment is the attachments of the message.
		Status      Status          `json:"status"`          // Status is the status of the message.
		CreatedAt   time.Time       `json:"created_at"`      // Time message was received.
	}

	EndpointMessage struct {
		Text        string               // Text is the message body.
		Attachments []EndpointAttachment // Attachments is a list of attachments.
	}

	MessageServicePort interface {
		// Create creates a new message and saves it.
		Create(subject, from string, to map[string]bool, text string) (*Message, error)
		// CreateAttachment adds an attachment to a message.
		CreateAttachment(msg *Message, name string, data []byte) (*Attachment, error)
		// Get a message with attachments.
		Get(uuid string) (*Message, error)
		// List messages with attachments.
		List(limit, offset int) ([]Message, error)
		// UpdateStatus updates the status of a message.
		UpdateStatus(msg *Message, status Status) error
	}

	// MessageRepositoryPort handles storing messages.
	MessageRepositoryPort interface {
		// CreateMessage saves a new message.
		Create(msg *Message) error
		// CountMessages returns the number of messages.
		Count() (int, error)
		// GetMessage returns a message by it's UUID.
		Get(uuid string) (*Message, error)
		// GetMessages returns a list of messages.
		List(limit, offset int) ([]Message, error)
		// UpdateMessage updates a message.
		Update(msg *Message, updateFN func(msg *Message) (*Message, error)) error
		// DeleteMessage deletes a message.
		Delete(msg *Message) error
	}

	Status uint8
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

func (m *Message) NewAttachment(name string, data []byte) (*Attachment, error) {
	var t AttachmentType
	contentType := http.DetectContentType(data)
	if contentType == "image/png" {
		t = TypePNG
	} else if contentType == "image/jpeg" {
		t = TypeJPEG
	} else {
		return nil, fmt.Errorf("%s: %v", contentType, ErrAttachmentInvalid)
	}

	att := Attachment{
		UUID:        uuid.New().String(),
		Name:        name,
		Type:        t,
		MessageUUID: m.UUID,
		Data:        data,
	}

	m.Attachments = append(m.Attachments, att)

	return &att, nil
}

func (em *EndpointMessage) IsEmpty() bool {
	return em.Text == "" && len(em.Attachments) == 0
}
