package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

var ErrNotImplemented = fmt.Errorf("not implemented")

var ErrNoRows = qrm.ErrNoRows

func oneRowAffected(res sql.Result) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRows
	}

	return nil
}

func muhTypeAffinity(date models.Time) string {
	return "\"" + date.Time().Format(time.RFC3339) + "\""
}

func Size(ctx context.Context, db database.Querier) (int64, error) {
	var res struct{ Size int64 }
	err := RawStatement("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();").
		QueryContext(ctx, db, &res)
	return res.Size, err
}

func Vacuum(ctx context.Context, db database.Querier) error {
	_, err := RawStatement("VACUUM;").ExecContext(ctx, db)
	return err
}
