package api

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/go-chi/chi/v5"
)

func MessagesGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		q := r.URL.Query()
		cursor, _ := strconv.ParseInt(q.Get("cursor"), 10, 64)
		ascending := q.Get("ascending") == "true"
		limit, _ := strconv.Atoi(q.Get("limit"))

		code := http.StatusOK
		res, err := a.MessageList(r.Context(), &dto.MessageListRequest{
			Cursor:    cursor,
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

func MessageGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

		code := http.StatusOK
		res, err := a.MessageGet(r.Context(), &dto.MessageGetRequest{ID: id})
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
func MessageDelete(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

		code := http.StatusOK
		err := a.MessageDelete(r.Context(), &dto.MessageDeleteRequest{ID: id})
		if err != nil {
			code = http.StatusInternalServerError
		}

		return Response{
			Error: err,
			Code:  code,
		}
	}
}

func MessageEventsGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		q := r.URL.Query()
		cursor, _ := strconv.ParseInt(q.Get("cursor"), 10, 64)
		ascending := q.Get("ascending") == "true"
		limit, _ := strconv.Atoi(q.Get("limit"))

		code := http.StatusOK
		res, err := a.MessageEventList(r.Context(), &dto.EventListRequest{
			EntityID:  id,
			Cursor:    cursor,
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
