package auth

import "fmt"

var ErrAuthInvalid = fmt.Errorf("invalid credentials")

type Service interface {
	Login(username, password string) error
}
