package controller

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
	"github.com/ItsNotGoodName/smtpbridge/core/version"
	"github.com/ItsNotGoodName/smtpbridge/left/http/view"
)

func IndexGet(envelopeService envelope.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envPage := paginate.NewPage(0, 10, false)
		envs, err := envelopeService.ListEnvelope(r.Context(), &envPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		attPage := paginate.NewPage(0, 0, false)
		_, err = envelopeService.ListAttachment(r.Context(), &attPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		view.Render(w, http.StatusOK, view.IndexPage, view.IndexData{
			EnvelopesCount:   envPage.Count,
			AttachmentsCount: attPage.Count,
			Envelopes:        envs,
			Version:          version.Current,
		})
	}
}
