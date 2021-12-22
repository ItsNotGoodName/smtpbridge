package app

import "fmt"

var (
	ErrInvalidEndpointType   = fmt.Errorf("invalid endpoint type")
	ErrInvalidEndpointConfig = fmt.Errorf("invalid endpoint config")
	ErrNoEndpoints           = fmt.Errorf("no endpoints")
	ErrMessageNotFound       = fmt.Errorf("message not found")
)
