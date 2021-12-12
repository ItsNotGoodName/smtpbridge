package port

import (
	"github.com/ItsNotGoodName/go-smtpbridge/app"
)

type AuthService interface {
	Login(username, password string) error
}

type MessageService interface {
	Handle(*app.Message) error
}

type MessageRepository interface {
	Save(*app.Message) error
}
