package auth

import (
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

func New(username, password string) models.Auth {
	anonymous := false
	if username == "" && password == "" {
		anonymous = true
	}

	return models.Auth{
		Anonymous: anonymous,
		Username:  username,
		Password:  password,
	}
}

func Check(a models.Auth, username, password string) bool {
	return strings.ToLower(username) == strings.ToLower(a.Username) && password == a.Password
}
