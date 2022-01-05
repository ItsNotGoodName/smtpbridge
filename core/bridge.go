package core

import "fmt"

var ErrBridgesNotFound = fmt.Errorf("bridges not found")

type (
	Bridge struct {
		Name            string
		Endpoints       []string
		OnlyText        bool
		OnlyAttachments bool
		Filters         []Filter
	}

	BridgeServicePort interface {
		// ListByMessage returns bridges that the message belongs to.
		ListByMessage(msg *Message) []*Bridge
	}
)

func NewBridge(name string, endpoints []string, onlyText, onlyAttachments bool, filters []Filter) *Bridge {
	return &Bridge{
		Name:            name,
		Endpoints:       endpoints,
		OnlyText:        onlyText,
		OnlyAttachments: onlyAttachments,
		Filters:         filters,
	}
}

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

func (b *Bridge) EndpointMessage(msg *Message, atts []Attachment) *EndpointMessage {
	if b.OnlyText && !b.OnlyAttachments {
		return &EndpointMessage{Text: msg.Text}
	}
	if b.OnlyAttachments && !b.OnlyText {
		return &EndpointMessage{Attachments: NewEndpointAttachments(atts)}
	}
	return &EndpointMessage{Text: msg.Text, Attachments: NewEndpointAttachments(atts)}
}
