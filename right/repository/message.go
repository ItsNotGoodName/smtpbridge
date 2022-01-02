package repository

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type Message struct {
	db             *storm.DB
	attachmentREPO domain.AttachmentRepositoryPort
}

func NewMessage(db *storm.DB, attachmentREPO domain.AttachmentRepositoryPort) *Message {
	return &Message{
		db:             db,
		attachmentREPO: attachmentREPO,
	}
}

func (m *Message) Create(msg *domain.Message) error {
	return m.db.Save(msg)
}

func (m *Message) Get(uuid string) (*domain.Message, error) {
	var msg domain.Message
	err := m.db.One("UUID", uuid, msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (m *Message) Update(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	tx, err := m.db.Begin(true)
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

func (m *Message) List(limit, offset int) ([]domain.Message, error) {
	var msgs []domain.Message
	err := m.db.Select().OrderBy("CreatedAt").Limit(limit).Skip(offset).Reverse().Find(&msgs)
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}

	return msgs, nil
}

func (m *Message) Count() (int, error) {
	count, err := m.db.Count(&domain.Message{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, err
}

func (m *Message) Delete(msg *domain.Message) error {
	tx, err := m.db.Begin(true)
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

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	// Delete attachment's data
	for _, att := range atts {
		if err := m.attachmentREPO.DeleteData(&att); err != nil {
			log.Println("database.DB.DeleteMessage: could not delete attachment file:", err)
		}
	}

	return nil
}
