package service

import (
	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Message struct {
	attachmentREPO core.AttachmentRepositoryPort
	messageREPO    core.MessageRepositoryPort
	size           int64
}

func NewMessage(
	cfg *config.Config,
	attachmentREPO core.AttachmentRepositoryPort,
	messageREPO core.MessageRepositoryPort,
) *Message {
	return &Message{
		attachmentREPO: attachmentREPO,
		messageREPO:    messageREPO,
		size:           cfg.DB.Size,
	}
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
	messages, err := m.messageREPO.List(limit, offset, true)
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

func (m *Message) LoadData(msg *core.Message) error {
	for i := range msg.Attachments {
		var err error
		msg.Attachments[i].Data, err = m.attachmentREPO.GetData(&msg.Attachments[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Message) UpdateStatus(msg *core.Message, status core.Status) error {
	return m.messageREPO.Update(msg, func(msg *core.Message) (*core.Message, error) {
		msg.Status = status
		return msg, nil
	})
}
