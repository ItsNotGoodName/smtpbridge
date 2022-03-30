package api

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func EventsGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		ascending := q.Get("ascending") == "true"
		limit, _ := strconv.Atoi(q.Get("limit"))

		code := http.StatusOK
		res, err := a.EventList(r.Context(), &dto.EventListRequest{
			Page:      page,
			Ascending: ascending,
			Limit:     limit,
		})
		if err != nil {
			code = http.StatusInternalServerError
		}

		return Response{
			Data:  res,
			Error: err,
			Code:  code,
		}
	}
}
