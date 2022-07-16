package auth

type Service interface {
	Login(username, password string) error
}
