package repo

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/go-jet/jet/v2/sqlite"
)

func Size(ctx context.Context, db database.Querier) (int64, error) {
	var res struct{ Size int64 }
	err := RawStatement("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();").
		QueryContext(ctx, db, &res)
	return res.Size, err
}

func Vacuum(ctx context.Context, db database.Querier) error {
	_, err := RawStatement("VACUUM;").
		ExecContext(ctx, db)
	return err
}
