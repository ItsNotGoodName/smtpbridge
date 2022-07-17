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
	// Get envelopes
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

	// Sort envelopes
	if ascending {
		sort.Slice(allEnvs, func(i, j int) bool {
			return allEnvs[i].Message.ID < allEnvs[j].Message.ID
		})
	} else {
		sort.Slice(allEnvs, func(i, j int) bool {
			return allEnvs[i].Message.ID > allEnvs[j].Message.ID
		})
	}

	// Slice envelopes
	envs := []envelope.Envelope{}
	end := offset + limit
	for i := offset; i < length && i < end; i++ {
		envs = append(envs, allEnvs[i])
	}

	return envs, length, nil
}

func (e *Envelope) CreateEnvelope(ctx context.Context, msg *envelope.Message, atts []envelope.Attachment) error {
	// Create IDs
	msg.ID = atomic.AddInt64(&e.lastMessageID, 1)
	for i := range atts {
		atts[i].ID = atomic.AddInt64(&e.lastAttachmentID, 1)
	}

	// Create envelope
	e.mu.Lock()
	e.messages[msg.ID] = *msg
	copy(e.attachments[msg.ID], atts)
	count := len(e.messages)
	e.mu.Unlock()

	// Delete oldest envelope if full
	if count > maxMessages {
		e.DeleteEnvelopeAndData(ctx, msg.ID-maxMessages)
	}

	return nil
}

func (e *Envelope) DeleteEnvelopeAndData(ctx context.Context, msgID int64) error {
	// Delete envelope
	e.mu.Lock()
	atts := e.attachments[msgID]
	e.deleteEnvelope(msgID)
	e.mu.Unlock()

	// Delete attachment data for envelopes
	for _, att := range atts {
		if err := e.dataStore.DeleteData(ctx, &att); err != nil {
			return err
		}
	}

	return nil
}

func (e *Envelope) GetEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	e.mu.Lock()
	env, err := e.getEnvelope(msgID)
	e.mu.Unlock()

	return env, err
}

func (e *Envelope) GetAndDeleteEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	e.mu.Lock()
	env, err := e.getEnvelope(msgID)
	if err != nil {
		e.mu.Unlock()
		return nil, err
	}

	e.deleteEnvelope(msgID)
	e.mu.Unlock()

	return env, nil
}

func (e *Envelope) getEnvelope(msgID int64) (*envelope.Envelope, error) {
	msg, ok := e.messages[msgID]
	if !ok {
		return nil, core.ErrMessageNotFound
	}

	atts, ok := e.attachments[msgID]
	if !ok {
		return nil, core.ErrMessageNotFound
	}

	return &envelope.Envelope{Message: msg, Attachments: atts}, nil
}

func (e *Envelope) deleteEnvelope(msgID int64) {
	delete(e.messages, msgID)
	delete(e.attachments, msgID)
}
