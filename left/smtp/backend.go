package smtp

import (
	"context"
	"io"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// Backend implements SMTP server methods.
type Backend struct {
	app dto.App
}

func (b Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if err := b.app.SMTPLogin(context.Background(), &dto.SMTPLoginRequest{}); err != nil {
		log.Println("smtp.AnonymousLogin: login failure:", err)
		return nil, smtp.ErrAuthRequired
	}
	return newSession(b.app), nil
}

func (b Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.app.SMTPLogin(context.Background(), &dto.SMTPLoginRequest{Username: username, Password: password}); err != nil {
		log.Println("smtp.Login: login failure:", err)
		return nil, err
	}
	return newSession(b.app), nil
}

func NewBackend(app dto.App) Backend {
	return Backend{app}
}

// A session is returned after EHLO.
type session struct {
	app  dto.App
	from string
	to   string
}

func newSession(app dto.App) *session {
	return &session{app: app}
}

func (s *session) Mail(from string, opts smtp.MailOptions) error {
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
		log.Println("smtp.session.Data: could not read email:", err)
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
	toMap := make(map[string]struct{})
	if to, err := e.AddressList("To"); err == nil {
		for _, t := range to {
			toMap[t.Address] = struct{}{}
			//log.Println("TO:", t.Address)
		}
	} else {
		log.Println("smtp.session.Data: could not get To from email:", err)
	}
	toMap[s.to] = struct{}{}

	req := dto.MessageHandleRequest{
		Subject: e.GetHeader("Subject"),
		From:    s.from,
		To:      toMap,
		Text:    e.Text,
	}
	for _, a := range e.Attachments {
		req.AddAttachment(a.FileName, a.Content)
	}

	go func() {
		err := s.app.MessageHandle(context.Background(), &req)
		if err != nil {
			log.Println("smtp.session.Data: could not handle message:", err)
		}
	}()

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
