package endpoint

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
)

func validate(f Factory, end models.Endpoint) error {
	if end.Internal && !end.InternalID.Valid {
		return fmt.Errorf("internal id is empty")
	}

	_, err := f.Build(end)
	if err != nil {
		return err
	}

	return nil
}

type CreateEndpoint struct {
	Name              string
	AttachmentDisable bool
	TextDisable       bool
	TitleTemplate     string
	BodyTemplate      string
	Kind              string
	Config            models.EndpointConfig
}

const DefaultTitleTemplate = "{{ .Message.Subject }}"
const DefaultBodyTemplate = "{{ .Message.Text }}"

func new(f Factory, r CreateEndpoint) models.Endpoint {
	return models.Endpoint{
		Internal:          false,
		InternalID:        sql.NullString{},
		Name:              r.Name,
		AttachmentDisable: r.AttachmentDisable,
		TextDisable:       r.TextDisable,
		TitleTemplate:     r.TitleTemplate,
		BodyTemplate:      r.BodyTemplate,
		Kind:              r.Kind,
		Config:            r.Config,
	}
}

func New(f Factory, r CreateEndpoint) (models.Endpoint, error) {
	end := new(f, r)
	return end, validate(f, end)
}

func NewInternal(f Factory, r CreateEndpoint, internalID string) (models.Endpoint, error) {
	if r.Name == "" {
		r.Name = internalID
	}

	end := new(f, r)

	end.Internal = true
	end.InternalID = sql.NullString{
		String: internalID,
		Valid:  true,
	}

	return end, validate(f, end)
}

type Sender interface {
	Send(ctx context.Context, env models.Envelope, tr senders.Transformer) error
}

type config struct {
	AttachmentsDisable bool
	TextDisable        bool
	TitleTemplate      *template.Template
	BodyTemplate       *template.Template
}

type Endpoint struct {
	ID     int64
	sender Sender
	config config
}

type FileStore interface {
	Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error)
	Path(ctx context.Context, att models.Attachment) (string, error)
}

func (e Endpoint) Send(ctx context.Context, fileStore FileStore, env models.Envelope) error {
	return e.sender.Send(ctx, env, payload{FileStore: fileStore, config: e.config})
}

type payload struct {
	FileStore
	config config
}

func (p payload) Body(ctx context.Context, env models.Envelope) (string, error) {
	if p.config.TextDisable {
		return "", nil
	}

	var buffer bytes.Buffer
	if err := p.config.BodyTemplate.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (p payload) Files(ctx context.Context, env models.Envelope) ([]models.Attachment, error) {
	if p.config.AttachmentsDisable {
		return []models.Attachment{}, nil
	}
	return env.Attachments, nil
}

func (p payload) Title(ctx context.Context, env models.Envelope) (string, error) {
	if p.config.TextDisable {
		return "", nil
	}

	var buffer bytes.Buffer
	if err := p.config.TitleTemplate.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
