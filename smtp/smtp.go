package smtp

import (
	"context"
	"net"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

type SMTP struct {
	server *smtp.Server
}

func (SMTP) String() string {
	return "smtp.SMTP"
}

func New(backend smtp.Backend, addr string, maxMessageBytes int64) SMTP {
	server := smtp.NewServer(backend)

	enableMechLogin(backend, server)

	server.Addr = addr
	server.Domain = "localhost"
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = maxMessageBytes
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	return SMTP{
		server: server,
	}
}

func (s SMTP) Serve(ctx context.Context) error {
	log.Info().Str("address", s.server.Addr).Msg("Starting SMTP server")

	l, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	errC := make(chan error, 1)

	go func() {
		err := s.server.Serve(l)
		if err != nil {
			errC <- err
		}
	}()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
	}

	s.server.Close()

	log.Info().Msg("Gracefully shutting down SMTP server...")
	return nil
}
