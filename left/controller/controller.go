package controller

import (
	"errors"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/left/view"
)

type Controller struct {
	envelopeService envelope.Service
	endpointService endpoint.Service
}

func New(envelopeService envelope.Service, endpointService endpoint.Service) *Controller {
	return &Controller{
		envelopeService: envelopeService,
		endpointService: endpointService,
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

	view.Render(w, http.StatusOK, view.EnvelopePage, view.EnvelopeData{Envelope: env, Tab: r.URL.Query().Get("tab"), Endpoints: c.endpointService.ListEndpoint()})
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

func (c *Controller) EndpointList(w http.ResponseWriter, r *http.Request) {
	view.Render(w, http.StatusOK, view.EndpointsPage, view.EndpointsData{Endpoints: c.endpointService.ListEndpoint()})
}

func (c *Controller) EnvelopeSendPost(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endpointName, ok := r.PostForm["endpoint"]
	if !ok || len(endpointName) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	noAttachments := false
	noAttachmentsStr, ok := r.PostForm["no attachments"]
	if ok && len(endpointName) == 1 {
		noAttachments = (noAttachmentsStr[0] == "on")
	}
	noText := false
	noTextStr, ok := r.PostForm["no text"]
	if ok && len(endpointName) == 1 {
		noText = (noTextStr[0] == "on")
	}
	if noText && noAttachments {
		http.Error(w, "cannot enable both no attachments and no text", http.StatusBadRequest)
		return
	}

	// Get envelope
	id := parseID(r)
	ctx := r.Context()
	env, err := c.envelopeService.GetEnvelope(ctx, id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	// Get endpoint
	end, err := c.endpointService.GetEndpoint(endpointName[0])
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrEndpointNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	// Convert envelope message to text
	text := ""
	if !noText {
		text, err = end.Text(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Convert envelope attachments to endpoint attachments
	atts := []endpoint.Attachment{}
	if !noAttachments {
		for _, att := range env.Attachments {
			data, err := c.envelopeService.GetData(ctx, &att)
			if err != nil && err != core.ErrDataNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			atts = append(atts, endpoint.NewAttachment(&att, data))
		}
	}

	// Send to endpoint
	if err := end.Sender.Send(ctx, text, atts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
