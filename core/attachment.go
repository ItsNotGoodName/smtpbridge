package core

import (
	"fmt"
	"io/fs"
)

const (
	TypePNG  AttachmentType = "png"
	TypeJPEG AttachmentType = "jpeg"
)

var ErrAttachmentInvalid = fmt.Errorf("invalid attachment")

type (
	Attachment struct {
		UUID        string
		Name        string
		Type        AttachmentType
		MessageUUID string
		Data        []byte
	}

	EndpointAttachment struct {
		Name string
		Type AttachmentType
		Data []byte
	}

	AttachmentRepositoryPort interface {
		// Create saves a new attachment.
		Create(att *Attachment) error
		// Get returns an attachment by it's UUID without data.
		Get(uuid string) (*Attachment, error)
		// GetData returns the data for an attachment.
		GetData(att *Attachment) ([]byte, error)
		// GetFS returns the attachment file system.
		GetFS() fs.FS
		// ListByMessage returns a list of attachments for a message without data.
		ListByMessage(msg *Message) ([]Attachment, error)
		// DeleteData deletes the data for an attachment.
		DeleteData(att *Attachment) error
	}

	AttachmentType string
)

func (a *Attachment) File() string {
	return fmt.Sprintf("%s.%s", a.UUID, a.Type)
}

func NewEndpointAttachments(atts []Attachment) []EndpointAttachment {
	eats := make([]EndpointAttachment, len(atts))
	for i, a := range atts {
		eats[i] = EndpointAttachment{
			Name: a.Name,
			Type: a.Type,
			Data: a.Data,
		}
	}

	return eats
}
