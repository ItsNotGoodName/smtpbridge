package app

type Bridge struct {
	Name            string   `json:"name" mapstructure:"name"`
	Endpoints       []string `json:"endpoints" mapstructure:"endpoints"`
	EmailFrom       string   `json:"email_from" mapstructure:"email_from"`
	EmailTo         string   `json:"email_to" mapstructure:"email_to"`
	OnlyText        bool     `json:"only_text" mapstructure:"only_text"`
	OnlyAttachments bool     `json:"only_attachments" mapstructure:"only_attachments"`
}

func (b *Bridge) Match(msg *Message) bool {
	if b.EmailTo != "" {
		if !msg.To[b.EmailTo] {
			return false
		}
	}
	if b.EmailFrom != "" {
		if msg.From != b.EmailFrom {
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
		return &EndpointMessage{Attachments: msg.Attachments}
	}
	return &EndpointMessage{Text: msg.Text, Attachments: msg.Attachments}
}
