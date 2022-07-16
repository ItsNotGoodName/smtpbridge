package router

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

type Handler struct {
	envelopeService envelope.Service
}

func NewHandler(envelopeSerivce envelope.Service) *Handler {
	return &Handler{
		envelopeService: envelopeSerivce,
	}
}

func (h *Handler) Index(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageQ, _ := strconv.Atoi(q.Get("page"))
	limitQ, _ := strconv.Atoi(q.Get("limit"))
	ascendingQ, _ := strconv.ParseBool(q.Get("ascending"))

	page := paginate.NewPage(pageQ, limitQ, ascendingQ)
	envs, err := h.envelopeService.ListEnvelope(r.Context(), &page)
	if err != nil {
		renderError(rw, http.StatusInternalServerError, err)
		return
	}

	renderJSON(rw, http.StatusOK, struct {
		Page      paginate.Page
		Envelopes []envelope.Envelope
	}{page, envs})
}
