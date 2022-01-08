package core

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
	StatusAll Status = iota
	StatusCreated
	StatusSent
	StatusFailed
	StatusSkipped
)

type (
	MessageParam struct {
		Limit   int
		Offset  int
		Status  Status
		Reverse bool
	}

	Message struct {
		UUID      string              // UUID of the message.
		From      string              // From is the email address of the sender.
		To        map[string]struct{} // To is the email addresses of the recipients.
		Subject   string              // Subject of the message.
		Text      string              // Text is the message body.
		Status    Status              // Status is the status of the message.
		CreatedAt time.Time           // Time message was received.
	}

	EndpointMessage struct {
		Text        string               // Text is the message body.
		Attachments []EndpointAttachment // Attachments is a list of attachments.
	}

	MessageServicePort interface {
		// Create a new message and saves it.
		Create(subject, from string, to map[string]struct{}, text string) (*Message, error)
		// CreateAttachment adds an attachment to a message and saves it.
		CreateAttachment(msg *Message, name string, data []byte) (*Attachment, error)
		// UpdateStatus updates the status of a message.
		UpdateStatus(msg *Message, status Status) error
	}

	MessageRepositoryPort interface {
		// Create saves a message.
		Create(msg *Message) error
		// Count returns the number of messages.
		Count(search *MessageParam) (int, error)
		// Get returns a message by it's UUID.
		Get(uuid string) (*Message, error)
		// List messages.
		List(search *MessageParam) ([]Message, error)
		// Update a message.
		Update(msg *Message, updateFN func(msg *Message) (*Message, error)) error
		// Delete a message.
		Delete(msg *Message) error
	}

	Status uint8
)

func NewMessage(subject, from string, to map[string]struct{}, text string) *Message {
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

// AttachmentDataValid returns the type of the attachment data.
func AttachmentDataValid(data []byte) (AttachmentType, error) {
	if len(data) == 0 {
		return "", ErrAttachmentNoData
	}

	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/png":
		return TypePNG, nil
	case "image/jpeg":
		return TypeJPEG, nil
	default:
		return "", fmt.Errorf("%s: %v", contentType, ErrAttachmentInvalid)
	}
}

func (m *Message) NewAttachment(name string, data []byte) (*Attachment, error) {
	attType, err := AttachmentDataValid(data)
	if err != nil {
		return nil, err
	}

	return &Attachment{
		UUID:        uuid.New().String(),
		Name:        name,
		Type:        attType,
		MessageUUID: m.UUID,
		data:        data,
	}, nil
}

func (em *EndpointMessage) IsEmpty() bool {
	return em.Text == "" && len(em.Attachments) == 0
}
