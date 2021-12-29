package database

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/asdine/storm"
)

type DB struct {
	db *storm.DB
}

func NewDB(dbPath string) *DB {
	db, err := storm.Open(dbPath)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	return &DB{
		db: db,
	}
}

func (db *DB) CreateMessage(msg *app.Message) error {
	return db.db.Save(msg)
}

func (db *DB) GetMessage(uuid string) (*app.Message, error) {
	var msg app.Message
	err := db.db.One("UUID", uuid, msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (db *DB) UpdateMessage(msg *app.Message, updateFN func(msg *app.Message) (*app.Message, error)) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingMSG app.Message
	if err := tx.One("UUID", msg.UUID, &existingMSG); err != nil {
		return err
	}

	updatedMSG, err := updateFN(&existingMSG)
	if err != nil {
		return err
	}

	err = tx.Save(updatedMSG)
	if err != nil {
		return err
	}

	return tx.Commit()
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
