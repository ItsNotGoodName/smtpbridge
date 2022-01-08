package repository

import "github.com/ItsNotGoodName/smtpbridge/core"

type MessageMock struct{}

func NewMessageMock() *MessageMock {
	return &MessageMock{}
}

func (MessageMock) Create(msg *core.Message) error {
	return nil
}

func (MessageMock) Count(search *core.MessageParam) (int, error) {
	return 0, nil
}

func (MessageMock) Get(uuid string) (*core.Message, error) {
	return nil, core.ErrMessageNotFound
}

func (MessageMock) List(search *core.MessageParam) ([]core.Message, error) {
	return []core.Message{}, nil
}

func (MessageMock) Update(msg *core.Message, updateFN func(msg *core.Message) (*core.Message, error)) error {
	return nil
}

func (MessageMock) Delete(msg *core.Message) error {
	return nil
}
