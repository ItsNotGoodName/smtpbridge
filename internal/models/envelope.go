package models

import (
	"strconv"
	"strings"
)

type Envelope struct {
	Message     Message
	Attachments []Attachment
}

type Message struct {
	ID        int64 `sql:"primary_key"`
	CreatedAt Time
	Date      Time
	Subject   string
	From      string
	To        MessageTo
	Text      string
	HTML      string
}

type MessageTo []string

func (ss MessageTo) EQ(strs ...string) bool {
	for _, str := range strs {
		for _, s := range ss {
			if s == str {
				return true
			}
		}
	}
	return false
}

type Attachment struct {
	ID        int64 `sql:"primary_key"`
	MessageID int64
	Name      string
	Mime      string
	Extension string
}

func (a Attachment) IsImage() bool {
	return strings.HasPrefix(a.Mime, "image/")
}

func (a Attachment) FileName() string {
	return strconv.FormatInt(a.ID, 10) + a.Extension
}
