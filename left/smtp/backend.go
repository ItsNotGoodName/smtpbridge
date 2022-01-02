package smtp

import (
	"io"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// Backend implements SMTP server methods.
type Backend struct {
	app *app.App
}

func (b Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if err := b.app.AuthLoginRequest(&app.AuthLoginRequest{}); err != nil {
		log.Println("smtp.AnonymousLogin: login failure:", smtp.ErrAuthRequired)
		return nil, smtp.ErrAuthRequired
	}
	return newSession(b.app), nil
}

func (b Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.app.AuthLoginRequest(&app.AuthLoginRequest{Username: username, Password: password}); err != nil {
		log.Println("smtp.Login: login failure:", err)
		return nil, err
	}
	return newSession(b.app), nil
}

func NewBackend(app *app.App) Backend {
	return Backend{app}
}

// A session is returned after EHLO.
type session struct {
	app  *app.App
	from string
	to   string
}

func newSession(app *app.App) *session {
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
		log.Println("smtp.Data: could not read email:", err)
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
	toMap := make(map[string]bool)
	if to, err := e.AddressList("To"); err == nil {
		for _, t := range to {
			toMap[t.Address] = true
			//log.Println("TO:", t.Address)
		}
	} else {
		log.Println("smtp.Data: could not get To from email:", err)
	}
	toMap[s.to] = true

	req := app.MessageCreateRequest{
		Subject: e.GetHeader("Subject"),
		From:    s.from,
		To:      toMap,
		Text:    e.Text,
	}
	for _, a := range e.Attachments {
		req.AddAttachment(a.FileName, a.Content)
	}

	msg, err := s.app.MessageCreate(&req)
	if err != nil {
		log.Println("smtp.Data: could not create message:", err)
		return err
	}

	go func() {
		err := s.app.MessageSend(&app.MessageSendRequest{Message: msg})
		if err != nil {
			log.Println("smtp.Data: could not send message:", err)
		}
	}()

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
