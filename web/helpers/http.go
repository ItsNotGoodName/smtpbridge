package helpers

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, err error, code ...int) error {
	if len(code) == 1 {
		return c.Status(code[0]).Render("500", fiber.Map{"Error": err})
	}

	if err == sql.ErrNoRows {
		return NotFound(c)
	}

	return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{"Error": err})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Render("404", nil)
}
