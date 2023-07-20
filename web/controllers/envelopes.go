package controllers

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
)

func Envelopes(c *fiber.Ctx, cc *core.Context) error {
	// Request
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return helpers.Error(c, err, http.StatusBadRequest)
	}

	perPage, err := strconv.Atoi(c.Query("perPage", "1"))
	if err != nil {
		return helpers.Error(c, err, http.StatusBadRequest)
	}

	// Execute
	filter := envelope.MessageFilter{
		Ascending:     c.Query("ascending") != "",
		Search:        c.Query("search"),
		SearchSubject: checkbox(c, "search-subject"),
		SearchText:    checkbox(c, "search-body"),
	}
	res, err := procs.EnvelopeMessageList(cc, pagination.NewPage(page, perPage), filter)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	queries := c.Queries()
	if res.PageResult.Page > res.PageResult.TotalPages {
		return c.Redirect("/envelopes?" + helpers.Query(queries, "page", res.PageResult.TotalPages))
	}

	return c.Render("envelopes", fiber.Map{
		"Queries":            queries,
		"Messages":           res.Messages,
		"MessagesPageResult": res.PageResult,
		"MessagesFilter":     filter,
	})
}

func EnvelopesView(c *fiber.Ctx, cc *core.Context, id int64) error {
	// Execute
	env, err := procs.EnvelopeGet(cc, id)
	if err != nil {
		return helpers.Error(c, err)
	}
	ends, err := procs.EndpointList(cc)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return c.Render("envelopes-view", fiber.Map{
		"Envelope":  env,
		"Endpoints": ends,
	})
}

func EnvelopesViewHTML(c *fiber.Ctx, cc *core.Context, id int64) error {
	// Execute
	html, err := procs.EnvelopeMessageHTMLGet(cc, id)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	c.Context().SetContentType("text/html; charset=utf-8")
	return c.SendString(html)
}

func EnvelopesNewGet(c *fiber.Ctx) error {
	// Render
	return c.Render("envelopes-new", fiber.Map{})
}

func EnvelopesNewPost(c *fiber.Ctx, cc *core.Context) error {
	// Request
	form, err := c.MultipartForm()
	if err != nil {
		return helpers.Error(c, err, http.StatusBadRequest)
	}
	var datts []envelope.DataAttachment
	for _, fh := range form.File["attachments"] {
		a, err := fh.Open()
		if err != nil {
			return helpers.Error(c, err, http.StatusBadRequest)
		}

		data, err := io.ReadAll(a)
		if err != nil {
			return helpers.Error(c, err, http.StatusBadRequest)
		}

		datts = append(datts, envelope.NewDataAttachment(fh.Filename, data))
	}

	msg := envelope.NewMessage(
		c.FormValue("from"),
		strings.Split(c.FormValue("to"), ","),
		c.FormValue("subject"),
		c.FormValue("body"),
		"",
		time.Now(),
	)

	// Execute
	_, err = procs.EnvelopeCreate(cc, msg, datts)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	return c.Redirect("/envelopes")
}

func EnvelopesDelete(c *fiber.Ctx, cc *core.Context, id int64) error {
	// Execute
	err := procs.EnvelopeDelete(cc, id)
	if err != nil {
		return helpers.Error(c, err)
	}

	// Response
	c.Set("HX-Redirect", "/envelopes")
	return c.SendStatus(http.StatusNoContent)
}
