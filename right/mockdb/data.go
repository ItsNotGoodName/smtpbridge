package mockdb

import (
	"context"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
)

type Data struct {}

func NewData() *Data {
	return &Data{}
}

func (d *Data) Create(ctx context.Context, att *attachment.Attachment) error {
	return nil
}

func (d *Data) FS() fs.FS {
	return d
}

func (d *Data) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (d *Data) Delete(ctx context.Context, att *attachment.Attachment) error {
	return nil
}

func (d *Data) Load(ctx context.Context, att *attachment.Attachment) error {
	return attachment.ErrNotFound
}

func (d *Data) Size(ctx context.Context) (int64, error) {
	return 0, nil
}
