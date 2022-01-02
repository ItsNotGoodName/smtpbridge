package domain

import "fmt"

var (
	ErrEndpointSendFailed    = fmt.Errorf("endpoint send failed")
	ErrEndpointInvalidType   = fmt.Errorf("invalid endpoint type")
	ErrEndpointInvalidConfig = fmt.Errorf("invalid endpoint config")
	ErrEndpointNotFound      = fmt.Errorf("endpoint not found")
	ErrEndpointNameConflict  = fmt.Errorf("endpoint name conflit")
)

type (
	// EndpointPort handles sending messages to an endpoint.
	EndpointPort interface {
		// Send sends the message to the endpoint.
		Send(msg *EndpointMessage) error
	}

	// MessageServicePort handles creating and sending messages.
	EndpointServicePort interface {
		SendBridges(msg *Message, bridges []Bridge) (Status, error)
	}

	EndpointRepositoryPort interface {
		Create(name, endpointType string, config map[string]string) error
		Get(name string) (EndpointPort, error)
	}
)
