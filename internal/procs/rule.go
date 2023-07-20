package procs

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
)

func RuleAggregateList(cc *core.Context) ([]rules.Aggregate, error) {
	rrules, err := db.RuleList(cc)
	if err != nil {
		return nil, err
	}

	aggregateRules := make([]rules.Aggregate, len(rrules))
	for i := range rrules {
		aggregateEndpoints, err := db.RuleEndpointList(cc, rrules[i].ID)
		if err != nil {
			return nil, err
		}

		aggregateRules[i] = rules.Aggregate{
			Rule:      rrules[i],
			Endpoints: aggregateEndpoints,
		}
	}

	return aggregateRules, nil
}

func RuleUpdateEnable(cc *core.Context, ruleID int64, enable bool) (rules.Rule, error) {
	return db.RuleUpdate(cc, ruleID, enable)
}
