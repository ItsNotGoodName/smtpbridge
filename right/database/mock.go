package database

import (
	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct {
}

func NewMock() *Mock {
	return &Mock{}
}

func (db *Mock) Create(msg *app.Message) error {
	return nil
}

func (db *Mock) Update(msg *app.Message) error {
	return nil
}
