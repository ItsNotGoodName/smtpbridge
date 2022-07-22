package controller

import (
	"errors"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/left/view"
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
	page := parsePage(r)
	envs, err := c.envelopeService.ListEnvelope(r.Context(), &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.Render(w, http.StatusOK, view.IndexPage, view.IndexData{Envelopes: envs, Page: page})
}

func (c *Controller) EnvelopeGet(w http.ResponseWriter, r *http.Request) {
	env, err := c.envelopeService.GetEnvelope(r.Context(), parseID(r))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	view.Render(w, http.StatusOK, view.EnvelopePage, view.EnvelopeData{Envelope: env, Tab: r.URL.Query().Get("tab")})
}

func (c *Controller) EnvelopeHTMLGet(w http.ResponseWriter, r *http.Request) {
	env, err := c.envelopeService.GetEnvelope(r.Context(), parseID(r))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	w.Write([]byte(env.Message.HTML))
}

func (c *Controller) EnvelopeDelete(w http.ResponseWriter, r *http.Request) {
	err := c.envelopeService.DeleteEnvelope(r.Context(), parseID(r))
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

func (c *Controller) AttachmentList(w http.ResponseWriter, r *http.Request) {
	page := parsePage(r)
	atts, err := c.envelopeService.ListAttachment(r.Context(), &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.Render(w, http.StatusOK, view.AttachmentsPage, view.AttachmentsData{Attachments: atts, Page: page})
}
