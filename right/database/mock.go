package database

import (
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (Mock) CreateMessage(msg *domain.Message) error {
	return nil
}

func (Mock) GetMessage(uuid string) (*domain.Message, error) {
	return nil, domain.ErrNotImplemented
}

func (Mock) UpdateMessage(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	_, err := updateFN(msg)
	return err
}

func (Mock) CreateAttachment(att *domain.Attachment) error {
	return nil
}

func (Mock) GetAttachment(uuid string) (*domain.Attachment, error) {
	return nil, domain.ErrNotImplemented
}

func (Mock) LoadAttachment(msg *domain.Message) error {
	return nil
}

func (Mock) DeleteMessage(msg *domain.Message) error {
	return nil
}

func (Mock) GetMessages(limit, offset int) ([]domain.Message, error) {
	return []domain.Message{}, nil
}

func (Mock) GetAttachmentData(att *domain.Attachment) ([]byte, error) {
	return nil, domain.ErrNotImplemented
}

func (m Mock) GetAttachmentFS() fs.FS {
	return m
}

func (Mock) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (Mock) GetAttachments(msg *domain.Message) ([]domain.Attachment, error) {
	return []domain.Attachment{}, nil
}
