package core

import "fmt"

var (
	ErrMessageNotFound    = fmt.Errorf("message not found")
	ErrAttachmentNotFound = fmt.Errorf("attachment not found")
	ErrDataNotFound       = fmt.Errorf("data not found")
	ErrAuthInvalid        = fmt.Errorf("auth invalid")
)
