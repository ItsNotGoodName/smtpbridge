package smtp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/mail"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
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
	fmt.Println("1.1", s.auth)

	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// s.log.Debug().Str("from", from).Msg("Mail")
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	// log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// s.log.Debug().Str("to", to).Msg("Rcpt")
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	// log.Println("Rcpt to:", to)
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
		log.Error().Err(err).Msg("Failed to read envelope")
		return err
	}

	//log.Println("SUBJECT:", e.GetHeader("Subject"))
	//log.Println("TEXT:", e.Text)
	//log.Println("HTML:", e.HTML)
	//log.Println("ATTACHMENTS:", len(e.Attachments))
	//for e := range e.Errors {
	//	log.Println("ERROR:", e)
	//}
	//log.Println("FROM:", e.GetHeader("From"))
	to := []string{s.to}
	if addresses, err := e.AddressList("To"); err == nil {
		for _, t := range addresses {
			to = append(to, t.Address)
			//log.Println("TO:", t.Address)
		}
	} else {
		log.Warn().Err(err).Msg("Failed to get 'To' from address list")
	}

	date, err := e.Date()
	if err != nil && err != mail.ErrHeaderNotPresent {
		log.Warn().Err(err).Str("data", e.GetHeader("Date")).Msg("Failed to parse date")
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
		log.Error().Err(err).Msg("Failed to create envelope")
		return err
	}

	s.tracer.Trace(ctx, trace.ActionEnvelopeCreated, trace.WithEnvelope(id))

	return nil
}

func (s *session) Reset() {
	// s.log.Debug().Msg("Reset")
}

func (s *session) Logout() error {
	// s.log.Debug().Msg("Logout")
	return nil
}
