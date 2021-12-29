package database

import "github.com/ItsNotGoodName/smtpbridge/app"

type DB struct {
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) CreateMessage(msg *app.Message) error {
	return app.ErrNotImplemented
}

func (db *DB) GetMessage(uuid string) (*app.Message, error) {
	return nil, app.ErrNotImplemented
}

func (db *DB) UpdateMessage(msg *app.Message, updateFN func(msg *app.Message) (*app.Message, error)) (*app.Message, error) {
	return nil, app.ErrNotImplemented
}

func (db *DB) CreateAttachment(att *app.Attachment) error {
	return app.ErrNotImplemented
}

func (db *DB) GetAttachment(uuid string) (*app.Attachment, error) {
	return nil, app.ErrNotImplemented
}

func (db *DB) LoadAttachment(msg *app.Message) error {
	return app.ErrNotImplemented
}
