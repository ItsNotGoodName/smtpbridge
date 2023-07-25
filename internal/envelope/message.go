package envelope

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
	"github.com/jaytaylor/html2text"
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
	text := r.Text
	if isHTML(r.Text) {
		var err error
		text, err = html2text.FromString(r.Text, html2text.Options{
			TextOnly: true,
		})
		if err != nil {
			text = r.Text
		}
	}

	return &Message{
		From:      r.From,
		To:        lo.Uniq(r.To),
		CreatedAt: time.Now(),
		Subject:   r.Subject,
		Text:      text,
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
