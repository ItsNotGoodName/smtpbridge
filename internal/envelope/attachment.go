package envelope

import (
	"net/http"
)

type DataAttachment struct {
	Attachment *Attachment
	Data       []byte
}

func NewDataAttachment(name string, data []byte) DataAttachment {
	mimeT := http.DetectContentType(data)
	extension := fileExtension(name, mimeT)
	return DataAttachment{
		Attachment: &Attachment{
			Name:      name,
			Mime:      mimeT,
			Extension: extension,
		},
		Data: data,
	}
}
