package database

import (
	"log"
	"os"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func (db *DB) CreateMessage(msg *domain.Message) error {
	return db.db.Save(msg)
}

func (db *DB) GetMessage(uuid string) (*domain.Message, error) {
	var msg domain.Message
	err := db.db.One("UUID", uuid, msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (db *DB) UpdateMessage(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingMSG domain.Message
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

func (db *DB) GetMessages(limit, offset int) ([]domain.Message, error) {
	var msgs []domain.Message
	err := db.db.Select().OrderBy("CreatedAt").Limit(limit).Skip(offset).Reverse().Find(&msgs)
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}

	return msgs, nil
}

func (db *DB) CountMessages() (int, error) {
	count, err := db.db.Count(&domain.Message{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, err
}

func (db *DB) DeleteMessage(msg *domain.Message) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := tx.Select(q.Eq("MessageUUID", msg.UUID))

	// List attachments
	var atts []domain.Attachment
	err = query.Find(&atts)
	if err != storm.ErrNotFound {
		if err != nil {
			return err
		}

		// Delete attachments
		err = query.Delete(&domain.Attachment{})
		if err != nil {
			return err
		}
	}

	// Delete message
	if err := tx.DeleteStruct(msg); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, att := range atts {
		if err := os.Remove(db.getAttachmentPath(&att)); err != nil {
			log.Println("database.DB.DeleteMessage: could not delete attachment file:", err)
		}
	}

	return nil
}
