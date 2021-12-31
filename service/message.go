package service

import "github.com/ItsNotGoodName/smtpbridge/domain"

type Message struct {
	dao domain.DAO
}

func NewMessage(dao domain.DAO) *Message {
	return &Message{dao}
}

func (m *Message) Get(uuid string) (*domain.Message, error) {
	msg, err := m.dao.Message.GetMessage(uuid)
	if err != nil {
		return nil, err
	}

	msg.Attachments, err = m.dao.Attachment.GetAttachments(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) List(limit, offset int) ([]domain.Message, error) {
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

func (m *Message) Create(subject, from string, to map[string]bool, text string) (*domain.Message, error) {
	msg := domain.NewMessage(subject, from, to, text)

	err := m.dao.Message.CreateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) CreateAttachment(msg *domain.Message, name string, data []byte) (*domain.Attachment, error) {
	att, err := domain.NewAttachment(msg, name, data)
	if err != nil {
		return nil, err
	}

	err = m.dao.Attachment.CreateAttachment(att)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (m *Message) UpdateStatus(msg *domain.Message, status domain.Status) error {
	return m.dao.Message.UpdateMessage(msg, func(msg *domain.Message) (*domain.Message, error) {
		msg.Status = status
		return msg, nil
	})
}
