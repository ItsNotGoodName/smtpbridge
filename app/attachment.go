package app

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Attachment struct {
	UUID string
	Name string
	Type AttachmentType
	Data []byte
}

type AttachmentType uint

const (
	TypePNG AttachmentType = iota
	TypeJPEG
)

func NewAttachment(name string, data []byte) (*Attachment, error) {
	var t AttachmentType
	contentType := http.DetectContentType(data)
	if contentType == "image/png" {
		t = TypePNG
	} else if contentType == "image/jpeg" {
		t = TypeJPEG
	} else {
		return nil, fmt.Errorf("invalid content type %s: %v", contentType, ErrInvalidAttachment)
	}

	return &Attachment{
		UUID: uuid.New().String(),
		Name: name,
		Type: t,
		Data: data,
	}, nil
}
