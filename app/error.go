package app

import "fmt"

var (
	ErrInvalidEndpointType   = fmt.Errorf("invalid endpoint type")
	ErrInvalidEndpointConfig = fmt.Errorf("invalid endpoint config")
	ErrNoEndpoints           = fmt.Errorf("no endpoints")
	ErrNoBridges             = fmt.Errorf("no bridges")
	ErrMessageNotFound       = fmt.Errorf("message not found")
	ErrInvalidAttachment     = fmt.Errorf("invalid attachment")
)
