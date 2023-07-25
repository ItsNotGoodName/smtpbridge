package smtp

import (
	"context"
	"io"
	"net"
	"net/mail"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/rs/zerolog/log"
)

// The backend implements SMTP server methods.
type backend struct {
	app core.App
}

func newBackend(app core.App) *backend {
	return &backend{app}
}

func (b *backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	return newSession(b.app.Context(context.Background()), state.Conn().RemoteAddr()), nil
}

// A Session is returned after EHLO.
type session struct {
	cc         core.Context
	from       string
	to         string
	remoteAddr net.Addr
}

func newSession(ctx core.Context, remoteAddr net.Addr) *session {
	return &session{cc: ctx, remoteAddr: remoteAddr}
}

func (s *session) AuthPlain(username, password string) error {
	return procs.AuthSMTPLogin(s.cc, username, password)
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *session) Rcpt(to string) error {
	// log.Println("Rcpt to:", to)
	s.to = to
	return nil
}

func (s *session) Data(r io.Reader) error {
	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Error().Err(err).Str("remoteAddr", s.remoteAddr.String()).Msg("failed to read envelope")
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
		log.Warn().Err(err).Str("remoteAddr", s.remoteAddr.String()).Msg("failed to get To from address list")
	}

	date, err := e.Date()
	if err != nil && err != mail.ErrHeaderNotPresent {
		log.Warn().Err(err).Str("remoteAddr", s.remoteAddr.String()).Str("data", e.GetHeader("Date")).Msg("failed to parse date")
	}

	// Attachment requests
	datts := []envelope.DataAttachment{}
	for _, a := range e.Attachments {
		datts = append(datts, envelope.NewDataAttachment(a.FileName, a.Content))
	}
	msg := envelope.NewMessage(envelope.CreateMessage{
		From:    s.from,
		To:      to,
		Subject: e.GetHeader("Subject"),
		Text:    e.Text,
		HTML:    e.HTML,
		Date:    date,
	})

	// Create envelope
	id, err := procs.EnvelopeCreate(s.cc, msg, datts)
	if err != nil {
		log.Error().Err(err).Str("remoteAddr", s.remoteAddr.String()).Msg("failed to create envelope")
		return err
	}

	log.Info().Str("remoteAddr", s.remoteAddr.String()).Int64("id", id).Msg("envelope created")

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
