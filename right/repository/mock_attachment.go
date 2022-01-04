package repository

import (
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type AttachmentMock struct{}

func NewAttachmentMock() *AttachmentMock {
	return &AttachmentMock{}
}

func (AttachmentMock) Create(att *core.Attachment) error {
	return nil
}

func (AttachmentMock) Get(uuid string) (*core.Attachment, error) {
	return nil, core.ErrAttachmentNotFound
}

func (AttachmentMock) GetData(att *core.Attachment) ([]byte, error) {
	return nil, core.ErrAttachmentNotFound
}

func (a AttachmentMock) GetFS() fs.FS {
	return a
}

func (a AttachmentMock) GetSizeAll() (int64, error) {
	return 0, nil
}

func (AttachmentMock) ListByMessage(msg *core.Message) ([]core.Attachment, error) {
	return []core.Attachment{}, nil
}

func (AttachmentMock) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (AttachmentMock) DeleteData(att *core.Attachment) error {
	return nil
}
