package procs

import (
	"fmt"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
)

var ErrorLogin = fmt.Errorf("login error")

func AuthHTTPAnonymous(cc core.Context) bool {
	return cc.Config.AuthHTTP.Username == "" && cc.Config.AuthHTTP.Password == ""
}

func AuthHTTPLogin(cc core.Context, username, password string) error {
	if AuthHTTPAnonymous(cc) {
		return nil
	}

	if strings.ToLower(username) != cc.Config.AuthHTTP.Username {
		return ErrorLogin
	}

	if strings.ToLower(password) != cc.Config.AuthHTTP.Password {
		return ErrorLogin
	}

	return nil
}

func SMTPLogin(cc core.Context, username, password string) error {
	if cc.Config.AuthSMTP.Username == "" && cc.Config.AuthSMTP.Password == "" {
		return nil
	}

	if strings.ToLower(username) != cc.Config.AuthSMTP.Username {
		return ErrorLogin
	}

	if strings.ToLower(password) != cc.Config.AuthSMTP.Password {
		return ErrorLogin
	}

	return nil
}
