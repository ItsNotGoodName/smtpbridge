package envelope

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/samber/lo"
)

type CreateMessage struct {
	Date    time.Time
	Subject string
	From    string
	To      []string
	Text    string
	HTML    string
}

func NewMessage(r CreateMessage) *Message {
	return &Message{
		From:      r.From,
		To:        lo.Uniq(r.To),
		CreatedAt: time.Now(),
		Subject:   r.Subject,
		Text:      r.Text,
		HTML:      r.HTML,
		Date:      r.Date,
	}
}

type MessageFilter struct {
	Ascending     bool
	Search        string
	SearchSubject bool
	SearchText    bool
}

type MessageListResult struct {
	Messages   []*Message
	PageResult pagination.PageResult
}

type AttachmentFilter struct {
	Ascending bool
}

type AttachmentListResult struct {
	Attachments []*Attachment
	PageResult  pagination.PageResult
}
