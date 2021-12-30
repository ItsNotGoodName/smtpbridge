package service

import (
	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Message struct {
	dao app.DAO
}

func NewMessage(dao app.DAO) *Message {
	return &Message{dao}
}

func (m *Message) List(limit, offset int) ([]app.Message, error) {
	messages, err := m.dao.Message.GetMessages(limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		messages[i].Attachments, err = m.dao.Attachment.GetAttachments(&messages[i])
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (m *Message) Create(subject, from string, to map[string]bool, text string) (*app.Message, error) {
	msg := app.NewMessage(subject, from, to, text)

	err := m.dao.Message.CreateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) CreateAttachment(msg *app.Message, name string, data []byte) (*app.Attachment, error) {
	att, err := app.NewAttachment(msg, name, data)
	if err != nil {
		return nil, err
	}

	err = m.dao.Attachment.CreateAttachment(att)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (m *Message) UpdateStatus(msg *app.Message, status app.Status) error {
	return m.dao.Message.UpdateMessage(msg, func(msg *app.Message) (*app.Message, error) {
		msg.Status = status
		return msg, nil
	})
}
