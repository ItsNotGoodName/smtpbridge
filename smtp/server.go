package smtp

import (
	"context"
	"net"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

// enableMechLogin enables the LOGIN mechanism which is used for legacy devices.
func enableMechLogin(be smtp.Backend, s *smtp.Server) {
	// Adapted from https://github.com/emersion/go-smtp/issues/41#issuecomment-493601465
	s.EnableAuth(sasl.Login, func(conn *smtp.Conn) sasl.Server {
		return sasl.NewLoginServer(func(username, password string) error {
			session, err := be.NewSession(conn)
			if err != nil {
				return err
			}

			return session.AuthPlain(username, password)
		})
	})
}

type SMTP struct {
	server   *smtp.Server
	shutdown context.CancelFunc
}

func New(app core.App, shutdown context.CancelFunc, addr string, maxMessageBytes int) SMTP {
	be := newBackend(app)
	server := smtp.NewServer(be)

	enableMechLogin(be, server)

	server.Addr = addr
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = maxMessageBytes
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	return SMTP{
		server:   server,
		shutdown: shutdown,
	}
}

func (s SMTP) Start() (net.Listener, error) {
	log.Info().Msg("Starting SMTP server on " + s.server.Addr)

	network := "tcp"
	if s.server.LMTP {
		network = "unix"
	}

	addr := s.server.Addr
	if !s.server.LMTP && addr == "" {
		addr = ":smtp"
	}

	l, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	go s.server.Serve(l)

	return l, nil
}

func (s SMTP) Background(ctx context.Context, doneC chan<- struct{}) {
	listener, err := s.Start() // Start SMTP server
	if err != nil {
		log.Error().Err(err).Msg("Failed to start SMTP server")
		doneC <- struct{}{} // Done
		s.shutdown()
		return
	}

	<-ctx.Done() // Wait for context
	log.Info().Msg("Gracefully shutting down SMTP server...")
	listener.Close()    // Shutdown SMTP server
	doneC <- struct{}{} // Done
}
