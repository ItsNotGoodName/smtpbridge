package service

import (
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Auth struct {
	anonymous bool
	username  string
	password  string
}

func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		anonymous: !cfg.SMTP.Auth,
		username:  strings.ToLower(cfg.SMTP.Username),
		password:  cfg.SMTP.Password,
	}
}

func (m *Auth) Login(username, password string) error {
	if m.anonymous {
		return nil
	}

	if strings.ToLower(username) != m.username || password != m.password {
		return core.ErrAuthInvalid
	}

	return nil
}

func (m *Auth) AnonymousLogin() bool {
	return m.anonymous
}
