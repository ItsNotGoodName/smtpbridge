package envelope

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/samber/lo"
)

func NewMessage(from string, to []string, subject, text, html string, date time.Time) *Message {
	return &Message{
		From:      from,
		To:        lo.Uniq(to),
		CreatedAt: time.Now(),
		Subject:   subject,
		Text:      text,
		HTML:      html,
		Date:      date,
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
