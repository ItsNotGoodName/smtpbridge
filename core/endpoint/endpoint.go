package endpoint

import (
	"context"
	"fmt"
)

type (
	Attachment struct {
		Name string
		Data []byte
	}

	Endpoint struct {
		Name   string
		Type   string
		Sender Sender
	}

	Sender interface {
		Send(ctx context.Context, text string, atts []Attachment) error
	}

	Service interface {
		CreateEndpoint(name string, endpointType string, config Config) error
		GetEndpoint(name string) (Endpoint, error)
		ListEndpoint() []Endpoint
	}

	Store interface {
		CreateEndpoint(endpoint Endpoint) error
		GetEndpoint(name string) (Endpoint, error)
		ListEndpoint() []Endpoint
	}

	Config map[string]string
)

func NewEndpoint(name string, endpointType string, sender Sender) Endpoint {
	return Endpoint{
		Name:   name,
		Type:   endpointType,
		Sender: sender,
	}
}

func (c Config) Require(keys []string) error {
	for _, key := range keys {
		if _, ok := c[key]; !ok {
			return fmt.Errorf("missing key: %s", key)
		}
	}

	return nil
}
