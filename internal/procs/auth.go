package procs

import (
	"fmt"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
)

var ErrorLogin = fmt.Errorf("login error")

func HTTPLogin(cc core.Context, username, password string) error {
	if cc.Config.AuthHTTP.Username == "" && cc.Config.AuthHTTP.Password == "" {
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
	if cc.Config.AuthSMTP.Username == "" && cc.Config.AuthHTTP.Password == "" {
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
