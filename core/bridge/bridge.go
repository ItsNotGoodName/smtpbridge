package bridge

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type (
	Bridge struct {
		Name      string     // Name is the name of the bridge.
		Endpoints []Endpoint // Endpoints are the endpoint.
		Filters   []Filter   // Filters are the filters.
	}

	Endpoint struct {
		NoText        bool
		NoAttachments bool
		Facade        *endpoint.Facade
	}

	Service interface {
		// ListByMessage returns bridges that the message belongs to.
		ListByMessage(msg *message.Message) []*Bridge
		// HandleMessage handles a message.
		HandleMessage(ctx context.Context, bridges []*Bridge, msg *message.Message, atts []attachment.Attachment) error
	}
)

func New(name string, facades []Endpoint, filters []Filter) *Bridge {
	return &Bridge{
		Name:      name,
		Endpoints: facades,
		Filters:   filters,
	}
}

func NewEndpoint(facade *endpoint.Facade, noText, noAttachments bool) Endpoint {
	return Endpoint{
		NoText:        noText,
		NoAttachments: noAttachments,
		Facade:        facade,
	}
}

func (b *Bridge) Match(msg *message.Message) bool {
	if len(b.Filters) == 0 {
		return true
	}
	for _, f := range b.Filters {
		if f.Match(msg) {
			return true
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
