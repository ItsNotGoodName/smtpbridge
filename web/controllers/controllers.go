package controllers

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/ItsNotGoodName/smtpbridge/web"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Index(c *fiber.Ctx, cc core.Context) error {
	// Execute
	storage, err := procs.StorageGet(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	messages, err := procs.EnvelopeMessageList(cc, pagination.NewPage(1, 5), envelope.MessageFilter{})
	if err != nil {
		return helpers.Error(c, err)
	}

	policy := procs.RetentionPolicyGet(cc)

	// Response
	return c.Render("index", fiber.Map{
		"Messages":        messages.Messages,
		"Storage":         storage,
		"RetentionPolicy": policy,
	})
}

func IndexStorageTable(c *fiber.Ctx, cc core.Context) error {
	// Execute
	storage, err := procs.StorageGet(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return c.Render("index", fiber.Map{
		"Storage": storage,
	}, "storage-table")
}

func IndexRecentEnvelopesTable(c *fiber.Ctx, cc core.Context) error {
	// Execute
	messages, err := procs.EnvelopeMessageList(cc, pagination.NewPage(1, 5), envelope.MessageFilter{})
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return c.Render("index", fiber.Map{
		"Messages": messages.Messages,
	}, "recent-envelopes-table")
}

func Files(app core.App) fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:   http.FS(app.File.FS),
		MaxAge: web.CacheControl,
	})
}

func Send(c *fiber.Ctx, cc core.Context) error {
	// Request
	envelope_id, err := strconv.ParseInt(c.FormValue("envelope"), 10, 64)
	if err != nil {
		return helpers.Error(c, err, http.StatusBadRequest)
	}
	endpoint_id, err := strconv.ParseInt(c.FormValue("endpoint"), 10, 64)
	if err != nil {
		return helpers.Error(c, err, http.StatusBadRequest)
	}

	// Execute
	err = procs.EndpointSend(cc, envelope_id, endpoint_id)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return c.SendStatus(http.StatusNoContent)
}

func Vacuum(c *fiber.Ctx, cc core.Context) error {
	// Execute
	err := procs.DatabaseVacuum(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	c.Set("HX-Trigger", "vacuumed")
	return c.SendStatus(http.StatusNoContent)
}

func Trim(c *fiber.Ctx, cc core.Context) error {
	// Execute
	err := procs.TrimStart(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	c.Set("HX-Trigger", "trimmed")
	return c.SendStatus(http.StatusNoContent)
}

func SomethingWentWrong(c *fiber.Ctx) error {
	return c.Render("something-went-wrong", fiber.Map{"Error": ""})
}
