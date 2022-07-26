package endpoint

import (
	"context"
	"fmt"
	"text/template"
)

type (
	Attachment struct {
		Name    string
		Data    []byte
		IsImage bool
	}

	Endpoint struct {
		Name               string
		Type               string
		TextDisable        bool
		textTemplate       *template.Template
		AttachmentsDisable bool
		sender             Sender
	}

	CreateEndpointRequest struct {
		Name               string
		Type               string
		Config             Config
		TextDisable        bool
		TextTemplate       string
		AttachmentsDisable bool
	}

	Sender interface {
		// Send text and attachments to endpoint. Text can be empty and atts can be length 0 but not both.
		Send(ctx context.Context, text string, atts []Attachment) error
	}

	Service interface {
		CreateEndpoint(req CreateEndpointRequest) error
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

func NewEndpoint(name string, endpointType string, textDisable bool, textTemplateStr string, attachmentsDisable bool, sender Sender) (Endpoint, error) {
	textTemplate, err := template.New(name).Parse(textTemplateStr)
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		Name:               name,
		Type:               endpointType,
		sender:             sender,
		TextDisable:        textDisable,
		textTemplate:       textTemplate,
		AttachmentsDisable: attachmentsDisable,
	}, nil
}

func (e Endpoint) SendRaw(ctx context.Context, text string, atts []Attachment) error {
	// Don't send there is nothing to send
	if text == "" && len(atts) == 0 {
		return nil
	}

	return e.sender.Send(ctx, text, atts)
}

func (e Endpoint) TextTemplate() string {
	return e.textTemplate.Root.String()
}

func (c Config) Require(keys []string) error {
	for _, key := range keys {
		if _, ok := c[key]; !ok {
			return fmt.Errorf("missing key: %s", key)
		}
	}

	return nil
}

func OnlyImages(atts []Attachment) []Attachment {
	imgAtts := []Attachment{}
	for _, a := range atts {
		if a.IsImage {
			imgAtts = append(imgAtts, a)
		}
	}

	return imgAtts
}
