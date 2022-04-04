package attachment

import (
	"context"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core/message"
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

	ServiceData interface {
		// Remote true if FS & URL is remote.
		Remote() bool
		// URI returns the URI prefix if it is not remote.
		URI() string
		// FS returns the attachment file system.
		FS() fs.FS
		// URL returns the attachment's data URL.
		URL(att *Attachment) string
		// File returns the attachment file name.
		File(att *Attachment) string
		// Create creates a new attachment data.
		Create(ctx context.Context, att *Attachment) error
		// Load loads the attachment data.
		Load(ctx context.Context, att *Attachment) error
		// Delete deletes the attachment data.
		Delete(ctx context.Context, att *Attachment) error
		// Size returns the size of all attachments in bytes.
		Size(ctx context.Context) (int64, error)
	}

	RepositoryData interface {
		// Remote true if fs & URL is remote.
		Remote() bool
		// FS returns the attachment file system.
		FS() fs.FS
		// File returns the attachment file name.
		File(att *Attachment) string
		// URL returns the attachment url.
		URL(att *Attachment) string
		// Create creates a new attachment data.
		Create(ctx context.Context, att *Attachment, data []byte) error
		// Get returns an attachment data by name.
		Get(ctx context.Context, att *Attachment) ([]byte, error)
		// Delete deletes the attachment file.
		Delete(ctx context.Context, att *Attachment) error
		// Size returns the size of all attachments in bytes.
		Size(ctx context.Context) (int64, error)
	}

	Type string
)
