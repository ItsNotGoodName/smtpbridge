package domain

import "fmt"

var (
	ErrEndpointSendFailed    = fmt.Errorf("endpoint send failed")
	ErrEndpointInvalidType   = fmt.Errorf("invalid endpoint type")
	ErrEndpointInvalidConfig = fmt.Errorf("invalid endpoint config")
	ErrEndpointNotFound      = fmt.Errorf("endpoint not found")
	ErrEndpointNameConflict  = fmt.Errorf("endpoint name conflict")
)

type (
	EndpointPort interface {
		// Send sends the message to the endpoint.
		Send(msg *EndpointMessage) error
	}

	EndpointServicePort interface {
		// SendBridges sends the message to the bridge's endpoints if they pass the filter.
		SendBridges(msg *Message, bridges []Bridge) (Status, error)
	}

	EndpointRepositoryPort interface {
		// Create initializes a new endpoint.
		Create(name, endpointType string, config map[string]string) error
		// Get returns an endpoint by it's name.
		Get(name string) (EndpointPort, error)
	}
)
