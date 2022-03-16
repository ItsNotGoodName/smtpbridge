package api

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func VersionGet(a dto.App) Handler {
	return func(rw http.ResponseWriter, r *http.Request) Response {
		return Response{
			Code: http.StatusOK,
			Data: a.Version(),
		}
	}
}
