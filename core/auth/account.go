package auth

import (
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type AccountService struct {
	username string
	password string
}

func NewAccountService(username, password string) *AccountService {
	return &AccountService{
		username: strings.ToLower(username),
		password: password,
	}
}

func (as *AccountService) Login(username, password string) error {
	if strings.ToLower(username) != as.username || password != as.password {
		return core.ErrAuthInvalid
	}

	return nil
}
