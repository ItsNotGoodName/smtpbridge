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

func (m *Message) UpdateStatus(msg *app.Message, status app.Status) error {
	_, err := m.messageREPO.UpdateMessage(msg.UUID, func(msg *app.Message) (*app.Message, error) {
		msg.Status = status
		return msg, nil
	})
	return err
}

func (m *Message) SendBridges(msg *app.Message, bridges []app.Bridge) error {
	if len(bridges) == 0 {
		return app.ErrBridgesNotFound
	}

	var errs []error
	sentCount := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if !emsg.IsEmpty() {
			for _, name := range bridge.Endpoints {
				endpoint, err := m.endpointREPO.Get(name)
				if err != nil {
					return err
				}
				err = endpoint.Send(emsg)
				if err != nil {
					errs = append(errs, err)
				} else {
					sentCount++
				}
			}
		}
	}

	for _, err := range errs {
		log.Println("service.Message.SendBridges:", err)
	}

	if sentCount == 0 {
		if err := m.UpdateStatus(msg, app.StatusUnsent); err != nil {
			return err
		}

		return app.ErrEndpointSendFailed
	}

	if len(errs) > 0 {
		return m.UpdateStatus(msg, app.StatusPartial)
	}

	return m.UpdateStatus(msg, app.StatusSent)
}

func (m *Message) Send(msg *app.Message) error {
	return m.SendBridges(msg, m.bridgeSVC.GetBridges(msg))
}
