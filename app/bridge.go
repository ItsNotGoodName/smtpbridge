package app

type Bridge struct {
	Name      string   `json:"name" mapstructure:"name"`
	Endpoints []string `json:"endpoints" mapstructure:"endpoints"`
	EmailFrom string   `json:"email_from" mapstructure:"email_from"`
	EmailTo   string   `json:"email_to" mapstructure:"email_to"`
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
