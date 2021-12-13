package app

import (
	"time"

	"github.com/google/uuid"
)

type Status uint8

const (
	StatusReceived      Status = iota // Message was received.
	StatusNoMatch                     // Message does not have corresponding Endpoint.
	StatusNotSent                     // Message was not sent.
	StatusPartiallySent               // Message was sent to atleast one Endpoint, but not all.
	StatusSent                        // Message was sent.
)

type Message struct {
	Time    time.Time       `json:"time"`    // Time message was received.
	UUID    string          `json:"uuid"`    // UUID of the message.
	Saved   bool            `json:"saved"`   // Saved represents whether the message has been saved.
	Status  Status          `json:"status"`  // Status of message.
	Subject string          `json:"subject"` // Subject of the message.
	From    string          `json:"from"`    // From is the email address of the sender.
	To      map[string]bool `json:"to"`      // To is the email addresses of the recipients.
	Text    string          `json:"text"`    // Text is the message body.
	//	Attachments []Attachment    `json:"attachments"` // Attachments is the list of attachments.
}

func NewMessage(subject, from string, to map[string]bool, text string) *Message {
	return &Message{
		Time:    time.Now(),
		UUID:    uuid.New().String(),
		Saved:   false,
		Status:  StatusReceived,
		Subject: subject,
		From:    from,
		To:      to,
		Text:    text,
	}
}

//const (
//	FileTypePNG uint = iota
//	FileTypeJPEG
//)
//
//type FileType uint
//
//type Attachment struct {
//	File string   `json:"file"` // File is the file to attach.
//	Type FileType `json:"type"` // Type is the type of the file.
//}
