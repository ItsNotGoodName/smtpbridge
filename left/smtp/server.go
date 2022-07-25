package smtp

import (
	"context"
	"log"
	"net"
	"time"

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

func New(b Backend, addr string, maxMessageBytes int) SMTP {
	s := smtp.NewServer(b)

	enableMechLogin(s, b)

	s.Addr = addr
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = maxMessageBytes
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	return SMTP{s}
}

func (s SMTP) Start() net.Listener {
	log.Println("smtp.SMTP.Start: SMTP server listening on", s.s.Addr)

	network := "tcp"
	if s.s.LMTP {
		network = "unix"
	}

	addr := s.s.Addr
	if !s.s.LMTP && addr == "" {
		addr = ":smtp"
	}

	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalln("smtp.SMTP.Start:", err)
	}

	go s.s.Serve(l)

	return l
}

func (s SMTP) Run(ctx context.Context, doneC chan<- struct{}) {
	l := s.Start()      // Start SMTP server
	<-ctx.Done()        // Wait for context
	l.Close()           // Shutdown SMTP server
	doneC <- struct{}{} // Done
}
