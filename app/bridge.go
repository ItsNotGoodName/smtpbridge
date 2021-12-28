package app

type Bridge struct {
	Name            string   `json:"name" mapstructure:"name"`
	Endpoints       []string `json:"endpoints" mapstructure:"endpoints"`
	OnlyText        bool     `json:"only_text" mapstructure:"only_text"`
	OnlyAttachments bool     `json:"only_attachments" mapstructure:"only_attachments"`
	Filters         []Filter `json:"filters" mapstructure:"filters"`
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

type Filter struct {
	To   string `json:"to,omitempty" mapstructure:"to,omitempty"`
	From string `json:"from,omitempty" mapstructure:"from,omitempty"`
}

func (f *Filter) Match(msg *Message) bool {
	// TODO: regex
	if f.To != "" {
		if !msg.To[f.To] {
			return false
		}
	}
	if f.From != "" {
		if msg.From != f.From {
			return false
		}
	}
	return true
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
