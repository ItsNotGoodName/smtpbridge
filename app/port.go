package app

type AuthServicePort interface {
	Login(username, password string) error
}

type BridgeServicePort interface {
	// GetEndpoints returns a list of endpoints that the message should be sent to.
	GetEndpoints(msg *Message) ([]EndpointPort, error)
}

type MessageServicePort interface {
	// Handle receives message from drivers.
	Handle(msg *Message) error
}

type MessageRepositoryPort interface {
	Create(msg *Message) error
	Update(msg *Message) error
}

type EndpointPort interface {
	// Capabilities() []Capability
	// Send sends a message to the given endpoint. It returns an error if the message could not be sent.
	Send(msg *Message) error
}

type EndpointFactoryPort interface {
	// Create creates a new sender.
	Create(senderType string, config map[string]string) (EndpointPort, error)
}
