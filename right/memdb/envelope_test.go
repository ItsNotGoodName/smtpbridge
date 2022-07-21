package memdb

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/stretchr/testify/assert"
)

const envelopeLimit int64 = 100

func TestEnvelopeCreateDelete(t *testing.T) {
	store := NewEnvelope(envelopeLimit)
	ctx := context.Background()
	msg, att := envelope.NewMessage("", []string{}, "", "", "", ""), []envelope.Attachment{*envelope.NewAttachment("", []byte{})}

	id, err := store.CreateEnvelope(ctx, msg, att)
	assert.Nil(t, err)

	env, err := store.GetEnvelope(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, store.lastAttachmentID, int64(1))
	assert.Equal(t, store.lastMessageID, int64(1))
	assert.Equal(t, env.Message.ID, int64(1))
	assert.Equal(t, env.Attachments[0].ID, int64(1))
	assert.Equal(t, env.Attachments[0].MessageID, int64(1))

	deleteErr := fmt.Errorf("delete error")
	err = store.DeleteEnvelope(ctx, id, func(env *envelope.Envelope) error { return deleteErr })
	assert.Equal(t, err, deleteErr)

	err = store.DeleteEnvelope(ctx, id, func(env *envelope.Envelope) error { return nil })
	assert.Nil(t, err)

	err = store.DeleteEnvelope(ctx, id, func(env *envelope.Envelope) error { return nil })
	assert.Equal(t, err, core.ErrMessageNotFound)

	_, err = store.GetEnvelope(ctx, id)
	assert.Equal(t, err, core.ErrMessageNotFound)
}

func TestEnvelopeCreateNoAttachments(t *testing.T) {
	store := NewEnvelope(envelopeLimit)
	ctx := context.Background()
	msg, att := envelope.NewMessage("", []string{}, "", "", "", ""), []envelope.Attachment{}

	id, err := store.CreateEnvelope(ctx, msg, att)
	assert.Nil(t, err)

	env, err := store.GetEnvelope(ctx, id)
	assert.Nil(t, err)

	assert.Equal(t, store.lastAttachmentID, int64(0))
	assert.Equal(t, store.lastMessageID, int64(1))
	assert.Equal(t, env.Message.ID, int64(1))
	assert.Len(t, env.Attachments, 0)
}

func TestEnvelopeListCount(t *testing.T) {
	store := NewEnvelope(envelopeLimit)
	ctx := context.Background()

	for i := 0; i < 12; i++ {
		msg, att := envelope.NewMessage(strconv.Itoa(i), []string{}, "", "", "", ""), []envelope.Attachment{*envelope.NewAttachment(strconv.Itoa(i), []byte{})}
		store.CreateEnvelope(ctx, msg, att)
	}

	envs, total, err := store.ListEnvelope(ctx, 0, 10, true)
	assert.Nil(t, err)
	assert.Equal(t, total, 12)
	assert.Len(t, envs, 10)
	assert.Equal(t, envs[0].Message.From, "0")
	assert.Equal(t, envs[9].Message.From, "9")

	count, err := store.CountEnvelope(ctx)
	assert.Nil(t, err)
	assert.Equal(t, count, 12)

	envs, total, err = store.ListEnvelope(ctx, 10, 2, true)
	assert.Nil(t, err)
	assert.Equal(t, total, 12)
	assert.Len(t, envs, 2)
	assert.Equal(t, envs[0].Message.From, "10")
	assert.Equal(t, envs[1].Message.From, "11")

	envs, total, err = store.ListEnvelope(ctx, 0, 2, false)
	assert.Nil(t, err)
	assert.Equal(t, total, 12)
	assert.Len(t, envs, 2)
	assert.Equal(t, envs[0].Message.From, "11")
	assert.Equal(t, envs[1].Message.From, "10")
}

func TestEnvelopeLimitCount(t *testing.T) {
	store := NewEnvelope(10)
	ctx := context.Background()

	msg, att := envelope.NewMessage("", []string{}, "", "", "", ""), []envelope.Attachment{}
	for i := 0; i < 12; i++ {
		store.CreateEnvelope(ctx, msg, att)
	}

	count, err := store.CountEnvelope(ctx)
	assert.Nil(t, err)
	assert.Equal(t, count, 10)
}
