package domain

import "fmt"

var (
	ErrEndpointSendFailed    = fmt.Errorf("endpoint send failed")
	ErrEndpointInvalidType   = fmt.Errorf("invalid endpoint type")
	ErrEndpointInvalidConfig = fmt.Errorf("invalid endpoint config")
	ErrEndpointNotFound      = fmt.Errorf("endpoint not found")
	ErrEndpointNameConflict  = fmt.Errorf("endpoint name conflit")
	ErrBridgesNotFound       = fmt.Errorf("bridges not found")
	ErrMessageNotFound       = fmt.Errorf("message not found")
	ErrMessageAlreadyExists  = fmt.Errorf("message already exists")
	ErrAttachmentInvalid     = fmt.Errorf("invalid attachment")
	ErrNotImplemented        = fmt.Errorf("not implemented")
)
