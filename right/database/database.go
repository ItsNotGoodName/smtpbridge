package database

import (
	"log"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type DB struct {
	db      *storm.DB
	attPath string
}

func NewDB(dbPath, attPath string) *DB {
	db, err := storm.Open(dbPath)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	err = os.MkdirAll(attPath, 0755)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	return &DB{
		db:      db,
		attPath: attPath,
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
	err := db.db.Save(att)
	if err != nil {
		return err
	}

	file := path.Join(db.attPath, att.UUID+att.EXT())

	return os.WriteFile(file, att.Data, 0644)
}

func (db *DB) GetAttachment(uuid string) (*app.Attachment, error) {
	var att app.Attachment
	err := db.db.One("UUID", uuid, att)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path.Join(db.attPath, att.UUID+att.EXT()))
	if err != nil {
		return nil, err
	}

	att.Data = data

	return &att, nil
}

func (db *DB) LoadAttachment(msg *app.Message) error {
	var atts []app.Attachment
	err := db.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&atts)
	if err != nil {
		return err
	}

	msg.Attachments = atts

	return nil
}
