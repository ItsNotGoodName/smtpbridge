package app

import (
	"io/fs"
)

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
	// UpdateStatus updates the status of a message.
	UpdateStatus(msg *Message, status Status) error
	// List messages with attachments.
	List(limit, offset int) ([]Message, error)
}

type EndpointServicePort interface {
	Send(msg *Message) error
	SendBridges(msg *Message, bridges []Bridge) error
}

// MessageRepositoryPort handles storing messages.
type MessageRepositoryPort interface {
	// CreateMessage saves a new message.
	CreateMessage(msg *Message) error
	// GetMessage returns a message by it's UUID.
	GetMessage(uuid string) (*Message, error)
	// GetMessages returns a list of messages.
	GetMessages(limit, offset int) ([]Message, error)
	// UpdateMessage updates a message.
	UpdateMessage(msg *Message, updateFN func(msg *Message) (*Message, error)) error
}

type AttachmentRepositoryPort interface {
	// CreateAttachment saves a new attachment.
	CreateAttachment(att *Attachment) error
	// GetAttachment returns an attachment by it's UUID.
	GetAttachment(uuid string) (*Attachment, error)
	// GetAttachmentData returns the data for an attachment.
	GetAttachmentData(att *Attachment) ([]byte, error)
	// GetAttachments returns a list of attachments for a message.
	GetAttachments(msg *Message) ([]Attachment, error)
	// GetAttachmentFile returns the filename of the attachment
	GetAttachmentFile(att *Attachment) string
	// GetFS returns the attachment file system
	GetAttachmentFS() fs.FS
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
