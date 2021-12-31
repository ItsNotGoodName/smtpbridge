package database

import (
	"io/fs"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func (db *DB) CreateAttachment(att *domain.Attachment) error {
	err := db.db.Save(att)
	if err != nil {
		return err
	}

	return os.WriteFile(db.getAttachmentPath(att), att.Data, 0644)
}

// getAttachmentPath returns the path to the attachment file on the file system.
func (db *DB) getAttachmentPath(att *domain.Attachment) string {
	return path.Join(db.attDir, att.File())
}

func (db *DB) GetAttachmentFS() fs.FS {
	return db.fs
}

func (db *DB) GetAttachment(uuid string) (*domain.Attachment, error) {
	var att domain.Attachment
	err := db.db.One("UUID", uuid, att)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (db *DB) GetAttachmentData(att *domain.Attachment) ([]byte, error) {
	data, err := os.ReadFile(db.getAttachmentPath(att))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *DB) GetAttachments(msg *domain.Message) ([]domain.Attachment, error) {
	var atts []domain.Attachment
	err := db.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&atts)
	if err != nil {
		if err == storm.ErrNotFound {
			return []domain.Attachment{}, nil
		}
		return nil, err
	}

	return atts, nil
}
