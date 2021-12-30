package database

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type DB struct {
	db     *storm.DB
	attDir string
}

func NewDB(dbFile, attDir string) *DB {
	db, err := storm.Open(dbFile)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	err = os.MkdirAll(attDir, 0755)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	return &DB{
		db:     db,
		attDir: attDir,
	}
}

// getAttachmentPath returns the path to the attachment file on the file system.
func (db *DB) getAttachmentPath(att *app.Attachment) string {
	return path.Join(db.attDir, db.GetAttachmentFile(att))
}

func (db *DB) GetAttachmentFile(att *app.Attachment) string {
	return fmt.Sprintf("%s.%s", att.UUID, att.Type)
}

func (db *DB) GetAttachmentFS() fs.FS {
	return os.DirFS(db.attDir)
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

func (db *DB) GetMessages(limit, offset int) ([]app.Message, error) {
	var msgs []app.Message
	err := db.db.All(&msgs, storm.Limit(limit), storm.Skip(offset), storm.Reverse())
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (db *DB) CreateAttachment(att *app.Attachment) error {
	err := db.db.Save(att)
	if err != nil {
		return err
	}

	return os.WriteFile(db.getAttachmentPath(att), att.Data, 0644)
}

func (db *DB) GetAttachment(uuid string) (*app.Attachment, error) {
	var att app.Attachment
	err := db.db.One("UUID", uuid, att)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (db *DB) GetAttachmentData(att *app.Attachment) ([]byte, error) {
	data, err := os.ReadFile(db.getAttachmentPath(att))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *DB) GetAttachments(msg *app.Message) ([]app.Attachment, error) {
	var atts []app.Attachment
	err := db.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&atts)
	if err != nil {
		return nil, err
	}

	return atts, nil
}
