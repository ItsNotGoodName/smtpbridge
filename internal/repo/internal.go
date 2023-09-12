package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	. "github.com/go-jet/jet/v2/sqlite"
)

// TODO: instead of doing INSERT ... ON CONFLICT DO UPDATE SET ..., run in a transaction

func InternalSync(
	ctx context.Context,
	db database.Querier,
	endpoints []models.Endpoint,
	rules []models.Rule,
	ruleToEndpoints map[string][]string,
) error {
	updatedAt := models.NewTime(time.Now())

	for _, end := range endpoints {
		if err := internalEndpointUpsert(ctx, db, end, updatedAt); err != nil {
			return err
		}
	}

	for _, rule := range rules {
		if err := internalRuleUpsert(ctx, db, rule, updatedAt); err != nil {
			return err
		}
	}

	for k, v := range ruleToEndpoints {
		if err := internalRuleEndpointsUpsert(ctx, db, k, v, updatedAt); err != nil {
			return err
		}
	}

	return internalDeleteOlderThan(ctx, db, updatedAt)
}

func internalEndpointUpsert(ctx context.Context, db database.Querier, r models.Endpoint, updatedAt models.Time) error {
	m := struct {
		models.Endpoint
		UpdatedAt models.Time
		CreatedAt models.Time
	}{
		Endpoint:  r,
		UpdatedAt: updatedAt,
		CreatedAt: models.NewTime(time.Now()),
	}
	_, err := Endpoints.
		INSERT(
			Endpoints.Internal,
			Endpoints.InternalID,
			Endpoints.Name,
			Endpoints.AttachmentDisable,
			Endpoints.TextDisable,
			Endpoints.TitleTemplate,
			Endpoints.BodyTemplate,
			Endpoints.Kind,
			Endpoints.Config,
			Endpoints.UpdatedAt,
			Endpoints.CreatedAt,
		).
		MODEL(m).
		ON_CONFLICT(Endpoints.InternalID).
		DO_UPDATE(SET(
			Endpoints.InternalID.SET(Endpoints.EXCLUDED.InternalID),
			Endpoints.Name.SET(Endpoints.EXCLUDED.Name),
			Endpoints.AttachmentDisable.SET(Endpoints.EXCLUDED.AttachmentDisable),
			Endpoints.TextDisable.SET(Endpoints.EXCLUDED.TextDisable),
			Endpoints.TitleTemplate.SET(Endpoints.EXCLUDED.TitleTemplate),
			Endpoints.BodyTemplate.SET(Endpoints.EXCLUDED.BodyTemplate),
			Endpoints.Kind.SET(Endpoints.EXCLUDED.Kind),
			Endpoints.Config.SET(Endpoints.EXCLUDED.Config),
			Endpoints.UpdatedAt.SET(Endpoints.EXCLUDED.UpdatedAt),
		)).
		ExecContext(ctx, db)
	return err
}

func internalRuleUpsert(ctx context.Context, db database.Querier, r models.Rule, updatedAt models.Time) error {
	m := struct {
		models.Rule
		UpdatedAt models.Time
		CreatedAt models.Time
	}{
		Rule:      r,
		UpdatedAt: updatedAt,
		CreatedAt: models.NewTime(time.Now()),
	}
	_, err := Rules.
		INSERT(
			Rules.Internal,
			Rules.InternalID,
			Rules.Name,
			Rules.Expression,
			Rules.Enable,
			Rules.UpdatedAt,
			Rules.CreatedAt,
		).
		MODEL(m).
		ON_CONFLICT(Rules.InternalID).
		DO_UPDATE(SET(
			Rules.Name.SET(Rules.EXCLUDED.Name),
			Rules.Expression.SET(Rules.EXCLUDED.Expression),
			Rules.UpdatedAt.SET(Rules.EXCLUDED.UpdatedAt),
		)).
		ExecContext(ctx, db)

	return err
}

func internalRuleEndpointsUpsert(ctx context.Context, db database.Querier, ruleInternalID string, endpointInternalIDs []string, updatedAt models.Time) error {
	if len(endpointInternalIDs) == 0 {
		return nil
	}

	for _, endpointInternalID := range endpointInternalIDs {
		// TODO: refactor this
		res, err := db.ExecContext(ctx, `
			INSERT INTO rules_to_endpoints (
				internal,
				rule_id,
				endpoint_id,
				updated_at,
				created_at
			) SELECT "1" AS internal, rules.id AS rule_id, endpoints.id AS endpoint_id, ?, ?
			FROM rules, endpoints
			WHERE rules.internal_id=? AND endpoints.internal_id IN (?) 
			ON CONFLICT (rule_id, endpoint_id) DO UPDATE SET updated_at=EXCLUDED.updated_at, internal=EXCLUDED.internal
		`, updatedAt, updatedAt, ruleInternalID, endpointInternalID)
		if err != nil {
			return err
		}
		count, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("endpoint not found '%s'", endpointInternalID)
		}
	}

	return nil
}

func internalDeleteOlderThan(ctx context.Context, db database.Querier, date models.Time) error {
	_, err := RulesToEndpoints.
		DELETE().
		WHERE(AND(
			RulesToEndpoints.Internal.IS_TRUE(),
			RulesToEndpoints.UpdatedAt.LT(RawTimestamp(muhTypeAffinity(date))),
		)).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	_, err = Rules.
		DELETE().
		WHERE(AND(
			Rules.Internal.IS_TRUE(),
			Rules.UpdatedAt.LT(RawTimestamp(muhTypeAffinity(date))),
		)).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	_, err = Endpoints.
		DELETE().
		WHERE(AND(
			Endpoints.Internal.IS_TRUE(),
			Endpoints.UpdatedAt.LT(RawTimestamp(muhTypeAffinity(date))),
		)).
		ExecContext(ctx, db)
	return err
}
