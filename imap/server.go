package imap

import (
	"context"
	"fmt"
	"net"

	"github.com/emersion/go-imap/server"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server *server.Server
	addr   string
}

func NewServer(backend Backend, addr string) Server {
	server := server.New(backend)

	server.AllowInsecureAuth = true
	server.ErrorLog = errorLogger{}

	return Server{
		server: server,
		addr:   addr,
	}
}

type errorLogger struct{}

// Printf implements imap.Logger.
func (errorLogger) Printf(format string, v ...interface{}) {
	log.Error().Str("source", "imap").Msgf(format, v...)
}

// Println implements imap.Logger.
func (errorLogger) Println(v ...interface{}) {
	log.Error().Str("source", "imap").Msg(fmt.Sprintln(v...))
}

func (s Server) Serve(ctx context.Context) error {
	log.Info().Str("address", s.addr).Msg("Starting IMAP server")

	l, err := net.Listen("tcp", s.addr)
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

	log.Info().Msg("Gracefully shutting down IMAP server...")
	return nil
}
