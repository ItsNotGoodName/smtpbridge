package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Message struct {
	endpointSVC    core.EndpointServicePort
	attachmentREPO core.AttachmentRepositoryPort
	messageREPO    core.MessageRepositoryPort
}

func NewMessage(
	endpointSVC core.EndpointServicePort,
	attachmentREPO core.AttachmentRepositoryPort,
	messageREPO core.MessageRepositoryPort,
) *Message {
	return &Message{
		endpointSVC:    endpointSVC,
		attachmentREPO: attachmentREPO,
		messageREPO:    messageREPO,
	}
}

func (m *Message) Get(uuid string) (*core.Message, error) {
	msg, err := m.messageREPO.Get(uuid)
	if err != nil {
		return nil, err
	}

	msg.Attachments, err = m.attachmentREPO.ListByMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) List(limit, offset int) ([]core.Message, error) {
	messages, err := m.messageREPO.List(limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		messages[i].Attachments, err = m.attachmentREPO.ListByMessage(&messages[i])
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (m *Message) Create(subject, from string, to map[string]struct{}, text string) (*core.Message, error) {
	msg := core.NewMessage(subject, from, to, text)

	err := m.messageREPO.Create(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) CreateAttachment(msg *core.Message, name string, data []byte) (*core.Attachment, error) {
	att, err := msg.NewAttachment(name, data)
	if err != nil {
		return nil, err
	}

	err = m.attachmentREPO.Create(att)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (m *Message) UpdateStatus(msg *core.Message, status core.Status) error {
	return m.messageREPO.Update(msg, func(msg *core.Message) (*core.Message, error) {
		msg.Status = status
		return msg, nil
	})
}

func (m *Message) Process(msg *core.Message, bridges []*core.Bridge) error {
	skipped := 0
	failed := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if emsg.IsEmpty() {
			skipped++
			continue
		}

		if err := m.endpointSVC.SendByEndpointNames(emsg, bridge.Endpoints); err != nil {
			failed++
			log.Println("service.Message.Process:", err)
		}
	}

	length := len(bridges)
	status := core.StatusSent
	if skipped == length {
		status = core.StatusSkipped
	} else if failed+skipped == len(bridges) {
		status = core.StatusFailed
	}

	return m.UpdateStatus(msg, status)
}
