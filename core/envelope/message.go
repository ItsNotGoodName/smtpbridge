package envelope

import "time"

type (
	Message struct {
		ID        int64
		To        map[string]struct{}
		From      string
		Subject   string
		Text      string
		HTML      string
		CreatedAt time.Time
	}
)

func NewMessage(From string, To []string, subject, text, html string) *Message {
	ToMap := make(map[string]struct{})
	for _, to := range To {
		ToMap[to] = struct{}{}
	}

	return &Message{
		To:        ToMap,
		From:      From,
		Subject:   subject,
		Text:      text,
		HTML:      html,
		CreatedAt: time.Now(),
	}
}
