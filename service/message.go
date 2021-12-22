package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Message struct {
	bridgeSVC   app.BridgeServicePort
	messageREPO app.MessageRepositoryPort
}

func NewMessage(bridgeSVC app.BridgeServicePort, messageREPO app.MessageRepositoryPort) *Message {
	return &Message{bridgeSVC, messageREPO}
}

func (m *Message) Create(subject, from string, to map[string]bool, text string) (*app.Message, error) {
	msg := app.NewMessage(subject, from, to, text)
	err := m.messageREPO.Create(msg)
	if err != nil {
		return nil, err
	}
	return msg, err
}

func (m *Message) AddAttachment(msg *app.Message, name string, content []byte) error {
	att, err := app.NewAttachment(name, content)
	if err != nil {
		return err
	}

	msg.Attachments = append(msg.Attachments, att)
	err = m.messageREPO.Update(msg)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) Send(msg *app.Message) error {
	endpoints := m.bridgeSVC.GetEndpoints(msg)
	if len(endpoints) == 0 {
		return app.ErrNoEndpoints
	}

	var errs []error
	for _, endpoint := range endpoints {
		err := endpoint.Send(msg)
		if err != nil {
			errs = append(errs, err)
			log.Println("service.Message.Handle:", err)
		}
	}
	// Return first error if messsage could not be sent to atleast one endpoint.
	if len(errs) == len(endpoints) {
		return errs[0]
	}

	return nil
}
