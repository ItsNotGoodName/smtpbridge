package service

import "github.com/ItsNotGoodName/smtpbridge/domain"

type Message struct {
	attachmentREPO domain.AttachmentRepositoryPort
	messageREPO    domain.MessageRepositoryPort
}

func NewMessage(
	attachmentREPO domain.AttachmentRepositoryPort,
	messageREPO domain.MessageRepositoryPort,
) *Message {
	return &Message{
		attachmentREPO: attachmentREPO,
		messageREPO:    messageREPO,
	}
}

func (m *Message) Get(uuid string) (*domain.Message, error) {
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

func (m *Message) List(limit, offset int) ([]domain.Message, error) {
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

func (m *Message) Create(subject, from string, to map[string]struct{}, text string) (*domain.Message, error) {
	msg := domain.NewMessage(subject, from, to, text)

	err := m.messageREPO.Create(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) CreateAttachment(msg *domain.Message, name string, data []byte) (*domain.Attachment, error) {
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

func (m *Message) UpdateStatus(msg *domain.Message, status domain.Status) error {
	return m.messageREPO.Update(msg, func(msg *domain.Message) (*domain.Message, error) {
		msg.Status = status
		return msg, nil
	})
}
