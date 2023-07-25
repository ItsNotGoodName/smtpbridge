package procs

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
)

var ErrorLogin = fmt.Errorf("login invalid")

func AuthHTTPAnonymous(cc core.Context) bool {
	return cc.Config.AuthHTTP.Anonymous
}

func AuthHTTPLogin(cc core.Context, username, password string) error {
	if cc.Config.AuthHTTP.Anonymous || cc.Config.AuthHTTP.Check(username, password) {
		return nil
	}

	return ErrorLogin
}

func AuthSMTPLogin(cc core.Context, username, password string) error {
	if cc.Config.AuthSMTP.Anonymous || cc.Config.AuthSMTP.Check(username, password) {
		return nil
	}

	return ErrorLogin
}
