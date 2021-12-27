package database

import (
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct {
	attMu sync.RWMutex
	att   map[string][]app.DataAttachment
	msgMu sync.RWMutex
	msg   map[string]app.Message
}

func NewMock() *Mock {
	return &Mock{
		attMu: sync.RWMutex{},
		att:   make(map[string][]app.DataAttachment),
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

func (db *Mock) GetMessage(uuid string) (*app.Message, error) {
	db.msgMu.RLock()
	msg := db.msg[uuid]
	db.msgMu.RUnlock()
	return &msg, nil
}

func (db *Mock) CreateAttachment(attachment *app.Attachment, data []byte) error {
	db.attMu.Lock()
	atts, ok := db.att[attachment.MessageUUID]
	if !ok {
		atts = []app.DataAttachment{}
	}

	atts = append(atts, *app.NewDataAttachment(attachment, data))
	db.att[attachment.MessageUUID] = atts
	db.attMu.Unlock()

	return nil
}
func (db *Mock) GetDataAttachmentsByMessage(message *app.Message) ([]app.DataAttachment, error) {
	db.attMu.RLock()
	atts, ok := db.att[message.UUID]
	if !ok {
		atts = []app.DataAttachment{}
	}
	db.attMu.RUnlock()

	return atts, nil
}
