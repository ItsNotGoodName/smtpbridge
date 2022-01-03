package core

import "fmt"

var ErrAuthInvalid = fmt.Errorf("invalid credentials")

type AuthServicePort interface {
	AnonymousLogin() bool
	Login(username, password string) error
}
