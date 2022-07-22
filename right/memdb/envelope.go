package memdb

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Envelope struct {
	limit            int64
	lastMessageID    int64
	lastAttachmentID int64

	mu          sync.Mutex
	messages    map[int64]envelope.Message
	attachments map[int64][]envelope.Attachment
}

func NewEnvelope(limit int64) *Envelope {
	return &Envelope{
		limit:       limit,
		messages:    make(map[int64]envelope.Message),
		attachments: map[int64][]envelope.Attachment{},
	}
}

func (e *Envelope) CountEnvelope(ctx context.Context) (int, error) {
	e.mu.Lock()
	count := len(e.messages)
	e.mu.Unlock()

	return count, nil
}

func (e *Envelope) CountAttachment(ctx context.Context) (int, error) {
	e.mu.Lock()
	count := 0
	for _, atts := range e.attachments {
		count += len(atts)
	}
	e.mu.Unlock()

	return count, nil
}

func (e *Envelope) ListAttachment(ctx context.Context, offset, limit int, ascending bool) ([]envelope.Attachment, int, error) {
	// Get attachments
	e.mu.Lock()
	count := 0
	all := []envelope.Attachment{}
	for _, atts := range e.attachments {
		for _, att := range atts {
			all = append(all, att)
			count++
		}
	}
	e.mu.Unlock()

	// Sort attachments
	if ascending {
		sort.Slice(all, func(i, j int) bool {
			return all[i].ID < all[j].ID
		})
	} else {
		sort.Slice(all, func(i, j int) bool {
			return all[i].ID > all[j].ID
		})
	}

	// Slice attachments
	atts := []envelope.Attachment{}
	end := offset + limit
	for i := offset; i < count && i < end; i++ {
		atts = append(atts, all[i])
	}

	return atts, count, nil
}

func (e *Envelope) ListEnvelope(ctx context.Context, offset, limit int, ascending bool) ([]envelope.Envelope, int, error) {
	// Get envelopes
	e.mu.Lock()
	count := len(e.messages)
	all := make([]envelope.Envelope, 0, count)
	for _, msg := range e.messages {
		all = append(all, envelope.Envelope{
			Message:     msg,
			Attachments: e.attachments[msg.ID],
		})
	}
	e.mu.Unlock()

	// Sort envelopes
	if ascending {
		sort.Slice(all, func(i, j int) bool {
			return all[i].Message.ID < all[j].Message.ID
		})
	} else {
		sort.Slice(all, func(i, j int) bool {
			return all[i].Message.ID > all[j].Message.ID
		})
	}

	// Slice envelopes
	envs := []envelope.Envelope{}
	end := offset + limit
	for i := offset; i < count && i < end; i++ {
		envs = append(envs, all[i])
	}

	return envs, count, nil
}

func (e *Envelope) CreateEnvelope(ctx context.Context, msg *envelope.Message, atts []envelope.Attachment) (int64, error) {
	// Create IDs
	msg.ID = atomic.AddInt64(&e.lastMessageID, 1)
	for i := range atts {
		atts[i].ID = atomic.AddInt64(&e.lastAttachmentID, 1)
		atts[i].MessageID = msg.ID
	}

	e.mu.Lock()
	// Create envelope
	e.messages[msg.ID] = *msg
	e.attachments[msg.ID] = atts
	e.deleteEnvelope(msg.ID - e.limit)
	e.mu.Unlock()

	return msg.ID, nil
}

func (e *Envelope) GetEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	e.mu.Lock()
	env, err := e.getEnvelope(msgID)
	e.mu.Unlock()

	return env, err
}

func (e *Envelope) getEnvelope(msgID int64) (*envelope.Envelope, error) {
	msg, ok := e.messages[msgID]
	if !ok {
		return nil, core.ErrMessageNotFound
	}

	return &envelope.Envelope{Message: msg, Attachments: e.attachments[msgID]}, nil
}

func (e *Envelope) DeleteEnvelope(ctx context.Context, msgID int64, fn func(env *envelope.Envelope) error) error {
	e.mu.Lock()
	env, err := e.getEnvelope(msgID)
	if err != nil {
		e.mu.Unlock()
		return err
	}

	if err := fn(env); err != nil {
		e.mu.Unlock()
		return err
	}

	e.deleteEnvelope(msgID)
	e.mu.Unlock()

	return nil
}

func (e *Envelope) deleteEnvelope(msgID int64) {
	delete(e.messages, msgID)
	delete(e.attachments, msgID)
}
