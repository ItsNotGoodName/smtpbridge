package domain

import "fmt"

var ErrAuthInvalid = fmt.Errorf("invalid credentials")

// AuthServicePort handles authenticating users.
type AuthServicePort interface {
	AnonymousLogin() bool
	Login(username, password string) error
}
