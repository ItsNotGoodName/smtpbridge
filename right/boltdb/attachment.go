package boltdb

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

func (d Database) CountAttachment(ctx context.Context) (int, error) {
	return countAttachments(d.db)
}

func (d Database) ListAttachment(ctx context.Context, offset, limit int, ascending bool) ([]envelope.Attachment, int, error) {
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

	// Get attachments
	var attsM []attachmentModel
	if err := query.Find(&attsM); err != nil && err != storm.ErrNotFound {
		return nil, 0, err
	}
	var atts []envelope.Attachment
	for _, attM := range attsM {
		atts = append(atts, *attachmentMC(&attM))
	}

	// Get attachments count
	count, err := countAttachments(tx)
	if err != nil {
		return nil, 0, err
	}

	return atts, count, nil
}

func listAttachment(tx storm.Node, msgID int64) ([]envelope.Attachment, error) {
	// Get attachments
	var attsM []attachmentModel
	if err := tx.Select(q.Eq("MessageID", msgID)).Find(&attsM); err != nil && err != storm.ErrNotFound {
		return nil, err
	}
	var atts []envelope.Attachment
	for _, attM := range attsM {
		atts = append(atts, *attachmentMC(&attM))
	}

	return atts, nil
}

func countAttachments(tx storm.Node) (int, error) {
	return count(tx, &attachmentModel{})
}
