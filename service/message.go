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

func (m *Message) UpdateStatus(msg *core.Message, status core.Status) error {
	return m.messageREPO.Update(msg, func(msg *core.Message) (*core.Message, error) {
		msg.Status = status
		return msg, nil
	})
}
