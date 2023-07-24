package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
)

func Rules(c *fiber.Ctx, cc core.Context) error {
	// Execute
	aggregateRules, err := procs.RuleAggregateList(cc)
	if err != nil {
		return h.Error(c, err)
	}

	// Response
	return h.Render(c, "rules", fiber.Map{
		"AggregateRules": aggregateRules,
	})
}

func RuleEnable(c *fiber.Ctx, cc core.Context, id int64) error {
	// Request
	enable := c.FormValue("enable") == "on"

	// Execute
	rule, err := procs.RuleUpdateEnable(cc, id, enable)
	if err != nil {
		return h.Error(c, err)
	}

	// Response
	return h.Render(c, "rules", fiber.Map{
		"Rule": rule,
	}, "rule-enable-form")
}
