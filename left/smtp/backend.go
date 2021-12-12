package smtp

import (
	"io"
	"log"

	"github.com/ItsNotGoodName/go-smtpbridge/app"
	"github.com/ItsNotGoodName/go-smtpbridge/port"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// backend implements SMTP server methods.
type backend struct {
	authSVC    port.AuthService
	messageSVC port.MessageService
}

func (b backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

func (b backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.authSVC.Login(username, password); err != nil {
		return nil, err
	}
	return NewSession(b.messageSVC), nil
}

func newBackend(auth port.AuthService, messageSVC port.MessageService) *backend {
	return &backend{auth, messageSVC}
}

// A session is returned after EHLO.
type session struct {
	messageSVC port.MessageService
	from       string
}

func NewSession(messageSVC port.MessageService) *session {
	return &session{messageSVC, ""}
}

func (s *session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	}

	log.Println("SUBJECT:", e.GetHeader("Subject"))
	log.Println("TEXT:", e.Text)
	log.Println("HTML:", e.HTML)
	log.Println("ATTACHMENTS:", len(e.Attachments))
	for e := range e.Errors {
		log.Println("ERROR:", e)
	}
	log.Println("FROM:", e.GetHeader("From"))
	toMap := make(map[string]bool)
	if to, err := e.AddressList("To"); err == nil {
		for _, t := range to {
			toMap[t.Address] = true
			log.Println("TO:", t.Address)
		}
	} else {
		log.Println("TO_ERROR:", err)
	}

	return s.messageSVC.Handle(&app.Message{Subject: e.GetHeader("Subject"), Text: e.Text, From: s.from, To: toMap})
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
