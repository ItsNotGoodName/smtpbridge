package repo

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	. "github.com/go-jet/jet/v2/sqlite"
)

var rulePJ ProjectionList = ProjectionList{
	Rules.ID.AS("rule.id"),
	Rules.Internal.AS("rule.internal"),
	Rules.InternalID.AS("rule.internal_id"),
	Rules.Name.AS("rule.name"),
	Rules.Expression.AS("rule.expression"),
	Rules.Enable.AS("rule.enable"),
}

var ruleCreatePJ ProjectionList = rulePJ.Except(Rules.ID)

func RuleCreate(ctx context.Context, db database.Querier, rule models.Rule, endpoints []int64) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	now := models.NewTime(time.Now())
	m := struct {
		models.Rule
		UpdatedAt models.Time
		CreatedAt models.Time
	}{
		Rule:      rule,
		UpdatedAt: now,
		CreatedAt: now,
	}
	res, err := Rules.INSERT(
		Rules.Name,
		Rules.Internal,
		Rules.InternalID,
		Rules.Expression,
		Rules.Enable,
		Rules.UpdatedAt,
		Rules.CreatedAt,
	).MODEL(m).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}
	ruleID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = ruleEndpointsSet(ctx, tx, ruleID, endpoints)
	if err != nil {
		return 0, err
	}

	return ruleID, tx.Commit()
}

func RuleUpdate(ctx context.Context, db database.Querier, rule models.Rule) error {
	m := struct {
		models.Rule
		UpdatedAt models.Time
	}{
		Rule:      rule,
		UpdatedAt: models.NewTime(time.Now()),
	}
	_, err := Rules.UPDATE(
		Rules.Name,
		Rules.Internal,
		Rules.InternalID,
		Rules.Expression,
		Rules.Enable,
		Rules.UpdatedAt,
	).MODEL(m).
		WHERE(Rules.ID.EQ(Int64(rule.ID))).
		ExecContext(ctx, db)
	return err
}

func RuleGet(ctx context.Context, db database.Querier, id int64) (models.Rule, error) {
	var res models.Rule
	err := Rules.SELECT(rulePJ).WHERE(Rules.ID.EQ(Int64(id))).QueryContext(ctx, db, &res)
	return res, err
}

func RuleList(ctx context.Context, db database.Querier) ([]models.Rule, error) {
	var res []models.Rule
	err := Rules.SELECT(rulePJ).WHERE(RawBool("1=1")).QueryContext(ctx, db, &res)
	return res, err
}

func RuleDelete(ctx context.Context, db database.Querier, id int64) error {
	res, err := Rules.DELETE().WHERE(Rules.ID.EQ(Int64(id))).ExecContext(ctx, db)
	if err != nil {
		return err
	}
	return oneRowAffected(res)
}

func RuleEndpointsGet(ctx context.Context, db database.Querier, id int64) (models.RuleEndpoints, error) {
	var res models.RuleEndpoints
	err := SELECT(rulePJ, endpointPJ).
		FROM(Rules.
			LEFT_JOIN(RulesToEndpoints, RulesToEndpoints.RuleID.EQ(Rules.ID)).
			LEFT_JOIN(Endpoints, Endpoints.ID.EQ(RulesToEndpoints.EndpointID)),
		).
		WHERE(Rules.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &res)
	return res, err
}

func RuleEndpointsList(ctx context.Context, db database.Querier) ([]models.RuleEndpoints, error) {
	var res []models.RuleEndpoints
	err := SELECT(rulePJ, endpointPJ).
		FROM(Rules.
			LEFT_JOIN(RulesToEndpoints, RulesToEndpoints.RuleID.EQ(Rules.ID)).
			LEFT_JOIN(Endpoints, Endpoints.ID.EQ(RulesToEndpoints.EndpointID)),
		).
		QueryContext(ctx, db, &res)
	return res, err
}

func RuleEndpointsSet(ctx context.Context, db database.Querier, ruleID int64, endpointIDs []int64) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = RulesToEndpoints.
		DELETE().
		WHERE(RulesToEndpoints.RuleID.EQ(Int64(ruleID)).AND(RulesToEndpoints.Internal.EQ(Bool(false)))).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	err = ruleEndpointsSet(ctx, tx, ruleID, endpointIDs)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ruleEndpointsSet(ctx context.Context, db database.QuerierTx, ruleID int64, endpointIDs []int64) error {
	if len(endpointIDs) == 0 {
		return nil
	}

	stmt := RulesToEndpoints.
		INSERT(
			RulesToEndpoints.Internal,
			RulesToEndpoints.RuleID,
			RulesToEndpoints.EndpointID,
		)
	for _, v := range endpointIDs {
		stmt = stmt.VALUES(false, ruleID, v)
	}

	_, err := stmt.ExecContext(ctx, db)
	return err
}
