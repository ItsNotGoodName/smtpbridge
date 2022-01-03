package repository

import "github.com/ItsNotGoodName/smtpbridge/core"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (Mock) Close() error {
	return nil
}

func (Mock) AttachmentRepository() core.AttachmentRepositoryPort {
	return NewAttachmentMock()
}

func (Mock) MessageRepository() core.MessageRepositoryPort {
	return NewMessageMock()
}
