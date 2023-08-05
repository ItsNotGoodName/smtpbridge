package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/dbgen/model"
	. "github.com/ItsNotGoodName/smtpbridge/internal/dbgen/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/uptrace/bun"
)

func InternalRuleUpsert(cc core.Context, r rules.Rule, updatedAt time.Time) error {
	_, err := Rules.INSERT(
		Rules.Internal,
		Rules.InternalID,
		Rules.Name,
		Rules.Expression,
		Rules.UpdatedAt,
		Rules.Enable,
	).MODEL(model.Rules{
		Internal:   true,
		InternalID: r.InternalID,
		Name:       r.Name,
		Expression: r.Expression,
		UpdatedAt:  updatedAt,
		Enable:     r.Enable,
	}).ON_CONFLICT(Rules.InternalID).DO_UPDATE(SET(
		Rules.Name.SET(Rules.EXCLUDED.Name),
		Rules.Expression.SET(Rules.EXCLUDED.Expression),
		Rules.UpdatedAt.SET(Rules.EXCLUDED.UpdatedAt),
	)).ExecContext(cc.Context(), cc.DB)

	return err
}

func InternalEndpointUpsert(cc core.Context, r endpoints.Endpoint, updatedAt time.Time) error {
	config, err := json.Marshal(r.Config)
	if err != nil {
		return err
	}

	_, err = Endpoints.INSERT(
		Endpoints.Internal,
		Endpoints.InternalID,
		Endpoints.Name,
		Endpoints.AttachmentDisable,
		Endpoints.TextDisable,
		Endpoints.BodyTemplate,
		Endpoints.Kind,
		Endpoints.Config,
		Endpoints.UpdatedAt,
	).MODEL(model.Endpoints{
		Internal:          true,
		InternalID:        r.InternalID,
		Name:              r.Name,
		AttachmentDisable: r.AttachmentDisable,
		TextDisable:       r.TextDisable,
		BodyTemplate:      r.BodyTemplate,
		Kind:              r.Kind,
		Config:            string(config),
		UpdatedAt:         updatedAt,
	}).ON_CONFLICT(Endpoints.InternalID).DO_UPDATE(SET(
		Endpoints.InternalID.SET(Endpoints.EXCLUDED.InternalID),
		Endpoints.Name.SET(Endpoints.EXCLUDED.Name),
		Endpoints.AttachmentDisable.SET(Endpoints.EXCLUDED.AttachmentDisable),
		Endpoints.TextDisable.SET(Endpoints.EXCLUDED.TextDisable),
		Endpoints.BodyTemplate.SET(Endpoints.EXCLUDED.BodyTemplate),
		Endpoints.Kind.SET(Endpoints.EXCLUDED.Kind),
		Endpoints.Config.SET(Endpoints.EXCLUDED.Config),
		Endpoints.UpdatedAt.SET(Endpoints.EXCLUDED.UpdatedAt),
	)).ExecContext(cc.Context(), cc.DB)

	return err
}

func InternalRuleEndpointsUpsert(cc core.Context, ruleInternalID string, endpointInternalIDs []string, updatedAt time.Time) error {
	if len(endpointInternalIDs) == 0 {
		return nil
	}

	for _, endpointInternalID := range endpointInternalIDs {
		rows, err := cc.DB.ExecContext(cc.Context(), `
			INSERT INTO rules_to_endpoints (
				internal,
				rule_id,
				endpoint_id,
				updated_at
			) SELECT "1" AS internal, rules.id AS rule_id, endpoints.id AS endpoint_id, ?
			FROM rules, endpoints
			WHERE rules.internal_id=? AND endpoints.internal_id IN (?) 
			ON CONFLICT (rule_id, endpoint_id) DO UPDATE SET updated_at=EXCLUDED.updated_at, internal=EXCLUDED.internal
		`, updatedAt.UTC(), ruleInternalID, endpointInternalID)
		if err != nil {
			return err
		}
		count, err := rows.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("endpoint not found '%s'", endpointInternalID)
		}
	}

	return nil
}

func InternalDeleteOlderThan(cc core.Context, date time.Time) error {
	return cc.DB.RunInTx(cc.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, table := range []string{"endpoints", "rules", "rules_to_endpoints"} {
			_, err := tx.NewDelete().Table(table).Where("updated_at < ?", date).Where("internal = true").Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
