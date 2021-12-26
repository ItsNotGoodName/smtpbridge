package smtp

import (
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// enableMechLogin enables the LOGIN mechanism which is used for legacy devices.
func enableMechLogin(s *smtp.Server, be smtp.Backend) {
	s.EnableAuth(sasl.Login, func(conn *smtp.Conn) sasl.Server {
		return sasl.NewLoginServer(func(username, password string) error {
			state := conn.State()
			session, err := be.Login(&state, username, password)
			if err != nil {
				return err
			}

			conn.SetSession(session)
			return nil
		})
	})
}

type SMTP struct {
	s *smtp.Server
}

func New(authSVC app.AuthServicePort, messageSVC app.MessageServicePort, config app.ConfigSMTP) SMTP {
	b := newBackend(authSVC, messageSVC)
	s := smtp.NewServer(b)

	enableMechLogin(s, b)

	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = config.Size
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	return SMTP{s}
}

func (s SMTP) Start(address string) {
	s.s.Addr = address
	log.Println("smtp.SMTP.Start: starting SMTP server on", address)
	if err := s.s.ListenAndServe(); err != nil {
		log.Fatalf("smtp.SMTP.Start: %s", err)
	}
}
