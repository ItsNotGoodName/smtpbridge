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

func (m *Message) AddAttachment(msg *app.Message, name string, data []byte) error {
	att, err := app.NewAttachment(msg, name, data)
	if err != nil {
		return err
	}

	return m.attachmentREPO.CreateAttachment(att, data)
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

	datts, err := m.attachmentREPO.GetDataAttachmentsByMessage(msg)
	if err != nil {
		return err
	}

	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg, datts)
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
