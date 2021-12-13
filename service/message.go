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

func (m *Message) Handle(msg *app.Message) error {
	// TODO: move creation to a seperate function
	m.messageREPO.Create(msg)

	endpoints, err := m.bridgeSVC.GetEndpoints(msg)
	if err != nil {
		return err
	}
	if len(endpoints) == 0 {
		return errors.New("no endpoints found")
	}

	var errs []error
	for _, endpoint := range endpoints {
		err := endpoint.Send(msg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == len(endpoints) {
		return errs[0]
	}

	// TODO: update message

	return nil
}
