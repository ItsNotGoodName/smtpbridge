package endpoint

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type (
	Attachment struct {
		Name    string
		Data    []byte
		IsImage bool
	}

	Endpoint struct {
		Name     string
		Type     string
		Sender   Sender
		template *template.Template
	}

	CreateEndpointRequest struct {
		Name     string
		Type     string
		Config   Config
		Template string
	}

	Sender interface {
		// Send text and attachments to endpoint. Text can be empty and atts can be length 0.
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

func NewEndpoint(name string, endpointType string, templateStr string, sender Sender) (Endpoint, error) {
	tmpl, err := template.New(name).Parse(templateStr)
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		Name:     name,
		Type:     endpointType,
		Sender:   sender,
		template: tmpl,
	}, nil
}

func (e Endpoint) Text(env *envelope.Envelope) (string, error) {
	var buffer bytes.Buffer
	if err := e.template.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (e Endpoint) SendText(ctx context.Context, text string) error {
	return e.Sender.Send(ctx, text, []Attachment{})
}

func NewAttachment(att *envelope.Attachment, data []byte) Attachment {
	return Attachment{
		Name:    att.Name,
		Data:    data,
		IsImage: att.IsImage(),
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

func FilterImages(atts []Attachment) []Attachment {
	imgAtts := []Attachment{}
	for _, a := range atts {
		if a.IsImage {
			imgAtts = append(imgAtts, a)
		}
	}

	return imgAtts
}
