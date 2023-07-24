package procs

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
)

func DatabaseVacuum(cc core.Context) error {
	return db.Vacuum(cc)
}

func RetentionPolicyGet(cc core.Context) models.RetentionPolicy {
	return cc.Config.RetentionPolicy
}

func InternalSync(cc core.Context, eendpoints []endpoints.Endpoint, rrules []rules.Rule, ruleToEndpoints map[string][]string) error {
	updatedAt := time.Now()

	for _, se := range eendpoints {
		if err := db.InternalEndpointUpsert(cc, se, updatedAt); err != nil {
			return err
		}
	}

	for _, se := range rrules {
		if err := db.InternalRuleUpsert(cc, se, updatedAt); err != nil {
			return err
		}
	}

	for k, v := range ruleToEndpoints {
		if err := db.InternalRuleEndpointsUpsert(cc, k, v, updatedAt); err != nil {
			return err
		}
	}

	return db.InternalDeleteOlderThan(cc, updatedAt)
}
