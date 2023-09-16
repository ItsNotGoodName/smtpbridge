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

func InternalSync(
	ctx context.Context,
	db database.Querier,
	endpoints []models.Endpoint,
	rules []models.Rule,
	ruleToEndpoints map[string][]string,
) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updatedAt := models.NewTime(time.Now())

	for _, end := range endpoints {
		if err := internalEndpointUpsert(ctx, tx, end, updatedAt); err != nil {
			return err
		}
	}

	for _, rule := range rules {
		if err := internalRuleUpsert(ctx, tx, rule, updatedAt); err != nil {
			return err
		}
	}

	if err := internalDeleteOlderThan(ctx, tx, updatedAt); err != nil {
		return err
	}

	for k, v := range ruleToEndpoints {
		if err := internalRuleEndpointsInsert(ctx, tx, k, v); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func internalEndpointUpsert(ctx context.Context, tx database.QuerierTx, r models.Endpoint, updatedAt models.Time) error {
	m := struct {
		models.Endpoint
		UpdatedAt models.Time
		CreatedAt models.Time
	}{
		Endpoint:  r,
		UpdatedAt: updatedAt,
		CreatedAt: models.NewTime(time.Now()),
	}

	res, err := Endpoints.
		UPDATE(
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
		).
		MODEL(m).
		WHERE(Endpoints.InternalID.EQ(String(r.InternalID.String))).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
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
			ExecContext(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func internalRuleUpsert(ctx context.Context, tx database.QuerierTx, r models.Rule, updatedAt models.Time) error {
	m := struct {
		models.Rule
		UpdatedAt models.Time
		CreatedAt models.Time
	}{
		Rule:      r,
		UpdatedAt: updatedAt,
		CreatedAt: models.NewTime(time.Now()),
	}

	res, err := Rules.
		UPDATE(
			Rules.Internal,
			Rules.InternalID,
			Rules.Name,
			Rules.Expression,
			Rules.UpdatedAt,
		).
		MODEL(m).
		WHERE(Rules.InternalID.EQ(String(r.InternalID.String))).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
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
			ExecContext(ctx, tx)
		if err != nil {
			return err
		}
	}

	return err
}

func internalRuleEndpointsInsert(ctx context.Context, tx database.QuerierTx, ruleInternalID string, endpointInternalIDs []string) error {
	_, err := RulesToEndpoints.
		DELETE().
		WHERE(AND(
			RulesToEndpoints.RuleID.IN(Rules.SELECT(Rules.ID).WHERE(Rules.InternalID.EQ(String(ruleInternalID)))),
			RulesToEndpoints.Internal.EQ(Bool(true)),
		)).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	for _, endpointInternalID := range endpointInternalIDs {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO rules_to_endpoints (
				internal,
				rule_id,
				endpoint_id
			) SELECT "1" AS internal, rules.id AS rule_id, endpoints.id AS endpoint_id
			FROM rules, endpoints
			WHERE rules.internal_id=? AND endpoints.internal_id IN (?)
			LIMIT 1
			ON CONFLICT (rule_id, endpoint_id) DO UPDATE SET internal=EXCLUDED.internal
		`, ruleInternalID, endpointInternalID)
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

func internalDeleteOlderThan(ctx context.Context, tx database.QuerierTx, date models.Time) error {
	_, err := Rules.
		DELETE().
		WHERE(AND(
			Rules.Internal.IS_TRUE(),
			Rules.UpdatedAt.LT(RawTimestamp(muhTypeAffinity(date))),
		)).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	_, err = Endpoints.
		DELETE().
		WHERE(AND(
			Endpoints.Internal.IS_TRUE(),
			Endpoints.UpdatedAt.LT(RawTimestamp(muhTypeAffinity(date))),
		)).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
