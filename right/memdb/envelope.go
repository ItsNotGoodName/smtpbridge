package memdb

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

const maxMessages = 30

type Envelope struct {
	lastMessageID    int64
	lastAttachmentID int64
	dataStore        envelope.DataStore

	mu          sync.Mutex
	messages    map[int64]envelope.Message
	attachments map[int64][]envelope.Attachment
}

func NewEnvelope(dataStore envelope.DataStore) *Envelope {
	return &Envelope{
		dataStore:   dataStore,
		messages:    make(map[int64]envelope.Message),
		attachments: map[int64][]envelope.Attachment{},
	}
}

func (e *Envelope) ListEnvelope(ctx context.Context, offset, limit int, ascending bool) ([]envelope.Envelope, int, error) {
	e.mu.Lock()
	length := len(e.messages)
	allEnvs := make([]envelope.Envelope, 0, length)
	for _, msg := range e.messages {
		allEnvs = append(allEnvs, envelope.Envelope{
			Message:     msg,
			Attachments: e.attachments[msg.ID],
		})
	}
	e.mu.Unlock()

	if ascending {
		sort.Slice(allEnvs, func(i, j int) bool {
			return allEnvs[i].Message.ID < allEnvs[j].Message.ID
		})
	} else {
		sort.Slice(allEnvs, func(i, j int) bool {
			return allEnvs[i].Message.ID > allEnvs[j].Message.ID
		})
	}

	envs := []envelope.Envelope{}
	end := offset + limit
	for i := offset; i < length && i < end; i++ {
		envs = append(envs, allEnvs[i])
	}

	return envs, length, nil
}

func (e *Envelope) CreateMessage(ctx context.Context, msg *envelope.Message) error {
	msg.ID = atomic.AddInt64(&e.lastMessageID, 1)

	e.mu.Lock()
	e.messages[msg.ID] = *msg
	count := len(e.messages)
	if count > maxMessages {
		e.mu.Unlock()
		e.cleanUpMessage(ctx, msg.ID-maxMessages)
	} else {
		e.mu.Unlock()
	}

	return nil
}

func (e *Envelope) cleanUpMessage(ctx context.Context, msgID int64) error {
	e.mu.Lock()
	atts := e.attachments[msgID]
	e.deleteMessage(msgID)

	for _, att := range atts {
		if err := e.dataStore.DeleteData(ctx, &att); err != nil {
			e.mu.Unlock()
			return err
		}
	}
	e.mu.Unlock()

	return nil
}

func (e *Envelope) deleteMessage(msgID int64) {
	delete(e.messages, msgID)
	delete(e.attachments, msgID)
}

func (e *Envelope) DeleteMessage(ctx context.Context, msgID int64) error {
	e.mu.Lock()
	e.deleteMessage(msgID)
	e.mu.Unlock()

	return nil
}

func (e *Envelope) GetMessage(ctx context.Context, id int64) (*envelope.Message, error) {
	e.mu.Lock()
	msg, ok := e.messages[id]
	e.mu.Unlock()

	if !ok {
		return nil, core.ErrMessageNotFound
	}

	return &msg, nil
}

func (e *Envelope) CreateAttachment(ctx context.Context, att *envelope.Attachment) error {
	att.ID = atomic.AddInt64(&e.lastAttachmentID, 1)

	e.mu.Lock()
	atts, ok := e.attachments[att.MessageID]
	if ok {
		atts = append(atts, *att)
	} else {
		atts = []envelope.Attachment{}
	}
	e.attachments[att.MessageID] = atts
	e.mu.Unlock()

	return nil
}

func (e *Envelope) GetAttachment(ctx context.Context, id int64) (*envelope.Attachment, error) {
	e.mu.Lock()
	for _, atts := range e.attachments {
		for _, att := range atts {
			if att.ID == id {
				e.mu.Unlock()
				return &att, nil
			}

		}
	}
	e.mu.Unlock()

	return nil, core.ErrAttachmentNotFound
}

func (e *Envelope) GetAttachmentsByMessageID(ctx context.Context, msgID int64) ([]envelope.Attachment, error) {
	e.mu.Lock()
	atts, ok := e.attachments[msgID]
	if !ok {
		e.mu.Unlock()
		return nil, core.ErrMessageNotFound
	}
	e.mu.Unlock()

	return atts, nil
}
