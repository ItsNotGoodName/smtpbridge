package models

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
