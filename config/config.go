package config

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rules"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/labstack/gommon/bytes"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabasePath         string
	AttachmentsDirectory string
	HTTPDisable          bool
	HTTPAddress          string
	HTTPBodyLimit        int
	SMTPDisable          bool
	SMTPAddress          string
	SMTPMaxMessageBytes  int
	Endpoints            []endpoints.Endpoint
	Rules                []rules.Rule
	RuleEndpoints        map[string][]string
	Config               *models.Config
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

	var attachmentsSize *int64
	if raw.Retention.AttachmentSize != nil {
		size, err := bytes.Parse(*raw.Retention.AttachmentSize)
		if err != nil {
			return Config{}, err
		}
		attachmentsSize = &size
	}
	var envelopeAge *time.Duration
	if raw.Retention.EnvelopeAge != nil {
		age, err := time.ParseDuration(*raw.Retention.EnvelopeAge)
		if err != nil {
			return Config{}, err
		}
		envelopeAge = &age
	}

	retentionPolicy := models.RetentionPolicy{
		EnvelopeCount:  raw.Retention.EnvelopeCount,
		AttachmentSize: attachmentsSize,
		EnvelopeAge:    envelopeAge,
		MinAge:         5 * time.Minute,
	}

	authSMTP := models.Auth{
		Username: raw.SMTP.Username,
		Password: raw.SMTP.Password,
	}

	authHTTP := models.Auth{
		Username: raw.HTTP.Username,
		Password: raw.HTTP.Password,
	}

	return Config{
		DatabasePath:         databasePath,
		AttachmentsDirectory: attachmentsDirectory,
		HTTPDisable:          raw.HTTP.Disable,
		HTTPAddress:          raw.HTTP.Host + ":" + strconv.Itoa(raw.HTTP.Port),
		HTTPBodyLimit:        int(maxBytesForEachPayload),
		SMTPDisable:          raw.SMTP.Disable,
		SMTPAddress:          raw.SMTP.Host + ":" + strconv.Itoa(raw.SMTP.Port),
		SMTPMaxMessageBytes:  int(maxBytesForEachPayload),
		Endpoints:            ends,
		Rules:                rrules,
		RuleEndpoints:        rulesToEndpoints,
		Config: &models.Config{
			RetentionPolicy: retentionPolicy,
			AuthSMTP:        authSMTP,
			AuthHTTP:        authHTTP,
		},
	}, nil
}

type Raw struct {
	MaxPayloadSize string `name:"max_payload_size" default:"25 MB"`
	DataDirectory  string `name:"data_directory" default:"smtpbridge_data" arg:""`
	Retention      struct {
		EnvelopeCount  *int    `name:"envelope_count"`
		EnvelopeAge    *string `name:"envelope_age"`
		AttachmentSize *string `name:"attachment_size"`
	} `embed:"" prefix:"retention-"`
	HTTP struct {
		Disable  bool
		Host     string
		Port     int `default:"8080"`
		Username string
		Password string
	} `embed:"" prefix:"http-"`
	SMTP struct {
		Disable  bool
		Host     string
		Port     int `default:"1025"`
		Username string
		Password string
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

const envConfigYamlKey = "SMTPBRIDGE_CONFIG_YAML"

func Read(cli CLI) (Raw, error) {
	var raw Raw
	envConfigYaml := os.Getenv(envConfigYamlKey)
	if envConfigYaml == "" {
		// Resolve config file
		var configFiles []string
		if cli.Config == nil {
			// DEPS: ../README.md
			configFile, err := resolve([]string{
				"config.yaml",
				"config.yml",
				".smtpbridge.yaml",
				".smtpbridge.yml",
				"~/.smtpbridge.yaml",
				"~/.smtpbridge.yml",
				"/etc/smtpbridge.yaml",
				"/etc/smtpbridge.yml",
			})
			if err != nil {
				return Raw{}, err
			}

			if configFile != "" {
				configFiles = []string{configFile}
			}
		} else if *cli.Config != "" {
			configFiles = []string{*cli.Config}
		}

		if len(configFiles) != 0 {
			// Load config file
			log.Info().Str("path", configFiles[0]).Msg("Reading config file")

			parser, err := kong.New(&raw, kong.Configuration(kongyaml.Loader, configFiles...))
			if err != nil {
				return Raw{}, err
			}

			_, err = parser.Parse([]string{})
			if err != nil {
				return Raw{}, err
			}
		}
	} else {
		// Load config file from ENV
		log.Info().Msgf("Reading config from environment variable %s", envConfigYamlKey)

		resolver, err := kongyaml.Loader(strings.NewReader(envConfigYaml))
		if err != nil {
			return Raw{}, err
		}

		parser, err := kong.New(&raw, kong.Resolvers(resolver))
		if err != nil {
			return Raw{}, err
		}

		_, err = parser.Parse([]string{})
		if err != nil {
			return Raw{}, err
		}
	}

	for endKey := range raw.Endpoints {
		for ruleKey, rrule := range raw.Rules {
			if ruleKey == endKey {
				rrule.Endpoints = append(rrule.Endpoints, endKey)
				raw.Rules[ruleKey] = rrule
			}
		}
	}

	if cli.DataDirectory != "" {
		raw.DataDirectory = cli.DataDirectory
	}

	if cli.SMTPDisable != nil {
		raw.SMTP.Disable = bool(*cli.SMTPDisable)
	}
	if cli.SMTPHost != nil {
		raw.SMTP.Host = string(*cli.SMTPHost)
	}
	if cli.SMTPPort != nil {
		raw.SMTP.Port = int(*cli.SMTPPort)
	}

	if cli.HTTPDisable != nil {
		raw.HTTP.Disable = bool(*cli.HTTPDisable)
	}
	if cli.HTTPHost != nil {
		raw.HTTP.Host = string(*cli.HTTPHost)
	}
	if cli.HTTPPort != nil {
		raw.HTTP.Port = int(*cli.HTTPPort)
	}

	return raw, nil
}

type CLI struct {
	Config        *string `name:"config" help:"Path to config file." type:"string"`
	DataDirectory string  `name:"data-directory" help:"Path to data directory." type:"path"`
	SMTPDisable   *bool   `name:"smtp-disable" help:"Disable SMTP server."`
	SMTPHost      *string `name:"smtp-host" help:"SMTP host address to listen on."`
	SMTPPort      *uint16 `name:"smtp-port" help:"SMTP port to listen on."`
	HTTPDisable   *bool   `name:"http-disable" help:"Disable HTTP server."`
	HTTPHost      *string `name:"http-host" help:"HTTP host address to listen on."`
	HTTPPort      *uint16 `name:"http-port" help:"HTTP port to listen on."`
	Version       bool    `name:"version" help:"Show version."`
}

func ReadAndParseCLI() CLI {
	cli := CLI{}
	kong.Parse(
		&cli,
		kong.Description("Bridge email to other messaging services."),
	)
	return cli
}
