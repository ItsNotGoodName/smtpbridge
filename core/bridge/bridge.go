package bridge

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type (
	Bridge struct {
		Name           string     // Name is the name of the bridge.
		Filters        []Filter   // Filters are the filters.
		Endpoints      []Endpoint // Endpoints are the endpoint.
		MinAttachments int        // MinAttachments is the minimum number of attachments required.
	}

	Endpoint struct {
		NoText        bool
		NoAttachments bool
		Facade        *endpoint.Facade
	}

	Service interface {
		// ListByEnvelope returns bridges that the envelope belongs to.
		ListByEnvelope(env envelope.Envelope) []*Bridge
		// HandleEnvelope handles an envelope.
		HandleEnvelope(ctx context.Context, bridges []*Bridge, env envelope.Envelope) error
	}
)

func New(name string, facades []Endpoint, filters []Filter, minAttachments int) *Bridge {
	return &Bridge{
		Name:           name,
		Endpoints:      facades,
		Filters:        filters,
		MinAttachments: minAttachments,
	}
}

func NewEndpoint(facade *endpoint.Facade, noText, noAttachments bool) Endpoint {
	return Endpoint{
		NoText:        noText,
		NoAttachments: noAttachments,
		Facade:        facade,
	}
}

func (b *Bridge) Match(env envelope.Envelope) bool {
	if len(env.Attachments) >= b.MinAttachments {
		if len(b.Filters) == 0 {
			return true
		}

		for _, f := range b.Filters {
			if f.Match(env.Message) {
				return true
			}
		}
	}

	return false
}

func (e Endpoint) Envelope(msg *endpoint.Message, atts []endpoint.Attachment) endpoint.Envelope {
	if e.NoAttachments && !e.NoText {
		return endpoint.NewEnvelope(msg, []endpoint.Attachment{})
	}

	if e.NoText && !e.NoAttachments {
		return endpoint.NewEnvelope(&endpoint.Message{ID: msg.ID}, atts)
	}

	return endpoint.NewEnvelope(msg, atts)
}
