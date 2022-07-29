package envelope

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type (
	Attachment struct {
		ID        int64
		MessageID int64
		Name      string
		Mime      string
		Extension string
	}
)

func NewAttachment(name string, data []byte) *Attachment {
	mimeT := http.DetectContentType(data)

	extension := ""
	extensions, err := mime.ExtensionsByType(mimeT)
	if err == nil && extensions != nil {
		extension = extensions[0]
		// Use extension from name if it is valid
		unknownExt := filepath.Ext(name)
		for _, ext := range extensions {
			if ext == unknownExt {
				extension = ext
				break
			}
		}
	}

	return &Attachment{
		Name:      name,
		Mime:      mimeT,
		Extension: extension,
	}
}

func (a *Attachment) IsImage() bool {
	return strings.HasPrefix(a.Mime, "image/")
}

func (a *Attachment) FileName() string {
	return strconv.FormatInt(a.ID, 10) + a.Extension
}

func AttachmentIDFromFileName(fileName string) (int64, error) {
	return strconv.ParseInt(strings.Split(fileName, ".")[0], 10, 64)
}
