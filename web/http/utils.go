package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

func csrfExtractor() func(c *fiber.Ctx) (string, error) {
	form := csrf.CsrfFromForm("csrf")
	header := csrf.CsrfFromHeader(csrf.HeaderName)

	return func(c *fiber.Ctx) (string, error) {
		if token, err := form(c); err == nil {
			return token, nil
		}

		return header(c)
	}
}
