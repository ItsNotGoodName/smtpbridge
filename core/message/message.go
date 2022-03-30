package message

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

var (
	ErrNotFound = fmt.Errorf("message not found")
)

type (
	Message struct {
		ID        int64               // ID of the message, should only increment.
		From      string              // From is the email address of the sender.
		To        map[string]struct{} // To is the email addresses of the recipients.
		Subject   string              // Subject of the message.
		Text      string              // Text is the message body.
		CreatedAt time.Time           // Time message was received.
		Processed bool                // Processed is true if the message has been processed.
	}

	Param struct {
		From    string              // From is the email address of the sender.
		To      map[string]struct{} // To is the email addresses of the recipients.
		Subject string              // Subject of the message.
		Text    string              // Text is the message text.
	}

	ListParam struct {
		Cursor   paginate.Cursor
		Messages []Message
	}

	Service interface {
		// Create create and saves a new message.
		Create(ctx context.Context, param *Param) (*Message, error)
		// Processed should called when message has been processed.
		Processed(ctx context.Context, msg *Message) error
	}

	Repository interface {
		// Create a new message.
		Create(ctx context.Context, msg *Message) error
		// Count number of messages.
		Count(ctx context.Context) (int, error)
		// Get message by ID.
		Get(ctx context.Context, id int64) (*Message, error)
		// List messages.
		List(ctx context.Context, param *ListParam) error
		// Update message.
		Update(ctx context.Context, msg *Message, updateFN func(msg *Message) (*Message, error)) error
		// Delete message.
		Delete(ctx context.Context, msg *Message) error
	}
)

func New(param *Param) *Message {
	return &Message{
		CreatedAt: time.Now(),
		Subject:   param.Subject,
		From:      param.From,
		To:        param.To,
		Text:      param.Text,
		Processed: false,
	}
}

func (m *Message) SetProcessed() {
	m.Processed = true
}
