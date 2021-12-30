package service

import (
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Auth struct {
	anonymous bool
	username  string
	password  string
}

func NewAuth(config *domain.ConfigAuth) *Auth {
	return &Auth{
		anonymous: !config.Enable,
		username:  strings.ToLower(config.Username),
		password:  config.Password,
	}
}

func (m *Auth) Login(username, password string) error {
	if m.anonymous {
		return nil
	}

	if strings.ToLower(username) != m.username || password != m.password {
		return domain.ErrAuthInvalid
	}

	return nil
}

func (m *Auth) AnonymousLogin() bool {
	return m.anonymous
}
