package envelope

import (
	"net/http"
)

const (
	DataTypeUnknown DataType = ""
	DataTypePNG     DataType = "png"
	DataTypeJPEG    DataType = "jpeg"
)

type (
	Attachment struct {
		ID        int64    `json:"id"`
		MessageID int64    `json:"message_id"`
		Name      string   `json:"name"`
		Type      DataType `json:"type"`
	}

	DataType string
)

func NewAttachment(messageID int64, name string, data []byte) *Attachment {
	return &Attachment{
		MessageID: messageID,
		Name:      name,
		Type:      ParseDataType(data),
	}
}

func ParseDataType(data []byte) DataType {
	if len(data) == 0 {
		return DataTypeUnknown
	}

	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/png":
		return DataTypePNG
	case "image/jpeg":
		return DataTypeJPEG
	default:
		return DataTypeUnknown
	}
}
