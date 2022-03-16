package mockdb

import (
	"context"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type Attachment struct{}

func NewAttachment() *Attachment {
	return &Attachment{}
}

func (Attachment) Create(ctx context.Context, att *attachment.Attachment) error {
	return nil
}

func (Attachment) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (Attachment) CountByMessage(ctx context.Context, msg *message.Message) (int, error) {
	return 0, nil
}

func (a Attachment) FS() fs.FS {
	return a
}

func (Attachment) Get(ctx context.Context, id int64) (*attachment.Attachment, error) {
	return nil, attachment.ErrNotFound
}

func (Attachment) ListByMessage(ctx context.Context, msg *message.Message) ([]attachment.Attachment, error) {
	return []attachment.Attachment{}, nil
}

func (Attachment) LoadData(ctx context.Context, att *attachment.Attachment) error {
	return nil
}

func (Attachment) Size(ctx context.Context) (int64, error) {
	return 0, nil
}

func (Attachment) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
