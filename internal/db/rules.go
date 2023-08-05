package db

import (
	"database/sql"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	. "github.com/ItsNotGoodName/smtpbridge/internal/dbgen/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
	. "github.com/go-jet/jet/v2/sqlite"
)

func RuleList(cc core.Context) ([]rules.Rule, error) {
	var rrules []rules.Rule
	err := cc.DB.NewSelect().Model(&rrules).Scan(cc.Context())
	return rrules, err
}

func RuleListEnable(cc core.Context) ([]rules.Rule, error) {
	var rrules []rules.Rule
	err := cc.DB.NewSelect().Model(&rrules).Where("enable = TRUE").Scan(cc.Context())
	return rrules, err
}

func RuleEndpointList(cc core.Context, ruleID int64) ([]rules.Endpoint, error) {
	var e []rules.Endpoint
	err := cc.DB.NewSelect().
		ColumnExpr("id, name, (rule_id IS NOT NULL) AS enable").
		TableExpr("endpoints").
		Join("LEFT JOIN rules_to_endpoints ON endpoints.id=rules_to_endpoints.endpoint_id AND rules_to_endpoints.rule_id=?", ruleID).
		Scan(cc.Context(), &e)

	return e, err
}

func RuleIsInternal(cc core.Context, ruleID int64) (bool, error) {
	rows, err := Rules.
		SELECT(Rules.Internal).
		WHERE(Rules.ID.EQ(Int64(ruleID))).
		Rows(cc.Context(), cc.DB)
	if err != nil {
		return false, err
	}
	var internal bool
	err = rows.Scan(&internal)
	return internal, err
}

func RuleUpdate(cc core.Context, ruleID int64, enable bool) (rules.Rule, error) {
	rule := rules.Rule{}
	res, err := cc.
		DB.
		NewUpdate().
		Model(&rule).
		Set("enable = ?", enable).
		Where("id = ?", ruleID).
		Returning("id, internal, internal_id, name, expression, enable").
		Exec(cc.Context(), &rule)
	if err != nil {
		return rule, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return rule, err
	}
	if rows == 0 {
		return rule, sql.ErrNoRows
	}

	return rule, err
}
