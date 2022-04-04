package endpoint

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

var (
	ErrNotFound      = fmt.Errorf("endpoint not found")
	ErrNameConflict  = fmt.Errorf("endpoint name conflict")
	ErrInvalidConfig = fmt.Errorf("endpoint config invalid")
	ErrInvalidType   = fmt.Errorf("endpoint type invalid")
	ErrSendFailed    = fmt.Errorf("endpoint send failed")
)

func NewFacade(name string, endpoint Endpoint) *Facade {
	return &Facade{
		Name:     name,
		Endpoint: endpoint,
	}
}

func NewEnvelope(msg *Message, atts []Attachment) Envelope {
	return Envelope{
		Message:     msg,
		Attachments: atts,
	}
}

func NewMessage(msg *message.Message) *Message {
	return &Message{
		ID:    msg.ID,
		Title: msg.Subject,
		Body:  msg.Text,
	}
}

func NewAttachments(atts []attachment.Attachment) ([]Attachment, error) {
	var a []Attachment
	for _, att := range atts {
		data, err := att.GetData()
		if err != nil {
			return nil, err
		}
		a = append(a, Attachment{
			Name: att.Name,
			Data: data,
		})
	}
	return a, nil
}

// Text returns the text of the message.
func (m *Message) Text() string {
	return m.Title + "\n" + m.Body
}

// Send sends a message to the endpoint.
func (f *Facade) Send(ctx context.Context, msg Envelope) error {
	if err := f.Endpoint.Send(ctx, msg); err != nil {
		return fmt.Errorf("%v: %w", f.Name, err)
	}
	return nil
}
