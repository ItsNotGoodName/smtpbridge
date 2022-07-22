package controller

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
	"github.com/go-chi/chi/v5"
)

func parsePage(r *http.Request) paginate.Page {
	q := r.URL.Query()
	pageQ, _ := strconv.Atoi(q.Get("page"))
	limitQ, _ := strconv.Atoi(q.Get("limit"))
	ascendingQ, _ := strconv.ParseBool(q.Get("ascending"))

	return paginate.NewPage(pageQ, limitQ, ascendingQ)
}

func parseID(r *http.Request) int64 {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	return id
}
