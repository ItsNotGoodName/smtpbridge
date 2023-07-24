package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
)

func checkbox(c *fiber.Ctx, key string) bool {
	isSet := c.Query("-"+key) != ""
	if isSet {
		return c.Query(key) != ""
	}

	return true
}
func redirect(c *fiber.Ctx, path string) error {
	if helpers.IsHTMXRequest(c) {
		c.Set("HX-Location", path)
		return nil
	} else {
		return c.Redirect(path)
	}
}
