package app

type AuthLoginRequest struct {
	Username string
	Password string
}

func (a *App) AuthLoginRequest(req *AuthLoginRequest) error {
	if a.authSVC.AnonymousLogin() {
		return nil
	}

	return a.authSVC.Login(req.Username, req.Password)
}
