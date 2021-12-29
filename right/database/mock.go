package database

import "github.com/ItsNotGoodName/smtpbridge/app"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (db *Mock) CreateMessage(msg *app.Message) error {
	return nil
}

func (db *Mock) GetMessage(uuid string) (*app.Message, error) {
	return nil, app.ErrNotImplemented
}

func (db *Mock) UpdateMessage(msg *app.Message, updateFN func(msg *app.Message) (*app.Message, error)) (*app.Message, error) {
	return updateFN(msg)
}

func (db *Mock) CreateAttachment(att *app.Attachment) error {
	return nil
}

func (db *Mock) GetAttachment(uuid string) (*app.Attachment, error) {
	return nil, app.ErrNotImplemented
}

func (db *Mock) LoadAttachment(msg *app.Message) error {
	return app.ErrNotImplemented
}
