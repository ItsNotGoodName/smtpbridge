package app

import "errors"

var ErrNotAuthorized = errors.New("not authorized")

type Message struct {
	Subject string          `json:"subject"` // Subject of the message.
	From    string          `json:"from"`    // From is the email address of the sender.
	To      map[string]bool `json:"to"`      // To is the email addresses of the recipients.
	Text    string          `json:"text"`    // Text is the message body.
	//	Attachments []Attachment    `json:"attachments"` // Attachments is the list of attachments.
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
