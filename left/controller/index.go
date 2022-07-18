package controller

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
	"github.com/ItsNotGoodName/smtpbridge/left/view"
)

func (c *Controller) Index(rw http.ResponseWriter, r *http.Request) {
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

	view.Render(rw, http.StatusOK, view.IndexData{Envelopes: envs, Page: page}, view.IndexPage)
}
