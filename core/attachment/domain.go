package attachment

import (
	"fmt"
	"net/http"
)

const (
	TypePNG  Type = "png"
	TypeJPEG Type = "jpeg"
)

var (
	ErrInvalid       = fmt.Errorf("attachment invalid")
	ErrDataEmpty     = fmt.Errorf("attachment data is empty")
	ErrDataNotLoaded = fmt.Errorf("attachment data not loaded")
	ErrNotFound      = fmt.Errorf("attachment not found")
)

func New(param *Param) (*Attachment, error) {
	attType, err := DataType(param.Data)
	if err != nil {
		return nil, err
	}

	return &Attachment{
		ID:        param.ID,
		Name:      param.Name,
		Type:      attType,
		MessageID: param.Message.ID,
		data:      param.Data,
	}, nil
}

// DataType returns the type of the attachment data.
func DataType(data []byte) (Type, error) {
	if len(data) == 0 {
		return "", ErrDataEmpty
	}

	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/png":
		return TypePNG, nil
	case "image/jpeg":
		return TypeJPEG, nil
	default:
		return "", fmt.Errorf("%s: %v", contentType, ErrInvalid)
	}
}

//// File returns the full attachment file name.
//func (a *Attachment) File() string {
//	return fmt.Sprintf("%d.%s", a.ID, a.Type)
//}

// GetData returns the attachment data.
func (a *Attachment) GetData() ([]byte, error) {
	if a.data == nil {
		return nil, ErrDataNotLoaded
	}

	return a.data, nil
}

// GetData returns the attachment data.
func (a *Attachment) SetData(data []byte) error {
	_, err := DataType(data)
	if err != nil {
		return err
	}

	a.data = data
	return nil
}
