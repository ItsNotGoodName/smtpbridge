package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
)

func PRecentEnvelopesTable(c *fiber.Ctx, cc *core.Context) error {
	// Execute
	messages, err := procs.EnvelopeMessageList(cc, pagination.NewPage(1, 5), envelope.MessageFilter{})
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return partial(c, "p/recent-envelopes-table", fiber.Map{
		"Messages": messages.Messages,
	})
}

func PStorageTable(c *fiber.Ctx, cc *core.Context) error {
	// Execute
	storage, err := procs.StorageGet(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return partial(c, "p/storage-table", fiber.Map{
		"Storage": storage,
	})
}
