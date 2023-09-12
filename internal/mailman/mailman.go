package mailman

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rule"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/rs/zerolog/log"
)

type Mailman struct {
	bus             core.Bus
	app             core.App
	fileStore       endpoint.FileStore
	endpointFactory endpoint.Factory
	queueLimit      int
}

func New(app core.App, bus core.Bus, fileStore endpoint.FileStore, endpointFactory endpoint.Factory) Mailman {
	return Mailman{
		app:             app,
		bus:             bus,
		fileStore:       fileStore,
		endpointFactory: endpointFactory,
		queueLimit:      100,
	}
}

func (m Mailman) Serve(ctx context.Context) error {
	idC := make(chan int64, m.queueLimit)
	release := m.bus.OnEnvelopeCreated(func(ctx context.Context, evt models.EventEnvelopeCreated) error {
		select {
		case idC <- evt.ID:
		default:
			m.app.Tracer(trace.SourceMailman).Trace(ctx,
				"mailman.overflow",
				trace.WithEnvelope(evt.ID),
				trace.WithError(fmt.Errorf("mailman is full")))
		}

		return nil
	})
	defer release()

	for {
		select {
		case <-ctx.Done():
			return nil
		case id := <-idC:
			tracer := m.app.Tracer(trace.SourceMailman).
				Sticky(trace.WithEnvelope(id))

			tracer.Trace(ctx, "mailman.start")
			err := m.send(ctx, tracer, id)
			if err != nil {
				tracer.Trace(ctx, "mailman.error", trace.WithError(err))
				log.Err(err).Int64("envelope-id", id).Msg("Failed to send envelope")
			}
			tracer.Trace(ctx, "mailman.end")
		}
	}
}

func (m Mailman) send(ctx context.Context, tracer trace.Tracer, envelopeID int64) error {
	// Get envelope
	env, err := m.app.EnvelopeGet(ctx, envelopeID)
	if err != nil {
		return err
	}

	// List all rules
	rules, err := m.app.RuleEndpointsList(ctx)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		tracer.Trace(ctx, "mailman.rules.skip.empty")
		return nil
	}

	sent := make(map[int64]struct{})
	for _, r := range rules {
		tracer := tracer.Sticky(trace.WithRule(r.Rule.ID))

		if len(r.Endpoints) == 0 {
			tracer.Trace(ctx, "mailman.rule.endpoints.skip.empty")
			continue
		}

		// Build rule
		rule, err := rule.Build(r.Rule)
		if err != nil {
			tracer.Trace(ctx, "mailman.rule.build.error", trace.WithError(err))
			continue
		}

		// Match rules
		ok, err := rule.Match(env)
		if err != nil {
			tracer.Trace(ctx, "mailman.rule.match.error", trace.WithError(err))
			continue
		}
		if !ok {
			tracer.Trace(ctx, "mailman.rule.match.fail")
			continue
		}

		tracer.Trace(ctx, "mailman.rule.match.pass")

		for _, e := range r.Endpoints {
			tracer := tracer.Sticky(trace.WithEndpoint(e.ID))

			// Prevent duplicate envelopes
			if _, ok := sent[e.ID]; ok {
				tracer.Trace(ctx, "mailman.rule.endpoint.skip.duplicate")
				continue
			}
			sent[e.ID] = struct{}{}

			// Build endpoint
			end, err := m.endpointFactory.Build(e)
			if err != nil {
				tracer.Trace(ctx, "mailman.rule.endpoint.build.error", trace.WithError(err))
				continue
			}

			start := time.Now()

			// Send envelope to endpoint
			err = end.Send(ctx, m.fileStore, env)
			if err != nil {
				tracer.Trace(ctx, "mailman.rule.endpoint.send.error", trace.WithError(err), trace.WithDuration(time.Now().Sub(start)))
			} else {
				tracer.Trace(ctx, "mailman.rule.endpoint.send", trace.WithDuration(time.Now().Sub(start)))
			}
		}
	}

	return nil
}
