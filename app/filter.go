package app

type Filter struct {
	To        string `json:"to,omitempty" yaml:"to,omitempty"`
	ToRegex   bool   `json:"to_regex,omitempty" yaml:"to_regex,omitempty"`
	From      string `json:"from,omitempty" yaml:"from,omitempty"`
	FromRegex bool   `json:"from_regex,omitempty" yaml:"from_regex,omitempty"`
}

func (f *Filter) Match(msg *Message) bool {
	// TODO: regex
	if msg.From == f.From {
		return true
	}
	if msg.To[f.To] {
		return true
	}
	return false
}
