package service

type MockAuth struct{}

func NewMockAuth() *MockAuth {
	return &MockAuth{}
}

func (MockAuth) Login(username, password string) error {
	return nil
}

func (MockAuth) AnonymousLogin() bool {
	return true
}
