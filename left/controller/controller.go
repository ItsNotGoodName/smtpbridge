package controller

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
	"github.com/ItsNotGoodName/smtpbridge/left/view"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	envelopeService envelope.Service
}

func New(envelopeService envelope.Service) *Controller {
	return &Controller{
		envelopeService: envelopeService,
	}
}

func (c *Controller) IndexGet(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageQ, _ := strconv.Atoi(q.Get("page"))
	limitQ, _ := strconv.Atoi(q.Get("limit"))
	ascendingQ, _ := strconv.ParseBool(q.Get("ascending"))

	page := paginate.NewPage(pageQ, limitQ, ascendingQ)
	envs, err := c.envelopeService.ListEnvelope(r.Context(), &page)
	if err != nil {
		view.RenderError(rw, http.StatusInternalServerError, err)
		return
	}

	view.Render(rw, http.StatusOK, view.IndexPage, view.IndexData{Envelopes: envs, Page: page})
}

func (c *Controller) EnvelopeGet(rw http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	tab := r.URL.Query().Get("tab")

	env, err := c.envelopeService.GetEnvelope(r.Context(), id)
	if err != nil {
		view.RenderError(rw, http.StatusInternalServerError, err)
		return
	}

	view.Render(rw, http.StatusOK, view.EnvelopePage, view.EnvelopeData{Envelope: env, Tab: tab})
}
