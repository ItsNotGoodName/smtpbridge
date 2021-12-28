package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Message struct {
	bridgeSVC      app.BridgeServicePort
	endpointREPO   app.EndpointRepositoryPort
	messageREPO    app.MessageRepositoryPort
	attachmentREPO app.AttachmentRepositoryPort
}

func NewMessage(bridgeSVC app.BridgeServicePort, endpointREPO app.EndpointRepositoryPort, messageREPO app.MessageRepositoryPort, attachmentREPO app.AttachmentRepositoryPort) *Message {
	return &Message{
		bridgeSVC:      bridgeSVC,
		endpointREPO:   endpointREPO,
		messageREPO:    messageREPO,
		attachmentREPO: attachmentREPO,
	}
}

func (m *Message) Create(subject, from string, to map[string]bool, text string) (*app.Message, error) {
	msg := app.NewMessage(subject, from, to, text)

	err := m.messageREPO.CreateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, err
}

func (m *Message) CreateAttachment(msg *app.Message, name string, data []byte) (*app.Attachment, error) {
	att, err := app.NewAttachment(msg, name, data)
	if err != nil {
		return nil, err
	}

	err = m.attachmentREPO.CreateAttachment(att)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (m *Message) send(msg *app.EndpointMessage, endpoint app.EndpointPort) {
	err := endpoint.Send(msg)
	if err != nil {
		log.Print("service.Message.send:", err)
	}
}

func (m *Message) SendBridges(msg *app.Message, bridges []app.Bridge) error {
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if !emsg.IsEmpty() {
			for _, name := range bridge.Endpoints {
				endpoint, err := m.endpointREPO.Get(name)
				if err != nil {
					return err
				}
				go m.send(emsg, endpoint)
			}
		}
	}

	return nil
}

func (m *Message) Send(msg *app.Message) error {
	bridges := m.bridgeSVC.GetBridges(msg)
	if len(bridges) == 0 {
		return app.ErrBridgesNotFound
	}

	return m.SendBridges(msg, bridges)
}
