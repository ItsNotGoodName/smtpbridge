package domain

import "fmt"

var ErrBridgesNotFound = fmt.Errorf("bridges not found")

type (
	Bridge struct {
		Name            string   `json:"name" mapstructure:"name"`
		Endpoints       []string `json:"endpoints" mapstructure:"endpoints"`
		OnlyText        bool     `json:"only_text" mapstructure:"only_text"`
		OnlyAttachments bool     `json:"only_attachments" mapstructure:"only_attachments"`
		Filters         []Filter `json:"filters" mapstructure:"filters"`
	}

	// BridgeServicePort handles finding endpoints for messages.
	BridgeServicePort interface {
		// GetBridges returns a list of bridges that the message belongs to.
		ListByMessage(msg *Message) []Bridge
	}
)

func (b *Bridge) Match(msg *Message) bool {
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

func (b *Bridge) EndpointMessage(msg *Message) *EndpointMessage {
	if b.OnlyText && !b.OnlyAttachments {
		return &EndpointMessage{Text: msg.Text}
	}
	if b.OnlyAttachments && !b.OnlyText {
		return &EndpointMessage{Attachments: NewEndpointAttachments(msg.Attachments)}
	}
	return &EndpointMessage{Text: msg.Text, Attachments: NewEndpointAttachments(msg.Attachments)}
}
