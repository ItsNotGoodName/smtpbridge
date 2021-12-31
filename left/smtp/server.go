package smtp

import (
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// enableMechLogin enables the LOGIN mechanism which is used for legacy devices.
func enableMechLogin(s *smtp.Server, be smtp.Backend) {
	// Taken from https://github.com/emersion/go-smtp/issues/41#issuecomment-493601465
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

func New(b Backend, config domain.ConfigSMTP) SMTP {
	s := smtp.NewServer(b)

	enableMechLogin(s, b)

	s.Addr = config.Host + ":" + config.PortStr
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = config.Size
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	return SMTP{s}
}

func (s SMTP) Start() {
	log.Println("smtp.SMTP.Start: SMTP server listening on", s.s.Addr)
	if err := s.s.ListenAndServe(); err != nil {
		log.Fatalln("smtp.SMTP.Start:", err)
	}
}
