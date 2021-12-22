package app

// AuthServicePort handles authenticating users.
type AuthServicePort interface {
	AnonymousLogin() bool
	Login(username, password string) error
}

// BridgeServicePort handles finding endpoints for messages.
type BridgeServicePort interface {
	GetEndpoints(msg *Message) []EndpointPort
}

// MessageServicePort handles creating and sending messages.
type MessageServicePort interface {
	// Create creates a new message and saves it.
	Create(subject, from string, to map[string]bool, text string) (*Message, error)
	// AddAttachment adds an attachment to a message.
	AddAttachment(msg *Message, name string, content []byte) error
	// Send finds endpoints for the message and sends to it.
	Send(msg *Message) error
}

// MessageRepositoryPort handles storing messages.
type MessageRepositoryPort interface {
	// Create saves the message.
	Create(msg *Message) error
	// Update updates the message.
	Update(msg *Message) error
}

// EndpointPort handles sending messages to an endpoint.
type EndpointPort interface {
	// Send sends the message to the endpoint.
	Send(msg *Message) error
}
