package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// checkbox allows 3 states for a input checkbox (checked, unchecked, and unset).
func checkbox(c *fiber.Ctx, key string) bool {
	isSet := c.Query("-"+key) != ""
	if isSet {
		return c.Query(key) != ""
	}

	return true
}
