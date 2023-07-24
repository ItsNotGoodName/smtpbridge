package endpoints

import (
	"fmt"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/rs/zerolog/log"
)

type CreateEndpoint struct {
	Internal          bool
	InternalID        string
	Name              string
	AttachmentDisable bool
	TextDisable       bool
	BodyTemplate      string
	Kind              string
	Config            map[string]string
}

type Endpoint struct {
	ID                int64 `bun:"id,pk,autoincrement"`
	Internal          bool
	InternalID        string
	Name              string
	AttachmentDisable bool
	TextDisable       bool
	BodyTemplate      string
	Kind              string
	Config            map[string]string
}

func New(r CreateEndpoint) (Endpoint, error) {
	if r.Internal && r.InternalID == "" {
		return Endpoint{}, fmt.Errorf("internal id is empty")
	}

	if r.Name == "" {
		return Endpoint{}, fmt.Errorf("name is empty")
	}

	if r.BodyTemplate == "" {
		r.BodyTemplate = "{{ .Message.Text }}"
	}

	end := Endpoint{
		Internal:          r.Internal,
		InternalID:        r.InternalID,
		Name:              r.Name,
		AttachmentDisable: r.AttachmentDisable,
		TextDisable:       r.TextDisable,
		BodyTemplate:      r.BodyTemplate,
		Kind:              r.Kind,
		Config:            r.Config,
	}

	_, err := end.Parse()
	return end, err
}

type Config struct {
	AttachmentsDisable bool
	TextDisable        bool
	BodyTemplate       *template.Template
}

type Sender interface {
	Send(cc core.Context, env envelope.Envelope, config Config) error
}

type ParsedEndpoint struct {
	ID     int64
	Config Config
	Sender Sender
}

func (e Endpoint) Parse() (ParsedEndpoint, error) {
	bodyTemplate, err := template.New("").Parse(e.BodyTemplate)
	if err != nil {
		return ParsedEndpoint{}, err
	}

	sender, err := SenderFactory(e.Kind, e.Config)
	if err != nil {
		return ParsedEndpoint{}, err
	}

	return ParsedEndpoint{
		ID: e.ID,
		Config: Config{
			TextDisable:        e.TextDisable,
			BodyTemplate:       bodyTemplate,
			AttachmentsDisable: e.AttachmentDisable,
		},
		Sender: sender,
	}, nil
}

func (pe ParsedEndpoint) Send(cc core.Context, env envelope.Envelope) error {
	if err := pe.Sender.Send(cc, env, pe.Config); err != nil {
		return err
	}

	log.Info().Int64("envelope-id", env.Message.ID).Int64("endpoint-id", pe.ID).Msg("Envelope sent to endpoint")

	return nil
}
