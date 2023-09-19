package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo/orm"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	. "github.com/go-jet/jet/v2/sqlite"
)

var messagesPJRaw ProjectionList = ProjectionList{
	Messages.ID.AS("id"),
	Messages.UUID.AS("uuid"),
	Messages.From.AS("from"),
	Messages.To.AS("to"),
	Messages.Subject.AS("subject"),
	Messages.Text.AS("text"),
	Messages.HTML.AS("html"),
	Messages.Date.AS("date"),
	Messages.CreatedAt.AS("created_at"),
}

var messagePJ ProjectionList = messagesPJRaw.As("message")

var attachmentPJ ProjectionList = ProjectionList{
	Attachments.ID.AS("attachment.id"),
	Attachments.UUID.AS("uuid"),
	Attachments.MessageID.AS("attachment.message_id"),
	Attachments.Name.AS("attachment.name"),
	Attachments.Mime.AS("attachment.mime"),
	Attachments.Extension.AS("attachment.extension"),
}

func EnvelopeCreate(ctx context.Context, db database.Querier, msg models.Message, atts []models.Attachment) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := Messages.
		INSERT(
			Messages.UUID,
			Messages.From,
			Messages.To,
			Messages.Subject,
			Messages.Text,
			Messages.HTML,
			Messages.Date,
			Messages.CreatedAt,
		).
		MODEL(msg).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}
	msgID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if len(atts) != 0 {
		stmt := Attachments.INSERT(
			Attachments.UUID,
			Attachments.MessageID,
			Attachments.Name,
			Attachments.Mime,
			Attachments.Extension,
		)

		for _, att := range atts {
			att.MessageID = msgID
			stmt = stmt.MODEL(att)
		}

		res, err = stmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return msgID, nil
}

func EnvelopeList(ctx context.Context, db database.Querier, page pagination.Page, req models.DTOEnvelopeListRequest) (models.DTOEnvelopeListResult, error) {
	subQuery := Messages.SELECT(
		messagesPJRaw,
		COUNT(Raw("*")).OVER().AS("count"),
	)
	// Order
	if req.Ascending {
		if req.Order == models.DTOEnvelopeFieldSubject {
			subQuery = subQuery.ORDER_BY(Messages.Subject.ASC())
		} else if req.Order == models.DTOEnvelopeFieldFrom {
			subQuery = subQuery.ORDER_BY(Messages.From.ASC())
		} else {
			subQuery = subQuery.ORDER_BY(Messages.ID.ASC())
		}
	} else {
		if req.Order == models.DTOEnvelopeFieldSubject {
			subQuery = subQuery.ORDER_BY(Messages.Subject.DESC())
		} else if req.Order == models.DTOEnvelopeFieldFrom {
			subQuery = subQuery.ORDER_BY(Messages.From.DESC())
		} else {
			subQuery = subQuery.ORDER_BY(Messages.ID.DESC())
		}
	}
	// Filter
	if req.Search != "" {
		var exp []BoolExpression
		if req.SearchText {
			exp = append(exp, Messages.Text.LIKE(RawString("?", map[string]interface{}{"?": "%" + req.Search + "%"})))
		}
		if req.SearchSubject {
			exp = append(exp, Messages.Subject.LIKE(RawString("?", map[string]interface{}{"?": "%" + req.Search + "%"})))
		}
		if len(exp) > 0 {
			subQuery = subQuery.WHERE(OR(exp...))
		} else {
			// Invalid state where the caller wants to search but has defined no fields to search
			subQuery = subQuery.WHERE(RawBool("1=0"))
		}
	}
	// Paginate
	subQuery = subQuery.
		LIMIT(int64(page.Limit())).
		OFFSET(int64(page.Offset()))

	var res struct {
		Count     int `sql:"primary_key"`
		Envelopes []models.Envelope
	}

	err := SELECT(messagePJ, attachmentPJ, Raw("messages.count")).
		FROM(subQuery.AsTable("messages").LEFT_JOIN(Attachments, Attachments.MessageID.EQ(Messages.ID))).
		QueryContext(ctx, db, &res)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return models.DTOEnvelopeListResult{}, err
	}

	pageResult := pagination.NewPageResult(page, res.Count)
	return models.DTOEnvelopeListResult{
		PageResult: pageResult,
		Envelopes:  res.Envelopes,
	}, nil
}

func EnvelopeGet(ctx context.Context, db database.Querier, id int64) (models.Envelope, error) {
	var env models.Envelope
	err := SELECT(messagePJ, attachmentPJ).
		FROM(Messages.LEFT_JOIN(Attachments, Attachments.MessageID.EQ(Int64(id)))).
		WHERE(Messages.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &env)
	return env, err
}

func EnvelopeCount(ctx context.Context, db database.Querier) (int, error) {
	s := orm.CountSelect(Messages)
	return orm.CountQuery(ctx, db, s)
}

func EnvelopeDelete(ctx context.Context, db database.Querier, id int64) error {
	res, err := Messages.
		DELETE().
		WHERE(Messages.ID.EQ(Int64(id))).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	return one(res)
}

func EnvelopeDrop(ctx context.Context, db database.Querier) (int64, error) {
	res, err := Messages.
		DELETE().
		WHERE(RawBool("1=1")).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func EnvelopeTrim(ctx context.Context, db database.Querier, age time.Time, keep int) (int64, error) {
	where := Messages.CreatedAt.LT(RawTimestamp(muhTypeAffinity(models.NewTime(age))))
	if keep != 0 {
		where = where.AND(Messages.ID.NOT_IN(
			Messages.
				SELECT(Messages.ID).
				ORDER_BY(Messages.ID.DESC()).
				LIMIT(int64(keep)),
		))
	}

	res, err := Messages.
		DELETE().
		WHERE(where).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func MessageGet(ctx context.Context, db database.Querier, id int64) (models.Message, error) {
	var res models.Message
	err := Messages.
		SELECT(messagePJ).
		WHERE(Messages.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &res)
	return res, err
}

func MessageHTMLGet(ctx context.Context, db database.Querier, id int64) (string, error) {
	var res struct{ HTML string }
	err := Messages.
		SELECT(Messages.HTML.AS("html")).
		WHERE(Messages.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &res)
	return res.HTML, err
}

func AttachmentGet(ctx context.Context, db database.Querier, id int64) (models.Attachment, error) {
	var res models.Attachment
	err := Attachments.
		SELECT(attachmentPJ).
		WHERE(Attachments.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &res)
	return res, err
}

func AttachmentList(ctx context.Context, db database.Querier, page pagination.Page, req models.DTOAttachmentListRequest) (models.DTOAttachmentListResult, error) {
	query := Attachments.
		SELECT(attachmentPJ, COUNT(Raw("*")).OVER().AS("count")).
		LIMIT(int64(page.Limit())).
		OFFSET(int64(page.Offset()))
	// Order
	if req.Ascending {
		query = query.ORDER_BY(Attachments.ID.ASC())
	} else {
		query = query.ORDER_BY(Attachments.ID.DESC())
	}

	var res struct {
		Count       int `sql:"primary_key"`
		Attachments []models.Attachment
	}

	err := query.QueryContext(ctx, db, &res)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return models.DTOAttachmentListResult{}, err
	}

	pageResult := pagination.NewPageResult(page, res.Count)

	return models.DTOAttachmentListResult{
		PageResult:  pageResult,
		Attachments: res.Attachments,
	}, nil
}

func AttachmentListByMessage(ctx context.Context, db database.Querier, messageID int64) ([]models.Attachment, error) {
	var res []models.Attachment
	err := Attachments.
		SELECT(attachmentPJ).
		WHERE(Attachments.MessageID.EQ(Int64(messageID))).
		QueryContext(ctx, db, &res)
	return res, err
}

func AttachmentListOrphan(ctx context.Context, db database.Querier, limit int) ([]models.Attachment, error) {
	var atts []models.Attachment
	err := Attachments.
		SELECT(attachmentPJ).
		WHERE(Attachments.MessageID.IS_NULL()).
		LIMIT(int64(limit)).
		QueryContext(ctx, db, &atts)
	return atts, err
}

func AttachmentCount(ctx context.Context, db database.Querier) (int, error) {
	s := orm.CountSelect(Attachments)
	return orm.CountQuery(ctx, db, s)
}

// AttachmentRemoveOrphan should only be called after the associated file has been deleted from the FileStore.
func AttachmentRemoveOrphan(ctx context.Context, db database.Querier, id int64) error {
	_, err := Attachments.
		DELETE().
		WHERE(AND(
			Attachments.ID.EQ(Int64(id)),
			Attachments.MessageID.IS_NULL(),
		)).
		ExecContext(ctx, db)
	return err
}
