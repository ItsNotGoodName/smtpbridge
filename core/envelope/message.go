package envelope

import "time"

type (
	Message struct {
		ID        int64
		CreatedAt time.Time
		Date      time.Time
		Subject   string
		From      string
		To        map[string]struct{}
		Text      string
		HTML      string
	}
)

func NewMessage(from string, to []string, subject, text, html string, date time.Time) *Message {
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
