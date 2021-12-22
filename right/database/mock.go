package database

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct {
	db     map[string]app.Message
	dbLock sync.RWMutex
}

func NewMock() *Mock {
	return &Mock{
		db:     make(map[string]app.Message),
		dbLock: sync.RWMutex{},
	}
}

func (db *Mock) Create(msg *app.Message) error {
	db.dbLock.Lock()
	db.db[msg.UUID] = *msg
	db.dbLock.Unlock()

	return nil
}

func (db *Mock) Update(msg *app.Message) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	_, ok := db.db[msg.UUID]
	if !ok {
		return app.ErrMessageNotFound
	}

	db.db[msg.UUID] = *msg

	return nil
}
