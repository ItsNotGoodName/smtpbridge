package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/labstack/gommon/bytes"
)

type Config struct {
	DatabasePath         string
	AttachmentsDirectory string
	HTTPAddress          string
	HTTPBodyLimit        int
	SMTPAddress          string
	SMTPMaxMessageBytes  int
	Endpoints            []endpoints.Endpoint
	Rules                []rules.Rule
	RuleEndpoints        map[string][]string
	RetentionPolicy      models.RetentionPolicy
}

func Parse(raw Raw) (Config, error) {
	dataDirectory := raw.DataDirectory
	if !path.IsAbs(raw.DataDirectory) {
		cwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}

		dataDirectory = path.Join(cwd, raw.DataDirectory)
	}

	if err := os.MkdirAll(dataDirectory, 0755); err != nil {
		return Config{}, err
	}

	databasePath := path.Join(dataDirectory, "smtpbridge.db")

	attachmentsDirectory := path.Join(dataDirectory, "attachments")
	if err := os.MkdirAll(attachmentsDirectory, 0755); err != nil {
		return Config{}, err
	}

	maxBytesForEachPayload, err := bytes.Parse(raw.MaxPayloadSize)
	if err != nil {
		return Config{}, err
	}

	var ends []endpoints.Endpoint
	for k, v := range raw.Endpoints {
		name := v.Name
		if name == "" {
			name = k
		}
		end, err := endpoints.New(
			endpoints.CreateEndpoint{
				Internal:          true,
				InternalID:        k,
				Name:              name,
				TextDisable:       v.Text_Disable,
				AttachmentDisable: v.Attachment_Disable,
				BodyTemplate:      v.Body_Template,
				Kind:              v.Kind,
				Config:            v.Config,
			},
		)
		if err != nil {
			return Config{}, err
		}

		ends = append(ends, end)
	}

	rulesToEndpoints := make(map[string][]string)
	var rrules []rules.Rule
	for k, rr := range raw.Rules {
		name := rr.Name
		if name == "" {
			name = k
		}
		rule, err := rules.New(rules.CreateRule{
			Internal:   true,
			InternalID: k,
			Name:       name,
			Expression: rr.Expression,
		})
		if err != nil {
			return Config{}, err
		}

		rulesToEndpoints[k] = rr.Endpoints
		rrules = append(rrules, rule)
	}

	var attachmentsSize int64
	if raw.Retention.AttachmentSize != "" {
		var err error
		attachmentsSize, err = bytes.Parse(raw.Retention.AttachmentSize)
		if err != nil {
			return Config{}, err
		}
	}
	var envelopeAge time.Duration
	if raw.Retention.EnvelopeAge != "" {
		var err error
		envelopeAge, err = time.ParseDuration(raw.Retention.EnvelopeAge)
		if err != nil {
			return Config{}, err
		}
	}
	retentionPolicy := models.RetentionPolicy{
		EnvelopeCount:  raw.Retention.EnvelopeCount,
		AttachmentSize: attachmentsSize,
		EnvelopeAge:    envelopeAge,
	}

	return Config{
		DatabasePath:         databasePath,
		AttachmentsDirectory: attachmentsDirectory,
		HTTPAddress:          raw.HTTP.Host + ":" + strconv.Itoa(raw.HTTP.Port),
		HTTPBodyLimit:        int(maxBytesForEachPayload),
		SMTPAddress:          raw.SMTP.Host + ":" + strconv.Itoa(raw.SMTP.Port),
		SMTPMaxMessageBytes:  int(maxBytesForEachPayload),
		Endpoints:            ends,
		Rules:                rrules,
		RuleEndpoints:        rulesToEndpoints,
		RetentionPolicy:      retentionPolicy,
	}, nil
}

type Raw struct {
	MaxPayloadSize string `name:"max_payload_size" default:"25 MB"`
	DataDirectory  string `name:"data_directory" default:"smtpbridge_data" arg:""`
	Retention      struct {
		EnvelopeCount  int    `name:"envelope_count"`
		EnvelopeAge    string `name:"envelope_age"`
		AttachmentSize string `name:"attachment_size"`
	} `embed:"" prefix:"retention-"`
	HTTP struct {
		Disable bool
		Host    string
		Port    int `default:"8080"`
	} `embed:"" prefix:"http-"`
	SMTP struct {
		Disable bool
		Host    string
		Port    int `default:"1025"`
	} `embed:"" prefix:"smtp-"`
	Endpoints map[string]RawEndpoint
	Rules     map[string]RawRule
}

type RawEndpoint struct {
	Name               string
	Kind               string
	Text_Disable       bool
	Body_Template      string
	Attachment_Disable bool
	Config             map[string]string
}

type RawRule struct {
	Name       string
	Expression string
	Endpoints  []string
}

type CLI struct {
	DataDirectory string `name:"data-directory" help:"Path to store data." type:"path" optional:""`
}

func Read() (Raw, error) {
	var raw Raw
	parser, err := kong.New(&raw, kong.Configuration(
		kongyaml.Loader,
		"config.yaml",
		"config.yml",
		".smtpbridge.yaml",
		".smtpbridge.yml",
		"~/.smtpbridge.yaml",
		"~/.smtpbridge.yml",
	))
	if err != nil {
		return Raw{}, err
	}

	_, err = parser.Parse([]string{})
	if err != nil {
		return Raw{}, err
	}

	for k, v := range raw.Endpoints {
		if v.Body_Template == "" {
			v.Body_Template = `SUBJECT: {{ .Message.Subject }}
FROM: {{ .Message.From }}

{{ .Message.Text }}`
			raw.Endpoints[k] = v
		}
	}

	cli := CLI{}

	kong.Parse(&cli)

	if cli.DataDirectory != "" {
		raw.DataDirectory = cli.DataDirectory
	}

	return raw, nil
}
