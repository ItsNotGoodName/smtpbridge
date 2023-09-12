package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	handler http.Handler
	address string
}

func NewServer(handler http.Handler, address string) Server {
	return Server{
		handler: handler,
		address: address,
	}
}

func (Server) String() string {
	return "http.Server"
}

func (s Server) Serve(ctx context.Context) error {
	server := &http.Server{Addr: s.address, Handler: s.handler}

	go func() {
		<-ctx.Done()

		log.Info().Msg("Gracefully shutting down HTTP server...")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				server.Close()
			}
		}()

		// Trigger graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Err(err).Msg("Failed to shutdown HTTP server")
		}
	}()

	log.Info().Str("address", s.address).Msg("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
