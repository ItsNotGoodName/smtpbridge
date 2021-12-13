package service

import (
	"errors"

	"github.com/ItsNotGoodName/go-smtpbridge/app"
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

func (m *Message) Handle(msg *app.Message) error {
	defer m.messageREPO.Update(msg)

	endpoints := m.bridgeSVC.GetEndpoints(msg)
	if len(endpoints) == 0 {
		msg.Status = app.StatusNoMatch
		return errors.New("no endpoints found")
	}

	// TODO: wait group and log errors
	var errs []error
	for _, endpoint := range endpoints {
		err := endpoint.Send(msg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == len(endpoints) {
		msg.Status = app.StatusNotSent
		return errs[0]
	}

	if len(errs) > 0 {
		msg.Status = app.StatusPartiallySent
	} else {
		msg.Status = app.StatusSent
	}

	return nil
}
