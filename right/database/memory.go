package database

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Memory struct {
	attMu sync.RWMutex
	att   map[string][]app.Attachment
	msgMu sync.RWMutex
	msg   map[string]app.Message
}

func NewMemory() *Memory {
	return &Memory{
		attMu: sync.RWMutex{},
		att:   make(map[string][]app.Attachment),
		msgMu: sync.RWMutex{},
		msg:   make(map[string]app.Message),
	}
}

func (db *Memory) CreateMessage(msg *app.Message) error {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()
	if _, ok := db.msg[msg.UUID]; ok {
		return app.ErrMessageAlreadyExists
	}

	db.msg[msg.UUID] = *msg

	return nil
}

func (db *Memory) UpdateMessage(msg *app.Message, updateFN func(msg *app.Message) (*app.Message, error)) (*app.Message, error) {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()

	dbMSG, ok := db.msg[msg.UUID]
	if !ok {
		return nil, app.ErrMessageNotFound
	}

	updatedMsg, err := updateFN(&dbMSG)
	if err != nil {
		return nil, err
	}

	db.msg[msg.UUID] = *updatedMsg

	return updatedMsg, nil
}

func (db *Memory) GetAttachment(uuid string) (*app.Attachment, error) {
	return nil, app.ErrNotImplemented
}

func (db *Memory) GetMessage(uuid string) (*app.Message, error) {
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

func (db *Memory) CreateAttachment(att *app.Attachment) error {
	db.attMu.Lock()
	atts, ok := db.att[att.MessageUUID]
	if !ok {
		atts = []app.Attachment{}
	}

	atts = append(atts, *att)
	db.att[att.MessageUUID] = atts
	db.attMu.Unlock()

	return nil
}

func (db *Memory) LoadAttachment(msg *app.Message) error {
	atts, err := db.GetAttachmentsByMessage(msg)
	if err != nil {
		return err
	}

	msg.Attachments = atts

	return nil
}

func (db *Memory) GetAttachmentsByMessage(message *app.Message) ([]app.Attachment, error) {
	db.attMu.RLock()
	atts, ok := db.att[message.UUID]
	if !ok {
		atts = []app.Attachment{}
	}
	db.attMu.RUnlock()

	return atts, nil
}
