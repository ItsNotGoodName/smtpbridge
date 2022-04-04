package mockdb

import (
	"context"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
)

type Data struct{}

func NewData() *Data {
	return &Data{}
}

func (Data) Remote() bool {
	return true
}

func (Data) URL(*attachment.Attachment) string {
	return ""
}

func (Data) File(*attachment.Attachment) string {
	return ""
}

func (Data) Create(context.Context, *attachment.Attachment, []byte) error {
	return nil
}

func (Data) Get(context.Context, *attachment.Attachment) ([]byte, error) {
	return nil, attachment.ErrNotFound
}

func (d *Data) FS() fs.FS {
	return d
}

func (Data) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (Data) Delete(context.Context, *attachment.Attachment) error {
	return nil
}

func (Data) Size(ctx context.Context) (int64, error) {
	return 0, nil
}
