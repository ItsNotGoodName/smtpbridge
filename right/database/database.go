package database

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Database struct {
}

func New() *Database {
	return &Database{}
}

func (d *Database) Create(msg *app.Message) error {
	log.Println("Create message:", msg.UUID)
	return nil
}

func (d *Database) Update(msg *app.Message) error {
	log.Println("Update message:", msg.UUID)
	return nil
}
