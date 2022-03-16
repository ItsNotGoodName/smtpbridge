package boltdb

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type attachmentModel struct {
	ID        int64  `json:"id" storm:"id,increment"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	MessageID int64  `json:"message_id" storm:"index"`
}

func convertAttachmentD(att *attachment.Attachment) *attachmentModel {
	return &attachmentModel{
		ID:        att.ID,
		Name:      att.Name,
		Type:      string(att.Type),
		MessageID: att.MessageID,
	}
}

func convertAttachmentM(attM *attachmentModel) *attachment.Attachment {
	return &attachment.Attachment{
		ID:        attM.ID,
		Name:      attM.Name,
		Type:      attachment.Type(attM.Type),
		MessageID: attM.MessageID,
	}
}

type Attachment struct {
	db             *storm.DB
	dataRepository attachment.DataRepository
}

func NewAttachment(db *Database, dataRepository attachment.DataRepository) *Attachment {
	return &Attachment{
		db:             db.db,
		dataRepository: dataRepository,
	}
}

func (a *Attachment) Create(ctx context.Context, att *attachment.Attachment) error {
	attD := convertAttachmentD(att)
	if err := a.db.Save(attD); err != nil {
		return err
	}

	att.ID = attD.ID

	return a.dataRepository.Create(ctx, att)
}

func (a *Attachment) Count(ctx context.Context) (int, error) {
	count, err := a.db.Count(&attachmentModel{})
	if err != nil {
		if err == storm.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (a *Attachment) CountByMessage(ctx context.Context, msg *message.Message) (int, error) {
	count, err := a.db.Select(q.Eq("MessageID", msg.ID)).Count(&attachmentModel{})
	if err != nil {
		if err == storm.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (a *Attachment) Get(ctx context.Context, id int64) (*attachment.Attachment, error) {
	var attM *attachmentModel
	err := a.db.One("ID", id, attM)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, attachment.ErrNotFound
		}
		return nil, err
	}

	return convertAttachmentM(attM), nil
}

func (a *Attachment) ListByMessage(ctx context.Context, msg *message.Message) ([]attachment.Attachment, error) {
	var attsM []attachmentModel
	err := a.db.Select(q.Eq("MessageID", msg.ID)).Find(&attsM)
	if err != nil {
		if err == storm.ErrNotFound {
			return []attachment.Attachment{}, nil
		}
		return nil, err
	}

	var atts []attachment.Attachment
	for _, attM := range attsM {
		atts = append(atts, *convertAttachmentM(&attM))
	}

	return atts, nil
}
