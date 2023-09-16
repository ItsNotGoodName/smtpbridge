package mailman

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rule"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/rs/zerolog/log"
)

type Mailman struct {
	id              int
	bus             core.Bus
	app             core.App
	fileStore       endpoint.FileStore
	endpointFactory endpoint.Factory
}

func New(id int, app core.App, bus core.Bus, fileStore endpoint.FileStore, endpointFactory endpoint.Factory) Mailman {
	return Mailman{
		id:              id,
		app:             app,
		bus:             bus,
		fileStore:       fileStore,
		endpointFactory: endpointFactory,
	}
}

func (m Mailman) Serve(ctx context.Context) error {
	checkC := make(chan struct{}, 1)
	checkC <- struct{}{}

	release := m.bus.OnMailmanEnqueued(func(ctx context.Context, evt models.EventMailmanEnqueued) error {
		select {
		case checkC <- struct{}{}:
		default:
		}

		return nil
	})
	defer release()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-checkC:
			tracer := m.app.
				Tracer(trace.SourceMailman).
				Sticky(trace.WithKV("mailman", m.id))

			tracer.Trace(ctx, "mailman.wake")

			for {
				maybeEnv, err := m.app.MailmanDequeue(ctx)
				if err != nil {
					tracer.Trace(ctx, "mailman.dequeue", trace.WithError(err))
					log.Err(err).Msg("Mailman failed to dequeue envelope")
					break
				}
				if maybeEnv == nil {
					break
				}
				env := *maybeEnv

				tracer := tracer.Sticky(trace.WithEnvelope(env.Message.ID))

				if err := m.send(ctx, tracer, env); err != nil {
					tracer.Trace(ctx, "mailman.error", trace.WithError(err))
					log.Err(err).Int64("envelope-id", env.Message.ID).Msg("Mailman failed to send envelope")
				}
			}

			tracer.Trace(ctx, "mailman.sleep")
		}
	}
}

func (m Mailman) send(ctx context.Context, tracer trace.Tracer, env models.Envelope) error {
	// List all rules
	rules, err := m.app.RuleEndpointsList(ctx)
	if err != nil {
		return err
	}

	if len(rules) == 0 {
		tracer.Trace(ctx, "mailman.rules.skip(empty)")
		return nil
	}

	sent := make(map[int64]struct{})
	for _, r := range rules {
		tracer := tracer.Sticky(trace.WithRule(r.Rule.ID))

		if len(r.Endpoints) == 0 {
			tracer.Trace(ctx, "mailman.rule.endpoints.skip(empty)")
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

		tracer.Trace(ctx, "mailman.rule.match.success")

		for _, e := range r.Endpoints {
			tracer := tracer.Sticky(trace.WithEndpoint(e.ID))

			// Prevent duplicate envelopes
			if _, ok := sent[e.ID]; ok {
				tracer.Trace(ctx, "mailman.rule.endpoint.skip(duplicate)")
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
				tracer.Trace(ctx, "mailman.rule.endpoint.send.success", trace.WithDuration(time.Now().Sub(start)))
			}
		}
	}

	return nil
}
