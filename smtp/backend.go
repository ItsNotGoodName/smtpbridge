package smtp

import (
	"bytes"
	"context"
	"io"
	"net/mail"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The Backend implements SMTP server methods.
type Backend struct {
	app core.App
	log zerolog.Logger
}

func NewBackend(app core.App) *Backend {
	log := log.With().Str("source", "smtp").Logger()
	return &Backend{
		app: app,
		log: log,
	}
}

func (b *Backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	address := state.Conn().RemoteAddr().String()
	tracer := b.app.Tracer(trace.SourceSMTP).Sticky(trace.WithAddress(address))
	log := b.log.With().Str("address", address).Logger()

	// log.Debug().Msg("NewSession")

	return &session{
		app:     b.app,
		auth:    b.app.AuthSMTPAnonymous(),
		log:     log,
		tracer:  tracer,
		address: address,
	}, nil
}

// A Session is returned after EHLO.
type session struct {
	auth    bool
	app     core.App
	log     zerolog.Logger
	tracer  trace.Tracer
	address string
	from    string
	to      string
}

func (s *session) AuthPlain(username, password string) error {
	// s.log.Debug().Str("username", username).Str("password", password).Msg("AuthPlain")
	err := s.app.AuthSMTPLogin(context.Background(), username, password)
	if err != nil {
		return smtp.ErrAuthFailed
	}

	s.auth = true

	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// s.log.Debug().Str("from", from).Msg("Mail")
	if !s.auth {
		return smtp.ErrAuthRequired
	}

	s.from = from

	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// s.log.Debug().Str("to", to).Msg("Rcpt")
	if !s.auth {
		return smtp.ErrAuthRequired
	}

	s.to = to

	return nil
}

func (s *session) Data(r io.Reader) error {
	// s.log.Debug().Msg("Data")
	if !s.auth {
		return smtp.ErrAuthRequired
	}

	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to read envelope")
		return err
	}

	to := []string{s.to}
	if addresses, err := e.AddressList("To"); err == nil {
		for _, t := range addresses {
			to = append(to, t.Address)
		}
	} else {
		s.log.Warn().Err(err).Msg("Failed to get 'To' from address list")
	}

	date, err := e.Date()
	if err != nil && err != mail.ErrHeaderNotPresent {
		s.log.Warn().Err(err).Str("data", e.GetHeader("Date")).Msg("Failed to parse date")
	}

	datts := []models.DTOAttachmentCreate{}
	for _, a := range e.Attachments {
		datts = append(datts, models.DTOAttachmentCreate{
			Name: a.FileName,
			Data: bytes.NewBuffer(a.Content),
		})
	}
	msg := models.DTOMessageCreate{
		From:    s.from,
		To:      to,
		Subject: e.GetHeader("Subject"),
		Text:    e.Text,
		HTML:    e.HTML,
		Date:    date,
	}

	// Create envelope
	ctx := context.Background()
	id, err := s.app.EnvelopeCreate(ctx, msg, datts)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to create envelope")
		return err
	}

	s.tracer.Trace(ctx, trace.ActionEnvelopeCreated, trace.WithEnvelope(id))
	s.log.Info().Int64("envelope-id", id).Msg("Envelope created")

	return nil
}

func (s *session) Reset() {
	// s.log.Debug().Msg("Reset")
}

func (s *session) Logout() error {
	// s.log.Debug().Msg("Logout")
	return nil
}

func (s *session) AuthMechanisms() []string {
	return []string{
		sasl.Login,
	}
}

func (s *session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewLoginServer(func(username, password string) error {
		return s.AuthPlain(username, password)
	}), nil
}
