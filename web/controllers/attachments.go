package controllers

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
)

func Attachments(c *fiber.Ctx, cc core.Context) error {
	// Request
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return h.Error(c, err, http.StatusBadRequest)
	}

	perPage, err := strconv.Atoi(c.Query("perPage", "1"))
	if err != nil {
		return h.Error(c, err, http.StatusBadRequest)
	}

	// Execute
	filter := envelope.AttachmentFilter{
		Ascending: c.Query("ascending") != "",
	}
	res, err := procs.EnvelopeAttachmentList(cc, pagination.NewPage(page, perPage), filter)
	if err != nil {
		return h.Error(c, err)
	}

	// Response
	queries := c.Queries()
	if res.PageResult.Page > res.PageResult.TotalPages {
		return c.Redirect("/attachments?" + h.Query(queries, "page", res.PageResult.TotalPages))
	}

	return h.Render(c, "attachments", fiber.Map{
		"Queries":               queries,
		"Attachments":           res.Attachments,
		"AttachmentsPageResult": res.PageResult,
		"AttachmentsFilter":     filter,
	})
}
