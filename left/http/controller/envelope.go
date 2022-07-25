package controller

import (
	"errors"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/left/http/view"
)

func IndexGet(envelopeService envelope.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := parsePage(r)
		envs, err := envelopeService.ListEnvelope(r.Context(), &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		view.Render(w, http.StatusOK, view.IndexPage, view.IndexData{Envelopes: envs, Page: page})
	}
}

func EnvelopeGet(envelopeService envelope.Service, endpointService endpoint.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env, err := envelopeService.GetEnvelope(r.Context(), parseID(r))
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, core.ErrMessageNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, err.Error(), code)
			return
		}

		view.Render(w, http.StatusOK, view.EnvelopePage, view.EnvelopeData{Envelope: env, Tab: r.URL.Query().Get("tab"), Endpoints: endpointService.ListEndpoint()})
	}
}

func EnvelopeHTMLGet(envelopeService envelope.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env, err := envelopeService.GetEnvelope(r.Context(), parseID(r))
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
}

func EnvelopeDelete(envelopeService envelope.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := envelopeService.DeleteEnvelope(r.Context(), parseID(r)); err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, core.ErrMessageNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, err.Error(), code)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func AttachmentList(envelopeService envelope.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := parsePage(r)
		atts, err := envelopeService.ListAttachment(r.Context(), &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		view.Render(w, http.StatusOK, view.AttachmentsPage, view.AttachmentsData{Attachments: atts, Page: page})
	}
}

func EnvelopeSendPost(envelopeService envelope.Service, endpointService endpoint.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		endpointFormValue := r.PostFormValue("endpoint")
		filterFormValue := r.PostFormValue("filter")
		noTextFormValue := filterFormValue == "no text"
		noAttachmentsFormValue := filterFormValue == "no attachments"

		// Get envelope
		ctx := r.Context()
		env, err := envelopeService.GetEnvelope(ctx, parseID(r))
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, core.ErrMessageNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, err.Error(), code)
			return
		}

		// Get endpoint
		end, err := endpointService.GetEndpoint(endpointFormValue)
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
		if !noTextFormValue {
			text, err = end.Text(env)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Convert envelope attachments to endpoint attachments
		atts := []endpoint.Attachment{}
		if !noAttachmentsFormValue {
			for _, att := range env.Attachments {
				data, err := envelopeService.GetData(ctx, &att)
				if err != nil {
					if errors.Is(err, core.ErrDataNotFound) {
						continue
					}
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
}
