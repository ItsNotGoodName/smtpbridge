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

func (m *Message) Handle(msg *app.Message) error {
	defer m.messageREPO.Update(msg)

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
