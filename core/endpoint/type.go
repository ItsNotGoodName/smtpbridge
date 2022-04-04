package endpoint

import "context"

type (
	Envelope struct {
		Message     *Message
		Attachments []Attachment
	}

	Message struct {
		ID    int64  // ID of the message.
		Title string // Title of the message.
		Body  string // Body of the message.
	}

	Attachment struct {
		Name string // Name is the name of the attachment.
		Data []byte // Data is the attachment data.
	}

	Facade struct {
		Name     string   // Name is the name of the endpoint.
		Endpoint Endpoint // Endpoint is the endpoint.
	}

	Endpoint interface {
		// Send sends the message to the endpoint.
		Send(ctx context.Context, env Envelope) error
	}

	SendRequest struct {
		Envelope Envelope
		Facade   *Facade
	}

	SendResponse struct {
		Envelope Envelope
		Facade   *Facade
		Error    error
	}

	Service interface {
		// Send sends the envelope to the endpoint.
		Send(ctx context.Context, req []SendRequest) <-chan SendResponse
	}

	Repository interface {
		// Create a new endpoint.
		Create(name, endpointType string, config map[string]string) error
		// Get returns an endpoint by name.
		Get(name string) (*Facade, error)
	}
)
