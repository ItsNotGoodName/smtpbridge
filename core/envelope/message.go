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
		Date      string
		CreatedAt time.Time
	}
)

func NewMessage(from string, to []string, subject, text, html string, date string) *Message {
	toMap := make(map[string]struct{})
	for _, s := range to {
		toMap[s] = struct{}{}
	}

	return &Message{
		To:        toMap,
		From:      from,
		Subject:   subject,
		Text:      text,
		HTML:      html,
		Date:      date,
		CreatedAt: time.Now(),
	}
}
