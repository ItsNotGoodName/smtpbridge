package procs

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/events"
	"github.com/rs/zerolog/log"
)

func MailmanBackground(ctx context.Context, app core.App) {
	evtC := make(chan core.EventEnvelopeCreated, 25)

	go mailman(app.SystemContext(ctx), evtC)

	events.OnEnvelopeCreated(app, func(cc core.Context, evt core.EventEnvelopeCreated) {
		select {
		case evtC <- evt:
		default:
			log.Warn().Int64("id", evt.ID).Msg("Mailman is full")
		}
	})
}

func mailman(cc core.Context, evtC <-chan core.EventEnvelopeCreated) {
	ctx := cc

	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-evtC:
			// Get envelope
			env, err := db.EnvelopeGet(cc, evt.ID)
			if err != nil {
				log.Err(err).Msg("Failed to read envelope from database")
				continue
			}

			// List all rules
			rrules, err := db.RuleListEnable(cc)
			if err != nil {
				log.Err(err).Msg("Failed to list rules from database")
				continue
			}

			sent := make(map[int64]struct{})
			// Iterate over each rule
			for _, rule := range rrules {
				// Parse rule
				parsedRule, err := rule.Parse()
				if err != nil {
					log.Err(err).Msg("Failed to parse rule")
					continue
				}

				// Run rule
				match, err := parsedRule.Match(env)
				if err != nil {
					log.Err(err).Msg("Failed to run rule")
					continue
				}
				if !match {
					log.Info().Msgf("Failed rule %d", rule.ID)
					continue
				}
				log.Info().Msgf("Passed rule %d", rule.ID)

				// Get all endpoints for rule
				ends, err := db.EndpointListByRule(cc, rule.ID)
				if err != nil {
					log.Err(err).Msg("Failed to list endpoints from database")
					continue
				}

				// Iterator over endpoints
				for _, end := range ends {
					if _, ok := sent[end.ID]; ok {
						continue
					}
					sent[end.ID] = struct{}{}

					// Parse endpoint
					parsedEndpoint, err := end.Parse()
					if err != nil {
						log.Err(err).Msg("Failed to parse endpoint")
						continue
					}

					// Send envelope to endpoint
					if err := parsedEndpoint.Send(cc, env); err != nil {
						log.Err(err).Int64("envelope-id", env.Message.ID).Int64("endpoint-id", end.ID).Msg("Failed to send envelope to endpoint")
					}
				}
			}
		}
	}
}
