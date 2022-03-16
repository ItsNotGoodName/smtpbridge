package api

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func InfoGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		code := http.StatusOK
		res, err := a.Info(r.Context())
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
