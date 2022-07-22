package smtp

import (
	"context"
	"io"
	"log"
	"net"
	"net/mail"

	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

// Backend implements SMTP server methods.
type Backend struct {
	envelopeService envelope.Service
	authService     auth.Service
}

func (b Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if err := b.authService.Login("", ""); err != nil {
		log.Println("smtp.AnonymousLogin: login failure:", err)
		return nil, smtp.ErrAuthRequired
	}
	return newSession(b.envelopeService, state.RemoteAddr), nil
}

func (b Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if err := b.authService.Login(username, password); err != nil {
		log.Println("smtp.Login: login failure:", err)
		return nil, err
	}
	return newSession(b.envelopeService, state.RemoteAddr), nil
}

func NewBackend(envelopeService envelope.Service, authService auth.Service) Backend {
	return Backend{
		envelopeService: envelopeService,
		authService:     authService,
	}
}

// A session is returned after EHLO.
type session struct {
	envelopeService envelope.Service
	from            string
	to              string
	remoteAddr      net.Addr
}

func newSession(envelopeService envelope.Service, remoteAddr net.Addr) *session {
	return &session{envelopeService: envelopeService, remoteAddr: remoteAddr}
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
		log.Printf("smtp.session.Data: remote %s: could not read envelope: %s", s.remoteAddr, err)
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
		log.Printf("smtp.session.Data: remote %s: could not get To from address list: %s", s.remoteAddr, err)
	}

	// Attachment requests
	attsReq := []envelope.CreateAttachmentRequest{}
	for _, a := range e.Attachments {
		attsReq = append(attsReq, envelope.CreateAttachmentRequest{
			Name: a.FileName,
			Data: a.Content,
		})
	}

	// Envelope request
	date, err := e.Date()
	if err != nil && err != mail.ErrHeaderNotPresent {
		log.Printf("smtp.session.Data: remote %s: could not parse date: %s: %s", s.remoteAddr, e.GetHeader("Date"), err)
	}

	envReq := &envelope.CreateEnvelopeRequest{
		From:       s.from,
		To:         to,
		Subject:    e.GetHeader("Subject"),
		Text:       e.Text,
		HTML:       e.HTML,
		Attachment: attsReq,
		Date:       date,
	}

	// Create envelope
	envID, err := s.envelopeService.CreateEnvelope(context.Background(), envReq)
	if err != nil {
		log.Printf("smtp.session.Data: remote %s: %s", s.remoteAddr, err)
		return nil
	}

	log.Printf("smtp.session.Data: remote %s: id %d: envelope created", s.remoteAddr, envID)
	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
