package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Message struct {
	bridgeSVC    app.BridgeServicePort
	endpointREPO app.EndpointRepositoryPort
	messageREPO  app.MessageRepositoryPort
}

func NewMessage(bridgeSVC app.BridgeServicePort, endpointREPO app.EndpointRepositoryPort, messageREPO app.MessageRepositoryPort) *Message {
	return &Message{
		bridgeSVC:    bridgeSVC,
		endpointREPO: endpointREPO,
		messageREPO:  messageREPO,
	}
}

func (m *Message) Create(subject, from string, to map[string]bool, text string) (*app.Message, error) {
	msg := app.NewMessage(subject, from, to, text)

	err := m.messageREPO.Create(msg)
	if err != nil {
		return nil, err
	}

	return msg, err
}

func (m *Message) AddAttachment(msg *app.Message, name string, data []byte) error {
	attachment, err := app.NewAttachment(name, data)
	if err != nil {
		return err
	}

	msg.Attachments = append(msg.Attachments, attachment)
	err = m.messageREPO.Update(msg)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) send(msg *app.EndpointMessage, endpoint app.EndpointPort) {
	err := endpoint.Send(msg)
	if err != nil {
		log.Print("service.Message.send:", err)
	}
}

func (m *Message) Send(msg *app.Message) error {
	bridges := m.bridgeSVC.GetBridges(msg)
	if len(bridges) == 0 {
		return app.ErrBridgesNotFound
	}

	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		for _, name := range bridge.Endpoints {
			endpoint, err := m.endpointREPO.Get(name)
			if err != nil {
				return err
			}
			go m.send(emsg, endpoint)
		}
	}

	return nil
}
