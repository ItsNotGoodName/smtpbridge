package bus

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/google/uuid"
	"github.com/mustafaturan/bus/v3"
	"github.com/rs/zerolog/log"
)

func logEmitErr(err error) {
	if err != nil {
		log.Err(err).Msg("Failed to emit bus event")
	}
}

type generator struct{}

func (generator) Generate() string {
	return uuid.NewString()
}

type Bus struct {
	bus *bus.Bus
}

func New() (Bus, error) {
	bus, err := bus.NewBus(generator{})
	if err != nil {
		return Bus{}, err
	}

	bus.RegisterTopics(
		TopicEnvelopeCreated,
		TopicEnvelopeDeleted,
		TopicMailmanEnqueued,
	)

	return Bus{
		bus: bus,
	}, nil
}

const (
	TopicEnvelopeCreated = "envelope.created"
	TopicEnvelopeDeleted = "envelope.deleted"
	TopicMailmanEnqueued = "mailman.enqueued"
)

// EnvelopeCreated implements core.Bus.
func (b Bus) EnvelopeCreated(ctx context.Context, id int64) {
	logEmitErr(b.bus.Emit(ctx, TopicEnvelopeCreated, id))
}

// OnEnvelopeCreated implements core.Bus.
func (b Bus) OnEnvelopeCreated(h func(ctx context.Context, evt models.EventEnvelopeCreated) error) func() {
	key := uuid.NewString()

	b.bus.RegisterHandler(key, bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			id := e.Data.(int64)
			h(ctx, models.EventEnvelopeCreated{
				ID: id,
			})
		},
		Matcher: TopicEnvelopeCreated,
	})

	return func() { b.bus.DeregisterHandler(key) }
}

// EnvelopeDeleted implements core.Bus.
func (b Bus) EnvelopeDeleted(ctx context.Context) {
	logEmitErr(b.bus.Emit(ctx, TopicEnvelopeDeleted, nil))
}

// OnEnvelopeDeleted implements core.Bus.
func (b Bus) OnEnvelopeDeleted(h func(ctx context.Context, evt models.EventEnvelopeDeleted) error) func() {
	key := uuid.NewString()

	b.bus.RegisterHandler(key, bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			h(ctx, models.EventEnvelopeDeleted{})
		},
		Matcher: TopicEnvelopeDeleted,
	})

	return func() { b.bus.DeregisterHandler(key) }
}

// MailmanEnqueued implements core.Bus.
func (b Bus) MailmanEnqueued(ctx context.Context) {
	logEmitErr(b.bus.Emit(ctx, TopicMailmanEnqueued, nil))
}

// OnMailmanEnqueued implements core.Bus.
func (b Bus) OnMailmanEnqueued(h func(ctx context.Context, evt models.EventMailmanEnqueued) error) func() {
	key := uuid.NewString()

	b.bus.RegisterHandler(key, bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			h(ctx, models.EventMailmanEnqueued{})
		},
		Matcher: TopicMailmanEnqueued,
	})

	return func() { b.bus.DeregisterHandler(key) }
}

var _ core.Bus = Bus{}
