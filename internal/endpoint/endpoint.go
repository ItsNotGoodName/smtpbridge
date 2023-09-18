package endpoint

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
)

func validate(f Factory, r models.Endpoint) error {
	if r.Internal && !r.InternalID.Valid {
		return fmt.Errorf("internal id is empty")
	}

	if r.Name == "" {
		return models.FieldError{Field: models.FieldName, Err: fmt.Errorf("cannot be empty")}
	}

	_, err := f.Build(r)
	if err != nil {
		return err
	}

	return nil
}

func new(f Factory, r models.DTOEndpointCreate) models.Endpoint {
	return models.Endpoint{
		Internal:          false,
		InternalID:        sql.NullString{},
		Name:              r.Name,
		AttachmentDisable: r.AttachmentDisable,
		TextDisable:       r.TextDisable,
		TitleTemplate:     r.TitleTemplate,
		BodyTemplate:      r.BodyTemplate,
		Kind:              r.Kind,
		Config:            Schema.Filter(r.Kind, r.Config),
	}
}

func New(f Factory, r models.DTOEndpointCreate) (models.Endpoint, error) {
	end := new(f, r)
	return end, validate(f, end)
}

func NewInternal(f Factory, r models.DTOEndpointCreate, internalID string) (models.Endpoint, error) {
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

func Update(f Factory, m models.Endpoint, req models.DTOEndpointUpdate) (models.Endpoint, error) {
	if m.Internal {
		return models.Endpoint{}, models.ErrInternalResource
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.AttachmentDisable != nil {
		m.AttachmentDisable = *req.AttachmentDisable
	}
	if req.TextDisable != nil {
		m.TextDisable = *req.TextDisable
	}
	if req.TitleTemplate != nil {
		m.TitleTemplate = *req.TitleTemplate
	}
	if req.BodyTemplate != nil {
		m.BodyTemplate = *req.BodyTemplate
	}
	if req.Kind != nil {
		m.Kind = *req.Kind
	}
	if req.Config != nil {
		m.Config = Schema.Filter(m.Kind, *req.Config)
	}

	return m, validate(f, m)
}

func Delete(r models.Endpoint) error {
	if r.Internal {
		return models.ErrInternalResource
	}
	return nil
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

type Factory struct {
	pythonExecutable  string
	appriseScriptPath string
	funcMap           template.FuncMap
}

func NewFactory(pythonExecutable string, appriseScriptPath string, funcMap template.FuncMap) Factory {
	return Factory{
		pythonExecutable:  pythonExecutable,
		appriseScriptPath: appriseScriptPath,
		funcMap:           funcMap,
	}
}

func (s Factory) Build(e models.Endpoint) (Endpoint, error) {
	sender, err := s.build(e.Kind, e.Config)
	if err != nil {
		if errors.Is(err, errInvalidSenderKind) {
			return Endpoint{}, models.FieldError{Field: models.FieldKind, Err: err}
		}
		return Endpoint{}, models.FieldError{Field: models.FieldConfig, Err: err}
	}

	titleTemplate, err := template.New("").Funcs(s.funcMap).Parse(e.TitleTemplate)
	if err != nil {
		return Endpoint{}, models.FieldError{Field: models.FieldTitleTemplate, Err: err}
	}

	bodyTemplate, err := template.New("").Funcs(s.funcMap).Parse(e.BodyTemplate)
	if err != nil {
		return Endpoint{}, models.FieldError{Field: models.FieldBodyTemplate, Err: err}
	}

	return Endpoint{
		ID:     e.ID,
		sender: sender,
		config: config{
			TextDisable:        e.TextDisable,
			TitleTemplate:      titleTemplate,
			BodyTemplate:       bodyTemplate,
			AttachmentsDisable: e.AttachmentDisable,
		},
	}, nil
}
