package core

import "fmt"

var (
	ErrMessageNotFound    = fmt.Errorf("message not found")
	ErrAttachmentNotFound = fmt.Errorf("attachment not found")
	ErrDataNotFound       = fmt.Errorf("data not found")
	ErrDataExists         = fmt.Errorf("data exists")
	ErrAuthInvalid        = fmt.Errorf("auth invalid")
)
