package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core"
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

func (c *Controller) IndexGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageQ, _ := strconv.Atoi(q.Get("page"))
	limitQ, _ := strconv.Atoi(q.Get("limit"))
	ascendingQ, _ := strconv.ParseBool(q.Get("ascending"))

	page := paginate.NewPage(pageQ, limitQ, ascendingQ)
	envs, err := c.envelopeService.ListEnvelope(r.Context(), &page)
	if err != nil {
		view.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	view.Render(w, http.StatusOK, view.IndexPage, view.IndexData{Envelopes: envs, Page: page})
}

func (c *Controller) EnvelopeGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	tab := r.URL.Query().Get("tab")

	env, err := c.envelopeService.GetEnvelope(r.Context(), id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		view.RenderError(w, code, err)
		return
	}

	view.Render(w, http.StatusOK, view.EnvelopePage, view.EnvelopeData{Envelope: env, Tab: tab})
}

func (c *Controller) EnvelopeDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	err := c.envelopeService.DeleteEnvelope(r.Context(), id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
