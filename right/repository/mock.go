package repository

import "github.com/ItsNotGoodName/smtpbridge/domain"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (Mock) Close() error {
	return nil
}

func (Mock) AttachmentRepository() domain.AttachmentRepositoryPort {
	return NewAttachmentMock()
}

func (Mock) MessageRepository() domain.MessageRepositoryPort {
	return NewMessageMock()
}
