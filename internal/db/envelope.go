package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/dbgen/model"
	. "github.com/ItsNotGoodName/smtpbridge/internal/dbgen/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/uptrace/bun"
)

func EnvelopeDeleteAll(cc core.Context) error {
	_, err := cc.DB.NewDelete().Model(&envelope.Message{}).Where("1=1").Exec(cc)
	return err
}

func EnvelopeDelete(cc core.Context, id int64) error {
	_, err := cc.DB.NewDelete().Model(&envelope.Message{}).Where("id = ?", id).Exec(cc)
	return err
}

func EnvelopeCreate(cc core.Context, msg *envelope.Message, atts []*envelope.Attachment) (int64, []int64, error) {
	err := cc.DB.RunInTx(cc, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(msg).Exec(ctx)
		if err != nil {
			return err
		}

		for _, att := range atts {
			att.MessageID = msg.ID
			_, err := tx.NewInsert().Model(att).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return err
	})
	if err != nil {
		return 0, nil, err
	}

	attIDS := make([]int64, len(atts))
	for i, att := range atts {
		attIDS[i] = att.ID
	}

	return msg.ID, attIDS, nil
}

func EnvelopeMessageList(cc core.Context, page pagination.Page, filter envelope.MessageFilter) (envelope.MessageListResult, error) {
	var msgs []*envelope.Message
	q := cc.DB.NewSelect().Model(&msgs).Limit(page.Limit()).Offset(page.Offset())

	if filter.Ascending {
		q = q.Order("created_at ASC")
	} else {
		q = q.Order("created_at DESC")
	}
	if filter.Search != "" {
		if filter.SearchSubject {
			q = q.WhereOr("subject LIKE ?", "%"+filter.Search+"%")
		}
		if filter.SearchText {
			q = q.WhereOr("text LIKE ?", "%"+filter.Search+"%")
		}
	}

	ctx := cc

	err := q.Scan(ctx, &msgs)
	if err != nil {
		return envelope.MessageListResult{}, err
	}

	count, err := q.Count(ctx)
	if err != nil {
		return envelope.MessageListResult{}, err
	}

	return envelope.MessageListResult{
		Messages:   msgs,
		PageResult: pagination.NewPageResult(page, count),
	}, nil
}

func EnvelopeGet(cc core.Context, id int64) (envelope.Envelope, error) {
	ctx := cc
	msg := &envelope.Message{}
	err := cc.DB.NewSelect().Model(msg).Where("id = ?", id).Scan(ctx, msg)
	if err != nil {
		return envelope.Envelope{}, err
	}

	atts := []*envelope.Attachment{}
	err = cc.DB.NewSelect().Model(&atts).Where("message_id = ?", id).Scan(ctx, &atts)
	if err != nil {
		return envelope.Envelope{}, err
	}

	return envelope.Envelope{
		Message:     msg,
		Attachments: atts,
	}, nil
}

func EnvelopeMessageHTMLGet(cc core.Context, id int64) (string, error) {
	res := model.Messages{}
	err := Messages.
		SELECT(Messages.HTML).
		WHERE(Messages.ID.EQ(Int64(id))).
		QueryContext(cc, cc, &res)
	return res.HTML, err
}

func EnvelopeCount(cc core.Context) (int, error) {
	return cc.DB.NewSelect().Model(&envelope.Message{}).Count(cc)
}

func EnvelopeAttachmentCount(cc core.Context) (int, error) {
	return cc.DB.NewSelect().Model(&envelope.Attachment{}).Where("message_id IS NOT NULL").Count(cc)
}

func EnvelopeAttachmentList(cc core.Context, page pagination.Page, filter envelope.AttachmentFilter) (envelope.AttachmentListResult, error) {
	var atts []*envelope.Attachment
	q := cc.DB.NewSelect().Model(&atts).Limit(page.Limit()).Offset(page.Offset()).Where("message_id IS NOT NULL")

	// Filter
	if filter.Ascending {
		q = q.Order("id ASC")
	} else {
		q = q.Order("id DESC")
	}

	ctx := cc

	// Scan
	err := q.Scan(ctx, &atts)
	if err != nil {
		return envelope.AttachmentListResult{}, err
	}

	// Count
	count, err := q.Count(ctx)
	if err != nil {
		return envelope.AttachmentListResult{}, err
	}

	return envelope.AttachmentListResult{
		Attachments: atts,
		PageResult:  pagination.NewPageResult(page, count),
	}, nil
}

func EnvelopeAttachmentListOrphan(cc core.Context, limit int) ([]*envelope.Attachment, error) {
	var atts []*envelope.Attachment
	err := cc.DB.NewSelect().Model(&atts).Limit(limit).Where("message_id IS NULL").Scan(cc)
	return atts, err
}

func EnvelopeDeleteUntilCount(cc core.Context, keep int, olderThan time.Time) (int64, error) {
	// DELETE FROM messages WHERE id NOT IN ( SELECT id FROM messages ORDER BY id DESC LIMIT ?) AND messages.created_at < ?;

	res, err := Messages.
		DELETE().
		WHERE(Messages.ID.NOT_IN(Messages.SELECT(Messages.ID).ORDER_BY(Messages.ID).LIMIT(int64(keep))).
			AND(Messages.CreatedAt.LT(DATETIME(olderThan))),
		).ExecContext(cc, cc.DB)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func EnvelopeDeleteOlderThan(cc core.Context, olderThan time.Time) (int64, error) {
	res, err := Messages.DELETE().WHERE(Messages.CreatedAt.LT(DATETIME(olderThan))).ExecContext(cc, cc.DB)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

func EnvelopeAttachmentDelete(cc core.Context, att *envelope.Attachment) error {
	err := files.DeleteFile(cc, att)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	_, err = cc.DB.NewDelete().Model(&envelope.Attachment{}).Where("id = ?", att.ID).Exec(cc)

	return err
}
