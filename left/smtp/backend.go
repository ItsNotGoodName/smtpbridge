package smtp

import (
	"io"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// Backend implements SMTP server methods.
type Backend struct {
	authSVC     domain.AuthServicePort
	bridgeSVC   domain.BridgeServicePort
	endpointSVC domain.EndpointServicePort
	messageSVC  domain.MessageServicePort
}

func (b Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if !b.authSVC.AnonymousLogin() {
		return nil, smtp.ErrAuthRequired
	}
	return newSession(b.bridgeSVC, b.endpointSVC, b.messageSVC), nil
}

func (b Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.authSVC.Login(username, password); err != nil {
		return nil, err
	}
	return newSession(b.bridgeSVC, b.endpointSVC, b.messageSVC), nil
}

func NewBackend(authSVC domain.AuthServicePort, bridgeSVC domain.BridgeServicePort, endpointSVC domain.EndpointServicePort, messageSVC domain.MessageServicePort) Backend {
	return Backend{authSVC, bridgeSVC, endpointSVC, messageSVC}
}

// A session is returned after EHLO.
type session struct {
	bridgeSVC   domain.BridgeServicePort
	endpointSVC domain.EndpointServicePort
	messageSVC  domain.MessageServicePort
	from        string
	to          string
}

func newSession(bridgeSVC domain.BridgeServicePort, endpointSVC domain.EndpointServicePort, messageSVC domain.MessageServicePort) *session {
	return &session{bridgeSVC: bridgeSVC, endpointSVC: endpointSVC, messageSVC: messageSVC}
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

	msg, err := s.messageSVC.Create(e.GetHeader("Subject"), s.from, toMap, e.Text)
	if err != nil {
		log.Println("ERROR: could not create message:", err)
		return err
	}

	for _, a := range e.Attachments {
		if _, err := s.messageSVC.CreateAttachment(msg, a.FileName, a.Content); err != nil {
			log.Println("ATTACHMENT_ERROR: could not add attachment:", err)
			return err
		}
	}

	go func() {
		err := s.endpointSVC.SendBridges(msg, s.bridgeSVC.GetBridges(msg))
		if err != nil {
			log.Println("ERROR: could not send message:", err)
		}
	}()

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
