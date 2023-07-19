package db

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
)

func RuleList(cc *core.Context) ([]rules.Rule, error) {
	var rrules []rules.Rule
	err := cc.DB.NewSelect().Model(&rrules).Scan(cc.Context())
	return rrules, err
}

func RuleEndpointList(cc *core.Context, ruleID int64) ([]rules.Endpoint, error) {
	var e []rules.Endpoint
	err := cc.DB.NewSelect().
		ColumnExpr("id, name, (rule_id IS NOT NULL) AS enable").
		TableExpr("endpoints").
		Join("LEFT JOIN rules_to_endpoints ON endpoints.id=rules_to_endpoints.endpoint_id AND rules_to_endpoints.rule_id=?", ruleID).
		Scan(cc.Context(), &e)

	return e, err
}
