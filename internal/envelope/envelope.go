package envelope

import (
	"strconv"
	"strings"
	"time"
)

type Envelope struct {
	Message     *Message
	Attachments []*Attachment
}

type Message struct {
	ID        int64 `bun:"id,pk,autoincrement"`
	CreatedAt time.Time
	Date      time.Time
	Subject   string
	From      string   `bun:"from_"`
	To        []string `bun:"to_"`
	Text      string
	HTML      string
}

func (e Message) IsTo(to string) bool {
	for _, t := range e.To {
		if t == to {
			return true
		}
	}
	return false
}

type Attachment struct {
	ID        int64 `bun:"id,pk,autoincrement"`
	MessageID int64
	Name      string
	Mime      string
	Extension string
}

func (a *Attachment) IsImage() bool {
	return strings.HasPrefix(a.Mime, "image/")
}

func (a *Attachment) FileName() string {
	return strconv.FormatInt(a.ID, 10) + a.Extension
}
