package database

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Memory struct {
	attMu sync.RWMutex
	att   map[string][]domain.Attachment
	msgMu sync.RWMutex
	msg   map[string]domain.Message
}

func NewMemory() *Memory {
	return &Memory{
		attMu: sync.RWMutex{},
		att:   make(map[string][]domain.Attachment),
		msgMu: sync.RWMutex{},
		msg:   make(map[string]domain.Message),
	}
}

func (db *Memory) CreateMessage(msg *domain.Message) error {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()
	if _, ok := db.msg[msg.UUID]; ok {
		return domain.ErrMessageAlreadyExists
	}

	db.msg[msg.UUID] = *msg

	return nil
}

func (db *Memory) UpdateMessage(msg *domain.Message, updateFN func(msg *domain.Message) (*domain.Message, error)) error {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()

	dbMSG, ok := db.msg[msg.UUID]
	if !ok {
		return domain.ErrMessageNotFound
	}

	updatedMsg, err := updateFN(&dbMSG)
	if err != nil {
		return err
	}

	db.msg[msg.UUID] = *updatedMsg

	return nil
}

func (db *Memory) GetAttachment(uuid string) (*domain.Attachment, error) {
	return nil, domain.ErrNotImplemented
}

func (db *Memory) GetMessage(uuid string) (*domain.Message, error) {
	db.msgMu.RLock()
	msg := db.msg[uuid]
	db.msgMu.RUnlock()
	atts, err := db.GetAttachmentsByMessage(&msg)
	if err != nil {
		return nil, err
	}

	msg.Attachments = atts

	return &msg, nil
}

func (db *Memory) CreateAttachment(att *domain.Attachment) error {
	db.attMu.Lock()
	atts, ok := db.att[att.MessageUUID]
	if !ok {
		atts = []domain.Attachment{}
	}

	atts = append(atts, *att)
	db.att[att.MessageUUID] = atts
	db.attMu.Unlock()

	return nil
}

func (db *Memory) LoadAttachment(msg *domain.Message) error {
	atts, err := db.GetAttachmentsByMessage(msg)
	if err != nil {
		return err
	}

	msg.Attachments = atts

	return nil
}

func (db *Memory) GetAttachmentsByMessage(message *domain.Message) ([]domain.Attachment, error) {
	db.attMu.RLock()
	atts, ok := db.att[message.UUID]
	if !ok {
		atts = []domain.Attachment{}
	}
	db.attMu.RUnlock()

	return atts, nil
}
