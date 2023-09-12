package repo

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	. "github.com/go-jet/jet/v2/sqlite"
)

var messagePJ ProjectionList = ProjectionList{
	Messages.ID.AS("message.id"),
	Messages.From.AS("message.from"),
	Messages.To.AS("message.to"),
	Messages.Subject.AS("message.subject"),
	Messages.Text.AS("message.text"),
	Messages.HTML.AS("message.html"),
	Messages.Date.AS("message.date"),
	Messages.CreatedAt.AS("message.created_at"),
}

var attachmentPJ ProjectionList = ProjectionList{
	Attachments.ID.AS("attachment.id"),
	Attachments.MessageID.AS("attachment.message_id"),
	Attachments.Name.AS("attachment.name"),
	Attachments.Mime.AS("attachment.mime"),
	Attachments.Extension.AS("attachment.extension"),
}

func EnvelopeCreate(ctx context.Context, db database.Querier, msg models.Message, atts []models.Attachment) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := Messages.
		INSERT(
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
	var res []models.Envelope

	subQuery := Messages.SELECT(Messages.ID)
	// This is what peak performance looks like
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

	query := SELECT(messagePJ, attachmentPJ).
		FROM(Messages.LEFT_JOIN(Attachments, Attachments.MessageID.EQ(Messages.ID))).
		WHERE(Messages.ID.IN(subQuery.LIMIT(int64(page.Limit())).OFFSET(int64(page.Offset()))))
	// s/subQuery/query/g
	if req.Ascending {
		if req.Order == models.DTOEnvelopeFieldSubject {
			query = query.ORDER_BY(Messages.Subject.ASC())
		} else if req.Order == models.DTOEnvelopeFieldFrom {
			query = query.ORDER_BY(Messages.From.ASC())
		} else {
			query = query.ORDER_BY(Messages.ID.ASC())
		}
	} else {
		if req.Order == models.DTOEnvelopeFieldSubject {
			query = query.ORDER_BY(Messages.Subject.DESC())
		} else if req.Order == models.DTOEnvelopeFieldFrom {
			query = query.ORDER_BY(Messages.From.DESC())
		} else {
			query = query.ORDER_BY(Messages.ID.DESC())
		}
	}

	err := query.QueryContext(ctx, db, &res)
	if err != nil {
		return models.DTOEnvelopeListResult{}, err
	}

	var resCount struct{ Count int }
	err = Messages.
		SELECT(COUNT(Raw("*")).AS("count")).
		QueryContext(ctx, db, &resCount)
	if err != nil {
		return models.DTOEnvelopeListResult{}, err
	}
	pageResult := pagination.NewPageResult(page, resCount.Count)

	return models.DTOEnvelopeListResult{
		PageResult: pageResult,
		Envelopes:  res,
	}, nil
}

func EnvelopeGet(ctx context.Context, db database.Querier, id int64) (models.Envelope, error) {
	var env models.Envelope
	err := SELECT(messagePJ, attachmentPJ).FROM(Messages.LEFT_JOIN(Attachments, Attachments.MessageID.EQ(Int64(id)))).WHERE(Messages.ID.EQ(Int64(id))).QueryContext(ctx, db, &env)
	return env, err
}

func EnvelopeCount(ctx context.Context, db database.Querier) (int, error) {
	var res struct{ Count int }
	err := Messages.
		SELECT(COUNT(Raw("*")).AS("count")).
		QueryContext(ctx, db, &res)
	return res.Count, err
}

func EnvelopeDelete(ctx context.Context, db database.Querier, id int64) error {
	res, err := Messages.
		DELETE().
		WHERE(Messages.ID.EQ(Int64(id))).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	return oneRowAffected(res)
}

func EnvelopeDrop(ctx context.Context, db database.Querier) error {
	_, err := Messages.
		DELETE().
		WHERE(RawBool("1=1")).
		ExecContext(ctx, db)
	return err
}

func EnvelopeTrim(ctx context.Context, db database.Querier, age time.Time, keep int) (int64, error) {
	q := Messages.CreatedAt.LT(RawTimestamp(muhTypeAffinity(models.NewTime(age))))
	if keep != 0 {
		q = q.AND(
			Messages.ID.NOT_IN(Messages.SELECT(Messages.ID).ORDER_BY(Messages.ID.DESC()).LIMIT(int64(keep))),
		)
	}

	res, err := Messages.
		DELETE().
		WHERE(q).
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
	withFilter := func(stmt SelectStatement) SelectStatement {
		return stmt
	}
	withOrder := func(stmt SelectStatement) SelectStatement {
		if req.Ascending {
			return stmt.ORDER_BY(Attachments.ID.ASC())
		}
		return stmt.ORDER_BY(Attachments.ID.DESC())
	}

	var res []models.Attachment

	query := withOrder(withFilter(Attachments.
		SELECT(attachmentPJ).
		LIMIT(int64(page.Limit())).
		OFFSET(int64(page.Offset()))))

	err := query.QueryContext(ctx, db, &res)
	if err != nil {
		return models.DTOAttachmentListResult{}, err
	}

	var resCount struct{ Count int }
	err = withFilter(Attachments.
		SELECT(COUNT(Raw("*")).AS("count"))).
		QueryContext(ctx, db, &resCount)
	if err != nil {
		return models.DTOAttachmentListResult{}, err
	}
	pageResult := pagination.NewPageResult(page, resCount.Count)

	return models.DTOAttachmentListResult{
		PageResult:  pageResult,
		Attachments: res,
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
	var res struct{ Count int }
	err := Attachments.
		SELECT(COUNT(Raw("*")).AS("count")).
		QueryContext(ctx, db, &res)
	return res.Count, err
}

// AttachmentRemove should only be called when it's MessageID is null and the associated file has been deleted from the FileStore.
func AttachmentRemove(ctx context.Context, db database.Querier, id int64) error {
	_, err := Attachments.
		DELETE().
		WHERE(Attachments.ID.EQ(Int64(id))).
		ExecContext(ctx, db)
	return err
}
