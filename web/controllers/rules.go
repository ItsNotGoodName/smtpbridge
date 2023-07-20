package controllers

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/gofiber/fiber/v2"
)

func Rules(c *fiber.Ctx, cc *core.Context) error {
	// Execute
	aggregateRules, err := procs.RuleAggregateList(cc)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Response
	return c.Render("rules", fiber.Map{
		"AggregateRules": aggregateRules,
	})
}

func RulesEnable(c *fiber.Ctx, cc *core.Context, id int64) error {
	rule, err := procs.RuleUpdateEnable(cc, id, c.FormValue("enable") == "on")
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Response
	return c.Render("p/rule-enable-form", rule)
}
