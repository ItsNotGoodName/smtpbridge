package service

type NoAuth struct{}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (NoAuth) Login(username, password string) error {
	return nil
}
