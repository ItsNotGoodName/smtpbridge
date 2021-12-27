package app

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Attachment struct {
	UUID        string
	Name        string
	Type        AttachmentType
	MessageUUID string
}

type AttachmentType uint

const (
	TypePNG AttachmentType = iota
	TypeJPEG
)

func NewAttachment(msg *Message, name string, data []byte) (*Attachment, error) {
	var t AttachmentType
	contentType := http.DetectContentType(data)
	if contentType == "image/png" {
		t = TypePNG
	} else if contentType == "image/jpeg" {
		t = TypeJPEG
	} else {
		return nil, fmt.Errorf("%s: %v", contentType, ErrAttachmentInvalid)
	}

	return &Attachment{
		UUID:        uuid.New().String(),
		Name:        name,
		Type:        t,
		MessageUUID: msg.UUID,
	}, nil
}

type DataAttachment struct {
	Name string
	Type AttachmentType
	Data []byte
}

func NewDataAttachment(attachment *Attachment, data []byte) *DataAttachment {
	return &DataAttachment{
		Name: attachment.Name,
		Type: attachment.Type,
		Data: data,
	}
}
