package core

import "fmt"

var (
	ErrAuthInvalid         = fmt.Errorf("auth invalid")
	ErrDataExists          = fmt.Errorf("data exists")
	ErrDataNotFound        = fmt.Errorf("data not found")
	ErrDataTooBig          = fmt.Errorf("data too big")
	ErrEndpointNameExists  = fmt.Errorf("endpoint name exists")
	ErrEndpointNotFound    = fmt.Errorf("endpoint not found")
	ErrEndpointTypeInvalid = fmt.Errorf("endpoint type invalid")
	ErrMessageNotFound     = fmt.Errorf("message not found")
)
