package auth

import (
	"strings"
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
		return ErrAuthInvalid
	}

	return nil
}
