package app

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Attachment struct {
	UUID        string         `json:"uuid" storm:"id"`
	Name        string         `json:"name"`
	Type        AttachmentType `json:"type"`
	MessageUUID string         `json:"message_uuid" storm:"index"`
	Data        []byte         `json:"-"`
}

type AttachmentType string

const (
	TypePNG  AttachmentType = "png"
	TypeJPEG AttachmentType = "jpeg"
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

	att := Attachment{
		UUID:        uuid.New().String(),
		Name:        name,
		Type:        t,
		MessageUUID: msg.UUID,
		Data:        data,
	}

	msg.Attachments = append(msg.Attachments, att)

	return &att, nil
}

type EndpointAttachment struct {
	Name string
	Type AttachmentType
	Data []byte
}

func NewEndpointAttachments(atts []Attachment) []EndpointAttachment {
	eats := make([]EndpointAttachment, len(atts))
	for i, a := range atts {
		eats[i] = EndpointAttachment{
			Name: a.Name,
			Type: a.Type,
			Data: a.Data,
		}
	}

	return eats
}
