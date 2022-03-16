package app

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func (a *App) SMTPLogin(ctx context.Context, req *dto.SMTPLoginRequest) error {
	return a.smtpAuthService.Login(req.Username, req.Password)
}
