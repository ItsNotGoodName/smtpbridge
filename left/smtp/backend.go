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
	authSVC     app.AuthServicePort
	endpointSVC app.EndpointServicePort
	messageSVC  app.MessageServicePort
}

func (b Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if !b.authSVC.AnonymousLogin() {
		return nil, smtp.ErrAuthRequired
	}
	return newSession(b.endpointSVC, b.messageSVC), nil
}

func (b Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.authSVC.Login(username, password); err != nil {
		return nil, err
	}
	return newSession(b.endpointSVC, b.messageSVC), nil
}

func NewBackend(authSVC app.AuthServicePort, endpointSVC app.EndpointServicePort, messageSVC app.MessageServicePort) Backend {
	return Backend{authSVC, endpointSVC, messageSVC}
}

// A session is returned after EHLO.
type session struct {
	messageSVC  app.MessageServicePort
	endpointSVC app.EndpointServicePort
	from        string
	to          string
}

func newSession(endpointSVC app.EndpointServicePort, messageSVC app.MessageServicePort) *session {
	return &session{messageSVC: messageSVC, endpointSVC: endpointSVC}
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
		log.Println("ERROR: could not read email:", err)
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
		log.Println("TO_ERROR: could not get To from email:", err)
	}
	toMap[s.to] = true

	m, err := s.messageSVC.Create(e.GetHeader("Subject"), s.from, toMap, e.Text)
	if err != nil {
		log.Println("ERROR: could not create message:", err)
		return err
	}

	for _, a := range e.Attachments {
		if _, err := s.messageSVC.CreateAttachment(m, a.FileName, a.Content); err != nil {
			log.Println("ATTACHMENT_ERROR: could not add attachment:", err)
			return err
		}
	}

	err = s.endpointSVC.Send(m)
	if err != nil {
		log.Println("ERROR: could not send message:", err)
		return err
	}

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
