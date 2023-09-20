package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo/orm"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/samber/lo"
)

var tracePJ ProjectionList = ProjectionList{
	Traces.ID.AS("trace.id"),
	Traces.RequestID.AS("trace.request_id"),
	Traces.Source.AS("trace.source"),
	Traces.Seq.AS("trace.seq"),
	Traces.Action.AS("trace.action"),
	Traces.Data.AS("trace.data"),
	Traces.Level.AS("trace.level"),
	Traces.CreatedAt.AS("trace.created_at"),
}

type TraceStore struct {
	db database.Querier
}

func NewTraceStore(db database.Querier) TraceStore {
	return TraceStore{
		db: db,
	}
}

func (r TraceStore) Save(ctx context.Context, trace models.Trace) error {
	_, err := TraceCreate(ctx, r.db, trace)
	return err
}

func TraceCreate(ctx context.Context, db database.Querier, r models.Trace) (int64, error) {
	res, err := Traces.
		INSERT(
			Traces.RequestID,
			Traces.Source,
			Traces.Seq,
			Traces.Action,
			Traces.Data,
			Traces.Level,
			Traces.CreatedAt,
		).
		MODEL(r).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func TraceList(ctx context.Context, db database.Querier, page pagination.Page, req models.DTOTraceListRequest) (models.DTOTraceListResult, error) {
	subQuery := Traces.
		SELECT(Traces.RequestID.AS(Traces.RequestID.Name())).
		DISTINCT()
	// Order
	subQuery = traceListOrder(subQuery, req)
	// Pagination
	subQuery = subQuery.
		LIMIT(int64(page.Limit())).
		OFFSET(int64(page.Offset()))

	var res []models.Trace
	query := SELECT(tracePJ).
		FROM(subQuery.AsTable("t").
			LEFT_JOIN(Traces, RawString(fmt.Sprintf("t.%s", Traces.RequestID.Name())).EQ(Traces.RequestID)))
	// Order
	query = traceListOrder(query, req)
	err := query.QueryContext(ctx, db, &res)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return models.DTOTraceListResult{}, err
	}

	traces := lo.PartitionBy(res, func(t models.Trace) string { return t.RequestID })

	count, err := orm.CountQuery(ctx, db, Traces.
		SELECT(COUNT(Raw(fmt.Sprintf("DISTINCT %s", Traces.RequestID.Name()))).AS("count")))
	if err != nil {
		return models.DTOTraceListResult{}, err
	}

	pageResult := pagination.NewPageResult(page, count)
	return models.DTOTraceListResult{
		PageResult: pageResult,
		Traces:     traces,
	}, nil
}

func traceListOrder(s SelectStatement, req models.DTOTraceListRequest) SelectStatement {
	if req.Ascending {
		return s.ORDER_BY(Traces.ID.ASC())
	}
	return s.ORDER_BY(Traces.ID.DESC())
}

func TraceDrop(ctx context.Context, db database.Querier) (int64, error) {
	res, err := Traces.
		DELETE().
		WHERE(RawBool("1=1")).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func TraceTrim(ctx context.Context, db database.Querier, age time.Time) (int64, error) {
	res, err := Traces.
		DELETE().
		WHERE(Traces.CreatedAt.LT(RawTimestamp(muhTypeAffinity(models.NewTime(age))))).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
