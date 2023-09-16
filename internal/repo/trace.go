package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
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
		SELECT(
			Traces.RequestID.AS("request_id"),
			COUNT(Raw("*")).OVER().AS("count"),
		).
		DISTINCT()
	// Order
	if req.Ascending {
		subQuery = subQuery.ORDER_BY(Traces.ID.ASC())
	} else {
		subQuery = subQuery.ORDER_BY(Traces.ID.DESC())
	}
	// Pagination
	subQuery = subQuery.
		LIMIT(int64(page.Limit())).
		OFFSET(int64(page.Offset()))

	var res struct {
		Count int `sql:"primary_key"`
		Trace []models.Trace
	}
	err := SELECT(tracePJ, Raw("t.count").AS("count")).
		FROM(subQuery.AsTable("t").
			LEFT_JOIN(Traces, RawString("t.request_id").EQ(Traces.RequestID))).
		QueryContext(ctx, db, &res)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return models.DTOTraceListResult{}, err
	}

	traces := lo.PartitionBy(res.Trace, func(t models.Trace) string { return t.RequestID })

	pageResult := pagination.NewPageResult(page, res.Count)
	return models.DTOTraceListResult{
		PageResult: pageResult,
		Traces:     traces,
	}, nil
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
