package repository

import (
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type AttachmentMock struct{}

func NewAttachmentMock() *AttachmentMock {
	return &AttachmentMock{}
}

func (AttachmentMock) Create(att *domain.Attachment) error {
	return nil
}

func (AttachmentMock) Get(uuid string) (*domain.Attachment, error) {
	return nil, domain.ErrNotImplemented
}

func (AttachmentMock) GetData(att *domain.Attachment) ([]byte, error) {
	return nil, domain.ErrNotImplemented
}

func (a AttachmentMock) GetFS() fs.FS {
	return a
}

func (AttachmentMock) ListByMessage(msg *domain.Message) ([]domain.Attachment, error) {
	return []domain.Attachment{}, nil
}

func (AttachmentMock) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (AttachmentMock) DeleteData(att *domain.Attachment) error {
	return nil
}
