package app

type AuthServicePort interface {
	AnonymousLogin() bool
	Login(username, password string) error
}

type BridgeServicePort interface {
	// GetEndpoints returns a list of endpoints that the message should be sent to.
	GetEndpoints(msg *Message) []EndpointPort
}

type MessageServicePort interface {
	// Handle receives message from drivers.
	Create(subject, from string, to map[string]bool, text string) (*Message, error)
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
