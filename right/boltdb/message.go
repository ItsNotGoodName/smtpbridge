package boltdb

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/entity"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type messageModel struct {
	ID        int64               `json:"id" storm:"id,increment"`
	From      string              `json:"from"`
	To        map[string]struct{} `json:"to"`
	Subject   string              `json:"subject"`
	Text      string              `json:"text"`
	CreatedAt time.Time           `json:"created_at"`
	Processed bool                `json:"processed"`
}

func convertMessageM(msg *messageModel) *message.Message {
	return &message.Message{
		ID:        msg.ID,
		From:      msg.From,
		To:        msg.To,
		Subject:   msg.Subject,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
		Processed: msg.Processed,
	}
}

func convertMessageD(msg *message.Message) *messageModel {
	return &messageModel{
		ID:        msg.ID,
		From:      msg.From,
		To:        msg.To,
		Subject:   msg.Subject,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
		Processed: msg.Processed,
	}
}

type Message struct {
	db             *storm.DB
	dataRepository attachment.DataRepository
}

func NewMessage(db *Database, dataRepository attachment.DataRepository) *Message {
	return &Message{
		db:             db.db,
		dataRepository: dataRepository,
	}
}

func (m *Message) Create(ctx context.Context, msg *message.Message) error {
	msgD := convertMessageD(msg)
	if err := m.db.Save(msgD); err != nil {
		return err
	}

	msg.ID = msgD.ID
	return nil
}

func (m *Message) Get(ctx context.Context, id int64) (*message.Message, error) {
	var msgM messageModel
	err := m.db.One("ID", id, &msgM)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, message.ErrNotFound
		}
		return nil, err
	}

	return convertMessageM(&msgM), nil
}

func (m *Message) Update(ctx context.Context, msg *message.Message, updateFN func(msg *message.Message) (*message.Message, error)) error {
	tx, err := m.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingMSGM messageModel
	if err := tx.One("ID", msg.ID, &existingMSGM); err != nil {
		if err == storm.ErrNotFound {
			return message.ErrNotFound
		}
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

func (m *Message) Count(ctx context.Context) (int, error) {
	count, err := m.db.Select().Count(&messageModel{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, err
}

func (m *Message) List(ctx context.Context, param *message.ListParam) error {
	var (
		q1 storm.Query
		q2 storm.Query
	)
	if param.Cursor.Ascending {
		q1 = m.db.Select(q.Gte("ID", param.Cursor.Cursor)).OrderBy("ID").Limit(param.Cursor.Limit + 1)
		q2 = m.db.Select(q.Lt("ID", param.Cursor.Cursor)).OrderBy("ID").Limit(param.Cursor.Limit).Reverse()
	} else {
		q1 = m.db.Select(q.Lte("ID", param.Cursor.Cursor)).OrderBy("ID").Limit(param.Cursor.Limit + 1).Reverse()
		q2 = m.db.Select(q.Gt("ID", param.Cursor.Cursor)).OrderBy("ID").Limit(param.Cursor.Limit)
	}

	var msgsM []messageModel
	if err := q1.Find(&msgsM); err != nil && err != storm.ErrNotFound {
		return err
	}
	if len(msgsM) == param.Cursor.Limit+1 {
		param.Cursor.SetNextCursor(msgsM[param.Cursor.Limit].ID)
		msgsM = msgsM[:param.Cursor.Limit]
	}

	var msgsMB []messageModel
	if err := q2.Find(&msgsMB); err != nil && err != storm.ErrNotFound {
		return err
	}
	if length := len(msgsMB); length > 0 {
		param.Cursor.SetBackCursor(msgsMB[length-1].ID)
	}

	var msgs []message.Message
	for _, msgM := range msgsM {
		msgs = append(msgs, *convertMessageM(&msgM))
	}
	param.Messages = msgs

	return nil
}

func (m *Message) Delete(ctx context.Context, msg *message.Message) error {
	tx, err := m.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := tx.Select(q.Eq("MessageID", msg.ID))

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

	// Delete events
	err = tx.Select(q.Eq("EntityID", msg.ID), q.Eq("Entity", entity.Message)).Delete(&eventModel{})
	if err != nil {
		return err
	}

	// Delete message
	if err := tx.DeleteStruct(convertMessageD(msg)); err != nil {
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	// TODO: Prevent orphaned attachment data

	// Delete attachment's data
	for _, attM := range attsM {
		if err := m.dataRepository.Delete(ctx, convertAttachmentM(&attM)); err != nil {
			log.Println("boltdb.Message.Delete: could not delete attachment file:", err)
		}
	}

	return nil
}
