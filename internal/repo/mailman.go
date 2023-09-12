package repo

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	. "github.com/go-jet/jet/v2/sqlite"
)

func MailmanEnqueue(ctx context.Context, db database.Querier, envelopeID int64) error {
	_, err := MailmanQueue.
		INSERT(
			MailmanQueue.MessageID,
			MailmanQueue.CreatedAt,
		).
		VALUES(
			envelopeID,
			models.NewTime(time.Now()),
		).
		ExecContext(ctx, db)
	return err
}

func MailmanDequeue(ctx context.Context, db database.Querier) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var res struct {
		MessageID int64
	}
	err = MailmanQueue.
		SELECT(MailmanQueue.MessageID.AS("message_id")).
		WHERE(RawBool("1=1")).
		ORDER_BY(MailmanQueue.CreatedAt.ASC()).
		LIMIT(1).
		QueryContext(ctx, tx, &res)
	if err != nil {
		return 0, err
	}

	_, err = MailmanQueue.
		DELETE().
		WHERE(MailmanQueue.MessageID.EQ(Int64(res.MessageID))).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return res.MessageID, nil
}
