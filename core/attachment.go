package core

import (
	"fmt"
	"io/fs"
)

const (
	TypePNG  AttachmentType = "png"
	TypeJPEG AttachmentType = "jpeg"
)

var (
	ErrAttachmentInvalid  = fmt.Errorf("invalid attachment")
	ErrAttachmentNotFound = fmt.Errorf("attachment not found")
)

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
		// Count returns the number of attachments.
		Count() (int, error)
		// CountByMessage returns the number of attachments for a message.
		CountByMessage(msg *Message) (int, error)
		// Get returns an attachment by it's UUID.
		Get(uuid string) (*Attachment, error)
		// GetFS returns the attachment file system.
		GetFS() fs.FS
		// GetSizeAll returns the size of all attachments in bytes.
		GetSizeAll() (int64, error)
		// LoadData sets the Data field of an attachment.
		LoadData(att *Attachment) error
		// ListByMessage returns a list of attachments for a message.
		ListByMessage(msg *Message) ([]Attachment, error)
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
