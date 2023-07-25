package models

import "strings"

type Auth struct {
	Anonymous bool
	Username  string
	Password  string
}

func NewAuth(username, password string) Auth {
	anonymous := false
	if username == "" && password == "" {
		anonymous = true
	}

	return Auth{
		Anonymous: anonymous,
		Username:  username,
		Password:  password,
	}
}

func (a Auth) Check(username, password string) bool {
	return strings.ToLower(username) == strings.ToLower(a.Username) && password == a.Password
}
