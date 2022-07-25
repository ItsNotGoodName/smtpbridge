package core

import "fmt"

var (
	ErrMessageNotFound     = fmt.Errorf("message not found")
	ErrDataNotFound        = fmt.Errorf("data not found")
	ErrDataExists          = fmt.Errorf("data exists")
	ErrAuthInvalid         = fmt.Errorf("auth invalid")
	ErrEndpointTypeInvalid = fmt.Errorf("endpoint type invalid")
	ErrEndpointNotFound    = fmt.Errorf("endpoint not found")
	ErrEndpointNameExists  = fmt.Errorf("endpoint name exists")
)
