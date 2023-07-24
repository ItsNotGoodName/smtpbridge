package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db/queries"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
	"github.com/uptrace/bun"
)

func InternalRuleUpsert(cc core.Context, r rules.Rule, updatedAt time.Time) error {
	return queries.New(cc.DB.DB).UpsertInternalRule(cc.Context(), queries.UpsertInternalRuleParams{
		InternalID: r.InternalID,
		Name:       r.Name,
		Expression: r.Expression,
		UpdatedAt:  updatedAt.UTC(),
		Enable:     r.Enable,
	})
}

func InternalEndpointUpsert(cc core.Context, end endpoints.Endpoint, updatedAt time.Time) error {
	config, err := json.Marshal(end.Config)
	if err != nil {
		return err
	}

	return queries.New(cc.DB.DB).UpsertInternalEndpoint(cc.Context(), queries.UpsertInternalEndpointParams{
		InternalID:        end.InternalID,
		Name:              end.Name,
		AttachmentDisable: end.AttachmentDisable,
		TextDisable:       end.TextDisable,
		BodyTemplate:      end.BodyTemplate,
		Kind:              end.Kind,
		Config:            config,
		UpdatedAt:         updatedAt.UTC(),
	})
}

func InternalRuleEndpointsUpsert(cc core.Context, ruleInternalID string, endpointInternalIDs []string, updatedAt time.Time) error {
	if len(endpointInternalIDs) == 0 {
		return nil
	}

	for _, v := range endpointInternalIDs {
		rows, err := queries.New(cc.DB.DB).UpsertInternalRuleToEndpoint(cc.Context(), queries.UpsertInternalRuleToEndpointParams{
			UpdatedAt:    updatedAt.UTC(),
			InternalID:   ruleInternalID,
			InternalID_2: v,
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("endpoint not found '%s'", v)
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
