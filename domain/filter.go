package domain

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
