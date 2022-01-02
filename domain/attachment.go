package domain

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
		UUID        string         `json:"uuid" storm:"id"`
		Name        string         `json:"name"`
		Type        AttachmentType `json:"type"`
		MessageUUID string         `json:"message_uuid" storm:"index"`
		Data        []byte         `json:"-"`
	}

	EndpointAttachment struct {
		Name string
		Type AttachmentType
		Data []byte
	}

	AttachmentRepositoryPort interface {
		// CreateAttachment saves a new attachment.
		Create(att *Attachment) error
		// GetAttachment returns an attachment by it's UUID.
		Get(uuid string) (*Attachment, error)
		// GetAttachmentData returns the data for an attachment.
		GetData(att *Attachment) ([]byte, error)
		// GetFS returns the attachment file system
		GetFS() fs.FS
		// GetAttachments returns a list of attachments for a message.
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
