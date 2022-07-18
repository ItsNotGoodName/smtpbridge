package envelope

import "time"

type (
	Message struct {
		ID        int64     `json:"id"`
		To        []string  `json:"to"`
		From      string    `json:"from"`
		Subject   string    `json:"subject"`
		Text      string    `json:"text"`
		HTML      string    `json:"html"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func NewMessage(From string, To []string, subject, text, html string) *Message {
	return &Message{
		To:        To,
		From:      From,
		Subject:   subject,
		Text:      text,
		HTML:      html,
		CreatedAt: time.Now(),
	}
}
