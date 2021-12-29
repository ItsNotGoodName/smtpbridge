package app

// AuthServicePort handles authenticating users.
type AuthServicePort interface {
	AnonymousLogin() bool
	Login(username, password string) error
}

// BridgeServicePort handles finding endpoints for messages.
type BridgeServicePort interface {
	// GetBridges returns a list of bridges that the message belongs to.
	GetBridges(msg *Message) []Bridge
}

// MessageServicePort handles creating and sending messages.
type MessageServicePort interface {
	// Create creates a new message and saves it.
	Create(subject, from string, to map[string]bool, text string) (*Message, error)
	// CreateAttachment adds an attachment to a message.
	CreateAttachment(msg *Message, name string, data []byte) (*Attachment, error)
	// AddAttachment adds an attachment to a message.
	Send(msg *Message) error
	SendBridges(msg *Message, bridges []Bridge) error
	UpdateStatus(msg *Message, status Status) error
}

// MessageRepositoryPort handles storing messages.
type MessageRepositoryPort interface {
	CreateMessage(msg *Message) error
	GetMessage(uuid string) (*Message, error)
	UpdateMessage(msg *Message, updateFN func(msg *Message) (*Message, error)) (*Message, error)
}

type AttachmentRepositoryPort interface {
	CreateAttachment(att *Attachment) error
	GetAttachment(uuid string) (*Attachment, error)
	LoadAttachment(msg *Message) error
}

// EndpointPort handles sending messages to an endpoint.
type EndpointPort interface {
	// Send sends the message to the endpoint.
	Send(msg *EndpointMessage) error
}

type EndpointRepositoryPort interface {
	Create(name, endpointType string, config map[string]string) error
	Get(name string) (EndpointPort, error)
}
