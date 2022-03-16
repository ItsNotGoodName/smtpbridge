package auth

type AnonymousService struct{}

func NewAnonymousService() *AnonymousService {
	return &AnonymousService{}
}

func (*AnonymousService) Login(username, password string) error {
	return nil
}
