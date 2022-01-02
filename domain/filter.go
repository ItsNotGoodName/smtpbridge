package domain

type Filter struct {
	To   string
	From string
}

func NewFilter(to, from string) *Filter {
	return &Filter{
		To:   to,
		From: from,
	}
}

func (f *Filter) Match(msg *Message) bool {
	// TODO: regex
	if f.To != "" {
		if _, ok := msg.To[f.To]; !ok {
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
