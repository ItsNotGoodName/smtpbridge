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

func (c *Controller) EndpointTestPost(w http.ResponseWriter, r *http.Request) {
	endpointName := r.URL.Query().Get("endpoint")
	end, err := c.endpointService.GetEndpoint(endpointName)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrEndpointNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	if err := end.Sender.Send(r.Context(), "SMTPBridge test message", []endpoint.Attachment{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) EnvelopeSendPost(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endpointName := r.PostFormValue("endpoint")
	filter := r.PostFormValue("filter")
	noText := filter == "no text"
	noAttachments := filter == "no attachments"

	// Get envelope
	ctx := r.Context()
	env, err := c.envelopeService.GetEnvelope(ctx, parseID(r))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, core.ErrMessageNotFound) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
		return
	}

	// Get endpoint
	end, err := c.endpointService.GetEndpoint(endpointName)
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
