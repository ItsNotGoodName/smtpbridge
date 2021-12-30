package database

import "github.com/ItsNotGoodName/smtpbridge/domain"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (db *Mock) CreateMessage(msg *domain.Message) error {
	return nil
}

func (db *Mock) GetMessage(uuid string) (*domain.Message, error) {
	return nil, domain.ErrNotImplemented
}

func (db *Mock) UpdateMessage(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	_, err := updateFN(msg)
	return err
}

func (db *Mock) CreateAttachment(att *domain.Attachment) error {
	return nil
}

func (db *Mock) GetAttachment(uuid string) (*domain.Attachment, error) {
	return nil, domain.ErrNotImplemented
}

func (db *Mock) LoadAttachment(msg *domain.Message) error {
	return domain.ErrNotImplemented
}
