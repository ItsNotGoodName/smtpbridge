package repository

import "github.com/ItsNotGoodName/smtpbridge/domain"

type MessageMock struct{}

func NewMessageMock() *MessageMock {
	return &MessageMock{}
}

func (MessageMock) Create(msg *domain.Message) error {
	return nil
}

func (MessageMock) Count() (int, error) {
	return 0, nil
}

func (MessageMock) Get(uuid string) (*domain.Message, error) {
	return nil, domain.ErrNotImplemented
}

func (MessageMock) List(limit, offset int) ([]domain.Message, error) {
	return []domain.Message{}, nil
}

func (MessageMock) Update(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	return nil
}

func (MessageMock) Delete(msg *domain.Message) error {
	return nil
}
