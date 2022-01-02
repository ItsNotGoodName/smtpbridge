package repository

import (
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type messageModel struct {
	UUID      string              `storm:"id"` // UUID of the message.
	From      string              ``           // From is the email address of the sender.
	To        map[string]struct{} ``           // To is the email addresses of the recipients.
	Subject   string              ``           // Subject of the message.
	Text      string              ``           // Text is the message body.
	Status    domain.Status       ``           // Status is the status of the message.
	CreatedAt time.Time           ``           // Time message was received.
}

func convertMessageM(msg *messageModel) *domain.Message {
	return &domain.Message{
		UUID:      msg.UUID,
		From:      msg.From,
		To:        msg.To,
		Subject:   msg.Subject,
		Text:      msg.Text,
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
	}
}

func convertMessageD(msg *domain.Message) *messageModel {
	return &messageModel{
		UUID:      msg.UUID,
		From:      msg.From,
		To:        msg.To,
		Subject:   msg.Subject,
		Text:      msg.Text,
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
	}
}

type Message struct {
	db             *storm.DB
	attachmentREPO *Attachment
}

func NewMessage(db *storm.DB, attachmentREPO *Attachment) *Message {
	return &Message{
		db:             db,
		attachmentREPO: attachmentREPO,
	}
}

func (m *Message) Create(msg *domain.Message) error {
	return m.db.Save(convertMessageD(msg))
}

func (m *Message) Get(uuid string) (*domain.Message, error) {
	var msgM messageModel
	err := m.db.One("UUID", uuid, &msgM)
	if err != nil {
		return nil, err
	}

	return convertMessageM(&msgM), nil
}

func (m *Message) Update(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	tx, err := m.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingMSGM messageModel
	if err := tx.One("UUID", msg.UUID, &existingMSGM); err != nil {
		return err
	}

	updatedMSG, err := updateFN(convertMessageM(&existingMSGM))
	if err != nil {
		return err
	}

	err = tx.Save(convertMessageD(updatedMSG))
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Message) List(limit, offset int) ([]domain.Message, error) {
	var msgsM []messageModel
	err := m.db.Select().OrderBy("CreatedAt").Limit(limit).Skip(offset).Reverse().Find(&msgsM)
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}

	var msgs []domain.Message
	for _, msgM := range msgsM {
		msgs = append(msgs, *convertMessageM(&msgM))
	}

	return msgs, nil
}

func (m *Message) Count() (int, error) {
	count, err := m.db.Count(&messageModel{})
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
	var attsM []attachmentModel
	err = query.Find(&attsM)
	if err != storm.ErrNotFound {
		if err != nil {
			return err
		}

		// Delete attachments
		err = query.Delete(&attachmentModel{})
		if err != nil {
			return err
		}
	}

	// Delete message
	if err := tx.DeleteStruct(convertMessageD(msg)); err != nil {
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	// Delete attachment's data
	for _, attM := range attsM {
		if err := m.attachmentREPO.DeleteData(convertAttachmentM(&attM)); err != nil {
			log.Println("database.DB.DeleteMessage: could not delete attachment file:", err)
		}
	}

	return nil
}
