package trace

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Store interface {
	Save(ctx context.Context, trace models.Trace) error
}

type Option func(t *models.Trace)

type Tracer struct {
	store     Store
	source    string
	requestID string
	seq       *int32
	sticky    []Option
}

func NewTracer(store Store, source string) Tracer {
	requestID := uuid.NewString()

	var seq int32
	return Tracer{
		store:     store,
		source:    source,
		requestID: requestID,
		seq:       &seq,
	}
}

func (t Tracer) Sticky(opts ...Option) Tracer {
	stickyLen := len(t.sticky)
	sticky := make([]Option, stickyLen, stickyLen+len(opts))
	copy(sticky, t.sticky)
	sticky = append(sticky, opts...)
	t.sticky = sticky
	return t
}

func (t Tracer) Trace(ctx context.Context, action string, options ...Option) {
	trace := models.Trace{
		Seq:       int(atomic.AddInt32(t.seq, 1)),
		Source:    t.source,
		RequestID: t.requestID,
		Action:    action,
		Data:      []models.TraceDataKV{},
		CreatedAt: models.NewTime(time.Now()),
		Level:     LevelInfo,
	}

	for _, sticky := range t.sticky {
		sticky(&trace)
	}

	for _, option := range options {
		option(&trace)
	}

	if err := t.store.Save(ctx, trace); err != nil {
		log.Err(err).Msg("Trace save failed")
	}
}

const (
	SourceHTTP    = "http"
	SourceSMTP    = "smtp"
	SourceMailman = "mailman"
	SourceCron    = "cron"
	SourceApp     = "app"
)

const (
	LevelInfo  models.TraceLevel = "info"
	LevelError models.TraceLevel = "error"
)

const (
	ActionEnvelopeCreated = "envelope.created"
)

func WithKV(key string, value any) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		})
	}
}

const KeyError = "error"

func WithError(err error) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyError,
			Value: err.Error(),
		})
		t.Level = LevelError
	}
}

const KeyAddress = "address"

func WithAddress(address string) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyAddress,
			Value: address,
		})
	}
}

const KeyEnvelope = "envelope"

func WithEnvelope(id int64) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyEnvelope,
			Value: strconv.FormatInt(id, 10),
		})
	}
}

const KeyEndpoint = "endpoint"

func WithEndpoint(id int64) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyEndpoint,
			Value: strconv.FormatInt(id, 10),
		})
	}
}

const KeyRule = "rule"

func WithRule(id int64) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyRule,
			Value: strconv.FormatInt(id, 10),
		})
	}
}

const KeyAttachment = "attachment"

func WithAttachment(id int64) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyAttachment,
			Value: strconv.FormatInt(id, 10),
		})
	}
}

const KeyDuration = "duration"

func WithDuration(d time.Duration) Option {
	return func(t *models.Trace) {
		t.Data = append(t.Data, models.TraceDataKV{
			Key:   KeyDuration,
			Value: d.String(),
		})
	}
}
