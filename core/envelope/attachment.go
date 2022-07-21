package envelope

import (
	"mime"
	"net/http"
	"strconv"
	"strings"
)

type (
	Attachment struct {
		ID        int64
		MessageID int64
		Name      string
		Mime      string
	}
)

func NewAttachment(name string, data []byte) *Attachment {
	return &Attachment{
		Name: name,
		Mime: http.DetectContentType(data),
	}
}

func (a *Attachment) IsImage() bool {
	return strings.HasPrefix(a.Mime, "image/")
}

func (a *Attachment) FileName() string {
	extension := ""
	extensions, err := mime.ExtensionsByType(a.Mime)
	if err == nil && extensions != nil {
		extension = extensions[0]
	}

	return strconv.FormatInt(a.ID, 10) + extension
}

func AttachmentIDFromFileName(fileName string) (int64, error) {
	return strconv.ParseInt(strings.Split(fileName, ".")[0], 10, 64)
}
