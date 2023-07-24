package inject

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func App(app core.App, fn func(c *fiber.Ctx, cc core.Context) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fn(c, app.Context(c.Context()))
	}
}

func AppID(app core.App, fn func(c *fiber.Ctx, cc core.Context, id int64) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return helpers.Error(c, err, http.StatusBadRequest)
		}

		return fn(c, app.Context(c.Context()), id)
	}
}

func AppStore(app core.App, store *session.Store, fn func(c *fiber.Ctx, cc core.Context, store *session.Store) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fn(c, app.Context(c.Context()), store)
	}
}
