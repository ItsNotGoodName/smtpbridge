package boltdb

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/asdine/storm/v3"
)

func (d Database) ListEnvelope(ctx context.Context, offset, limit int, ascending bool) ([]envelope.Envelope, int, error) {
	tx, err := d.db.Begin(false)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	// Query
	query := d.db.Select().OrderBy("ID").Limit(limit).Skip(offset)
	if !ascending {
		query = query.Reverse()
	}

	// Get messages and attachments
	var msgsM []messageModel
	if err := query.Find(&msgsM); err != nil && err != storm.ErrNotFound {
		return nil, 0, err
	}
	var envs []envelope.Envelope
	for _, msgM := range msgsM {
		msg := messageMC(&msgM)
		atts, err := listAttachment(tx, msg.ID)
		if err != nil {
			return nil, 0, err
		}

		envs = append(envs, envelope.Envelope{
			Message:     *msg,
			Attachments: atts,
		})
	}

	// Get envelopes count
	count, err := countMessages(tx)
	if err != nil {
		return nil, 0, err
	}

	return envs, count, nil
}

func (d Database) CreateEnvelope(ctx context.Context, msg *envelope.Message, atts []envelope.Attachment) (int64, error) {
	tx, err := d.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Create message
	msgM := messageCM(msg)
	if err := tx.Save(msgM); err != nil {
		return 0, err
	}
	msg.ID = msgM.ID

	// Create attachments
	for i := range atts {
		atts[i].MessageID = msg.ID

		attM := attachmentCM(&atts[i])
		if err := tx.Save(attM); err != nil {
			return 0, err
		}

		atts[i].ID = attM.ID
	}

	return msg.ID, tx.Commit()
}

func (d Database) CountEnvelope(ctx context.Context) (int, error) {
	return countMessages(d.db)
}

func (d Database) GetEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	return getEnvelope(d.db, msgID)
}

func (d Database) DeleteEnvelope(ctx context.Context, msgID int64, fn func(env *envelope.Envelope) error) error {
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get envelope
	env, err := getEnvelope(tx, msgID)
	if err != nil {
		return err
	}

	// Delete message
	if tx.DeleteStruct(messageCM(&env.Message)); err != nil {
		return err
	}

	// Delete attachments
	for _, att := range env.Attachments {
		tx.DeleteStruct(attachmentCM(&att))
	}

	// Run hook
	if err := fn(env); err != nil {
		return err
	}

	return tx.Commit()
}

func getEnvelope(tx storm.Node, msgID int64) (*envelope.Envelope, error) {
	// Get message
	var msgM messageModel
	err := tx.One("ID", msgID, &msgM)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, core.ErrMessageNotFound
		}
		return nil, err
	}

	// Get attachments
	atts, err := listAttachment(tx, msgID)
	if err != nil {
		return nil, err
	}

	return &envelope.Envelope{
		Message:     *messageMC(&msgM),
		Attachments: atts,
	}, nil
}

func countMessages(tx storm.Node) (int, error) {
	return count(tx, &messageModel{})
}
