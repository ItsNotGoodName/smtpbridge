package app

import "fmt"

var (
	ErrEndpointInvalidType   = fmt.Errorf("invalid endpoint type")
	ErrEndpointInvalidConfig = fmt.Errorf("invalid endpoint config")
	ErrEndpointNotFound      = fmt.Errorf("endpoint not found")
	ErrEndpointNameConflict  = fmt.Errorf("endpoint name conflit")
	ErrBridgesNotFound       = fmt.Errorf("bridges not found")
	ErrMessageNotFound       = fmt.Errorf("message not found")
	ErrAttachmentInvalid     = fmt.Errorf("invalid attachment")
)
