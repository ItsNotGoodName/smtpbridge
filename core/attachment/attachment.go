package attachment

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

const (
	TypePNG  Type = "png"
	TypeJPEG Type = "jpeg"
)

var (
	ErrInvalid       = fmt.Errorf("invalid attachment")
	ErrDataNotLoaded = fmt.Errorf("attachment data not loaded")
	ErrNotFound      = fmt.Errorf("attachment not found")
)

type (
	Attachment struct {
		ID        int64  // ID is the attachment ID.
		Name      string // Name is the attachment name.
		Type      Type   // Type is the attachment type.
		MessageID int64  // MessageID is the message ID that the attachment belongs to.
		data      []byte // data is the attachment data.
	}

	Param struct {
		ID      int64            // ID of the message.
		Name    string           // Name of the message.
		Message *message.Message // Message of the message.
		Data    []byte           // data of the message.
	}

	Service interface {
		Create(ctx context.Context, param *Param) (*Attachment, error)
	}

	DataRepository interface {
		// FS returns the attachment file system.
		FS() fs.FS
		// Size returns the size of all attachments in bytes.
		Size(ctx context.Context) (int64, error)
		// Create creates a new attachment data.
		Create(ctx context.Context, att *Attachment) error
		// Load loads the attachment data.
		Load(ctx context.Context, att *Attachment) error
		// Delete deletes the attachment data.
		Delete(ctx context.Context, att *Attachment) error
	}

	Repository interface {
		// Create creates a new attachment.
		Create(ctx context.Context, att *Attachment) error
		// Count returns the number of attachments.
		Count(ctx context.Context) (int, error)
		// CountByMessage returns the number of attachments for a message.
		CountByMessage(ctx context.Context, msg *message.Message) (int, error)
		// Get returns an attachment by ID.
		Get(ctx context.Context, id int64) (*Attachment, error)
		// ListByMessage returns a list of attachments for a message.
		ListByMessage(ctx context.Context, msg *message.Message) ([]Attachment, error)
	}

	Type string
)

func New(param *Param) (*Attachment, error) {
	attType, err := DataType(param.Data)
	if err != nil {
		return nil, err
	}

	return &Attachment{
		ID:        param.ID,
		Name:      param.Name,
		Type:      attType,
		MessageID: param.Message.ID,
		data:      param.Data,
	}, nil
}

// DataType returns the type of the attachment data.
func DataType(data []byte) (Type, error) {
	if len(data) == 0 {
		return "", ErrDataNotLoaded
	}

	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/png":
		return TypePNG, nil
	case "image/jpeg":
		return TypeJPEG, nil
	default:
		return "", fmt.Errorf("%s: %v", contentType, ErrInvalid)
	}
}

// File returns the full attachment file name.
func (a *Attachment) File() string {
	return fmt.Sprintf("%d.%s", a.ID, a.Type)
}

// GetData returns the attachment data.
func (a *Attachment) GetData() ([]byte, error) {
	if a.data == nil {
		return nil, ErrDataNotLoaded
	}

	return a.data, nil
}

// GetData returns the attachment data.
func (a *Attachment) SetData(data []byte) error {
	_, err := DataType(data)
	if err != nil {
		return err
	}

	a.data = data
	return nil
}
