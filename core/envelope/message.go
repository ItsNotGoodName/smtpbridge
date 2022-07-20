package envelope

import "time"

type (
	Message struct {
		ID        int64               `json:"id"`
		To        map[string]struct{} `json:"to"`
		From      string              `json:"from"`
		Subject   string              `json:"subject"`
		Text      string              `json:"text"`
		HTML      string              `json:"html"`
		Date      string              `json:"date"`
		CreatedAt time.Time           `json:"created_at"`
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
