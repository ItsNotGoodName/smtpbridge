package controllers

import "github.com/gofiber/fiber/v2"

func checkbox(c *fiber.Ctx, key string) bool {
	isSet := c.Query("-"+key) != ""
	if isSet {
		return c.Query(key) != ""
	}

	return true

}
