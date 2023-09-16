package orm

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
)

func CountSelect(t sqlite.Table) SelectStatement {
	return t.SELECT(COUNT(Raw("*")).AS("count"))
}

func CountQuery(ctx context.Context, db database.Querier, s SelectStatement) (int, error) {
	var res struct{ Count int }
	err := s.QueryContext(ctx, db, &res)
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}
