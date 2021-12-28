package database

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct {
	attMu sync.RWMutex
	att   map[string][]app.Attachment
	msgMu sync.RWMutex
	msg   map[string]app.Message
}

func NewMock() *Mock {
	return &Mock{
		attMu: sync.RWMutex{},
		att:   make(map[string][]app.Attachment),
		msgMu: sync.RWMutex{},
		msg:   make(map[string]app.Message),
	}
}

func (db *Mock) CreateMessage(msg *app.Message) error {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()
	if _, ok := db.msg[msg.UUID]; ok {
		return app.ErrMessageAlreadyExists
	}

	db.msg[msg.UUID] = *msg

	return nil
}

func (db *Mock) UpdateMessage(uuid string, updateFN func(msg *app.Message) (*app.Message, error)) (*app.Message, error) {
	db.msgMu.Lock()
	defer db.msgMu.Unlock()

	msg, ok := db.msg[uuid]
	if !ok {
		return nil, app.ErrMessageNotFound
	}

	updatedMsg, err := updateFN(&msg)
	if err != nil {
		return nil, err
	}

	db.msg[uuid] = *updatedMsg

	return updatedMsg, nil
}

func (db *Mock) GetAttachment(uuid string) (*app.Attachment, error) {
	return nil, app.ErrNotImplemented
}

func (db *Mock) GetMessage(uuid string) (*app.Message, error) {
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

func (db *Mock) CreateAttachment(att *app.Attachment) error {
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

func (db *Mock) LoadAttachment(msg *app.Message) error {
	atts, err := db.GetAttachmentsByMessage(msg)
	if err != nil {
		return err
	}

	msg.Attachments = atts

	return nil
}

func (db *Mock) GetAttachmentsByMessage(message *app.Message) ([]app.Attachment, error) {
	db.attMu.RLock()
	atts, ok := db.att[message.UUID]
	if !ok {
		atts = []app.Attachment{}
	}
	db.attMu.RUnlock()

	return atts, nil
}
