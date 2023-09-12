package smtp

import (
	"bytes"
	"context"
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

	return &session{
		app:     b.app,
		log:     log,
		tracer:  tracer,
		address: address,
	}, nil
}

// A Session is returned after EHLO.
type session struct {
	app     core.App
	log     zerolog.Logger
	tracer  trace.Tracer
	address string
	from    string
	to      string
}

func (s *session) AuthPlain(username, password string) error {
	return s.app.AuthSMTPLogin(context.Background(), username, password)
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// log.Println("Rcpt to:", to)
	s.to = to
	return nil
}

func (s *session) Data(r io.Reader) error {
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

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
