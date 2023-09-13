package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/auth"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/rule"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/basicflag"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/gommon/bytes"
	"github.com/rs/zerolog/log"
)

const (
	TimeFormat12H = "12h"
	TimeFormat24H = "24h"
)

type Config struct {
	Debug                bool
	TimeHourFormat       string
	DatabasePath         string
	AttachmentsDirectory string
	CSRFSecretPath       string
	SessionSecretPath    string
	SessionsDirectory    string
	HTTPDisable          bool
	HTTPAddress          string
	HTTPPort             uint16
	HTTPBodyLimit        int64
	HTTPURL              string
	SMTPDisable          bool
	SMTPAddress          string
	SMTPMaxMessageBytes  int64
	// IMAPDisable            bool
	// IMAPAddress            string
	Config                 *models.Config
	EndpointFactory        endpoint.Factory
	InternalEndpoints      []models.Endpoint
	InternalRules          []models.Rule
	InternalRuleToEndpoint map[string][]string
}

type Raw struct {
	Debug                   bool   `koanf:"debug"`
	TimeFormat              string `koanf:"time_format"`
	MaxPayloadSize          string `koanf:"max_payload_size"`
	DataDirectory           string `koanf:"data_directory"`
	PythonExecutable        string `koanf:"python_executable"`
	RetentionEnvelopeCount  string `koanf:"retention.envelope_count"`
	RetentionEnvelopeAge    string `koanf:"retention.envelope_age"`
	RetentionAttachmentSize string `koanf:"retention.attachment_size"`
	SMTPDisable             bool   `koanf:"smtp.disable"`
	SMTPHost                string `koanf:"smtp.host"`
	SMTPPort                uint16 `koanf:"smtp.port"`
	SMTPUsername            string `koanf:"smtp.username"`
	SMTPPassword            string `koanf:"smtp.password"`
	HTTPDisable             bool   `koanf:"http.disable"`
	HTTPHost                string `koanf:"http.host"`
	HTTPPort                uint16 `koanf:"http.port"`
	HTTPUsername            string `koanf:"http.username"`
	HTTPPassword            string `koanf:"http.password"`
	HTTPURL                 string `koanf:"http.url"`
	// IMAPDisable             bool    `koanf:"imap.disable"`
	// IMAPHost                string  `koanf:"imap.host"`
	// IMAPPort                uint16  `koanf:"imap.port"`
	Endpoints map[string]RawEndpoint
	Rules     map[string]RawRule
}

type RawEndpoint struct {
	Name              string            `koanf:"name"`
	Kind              string            `koanf:"kind"`
	TextDisable       bool              `koanf:"text_disable"`
	TitleTemplate     string            `koanf:"title_template"`
	BodyTemplate      string            `koanf:"body_template"`
	AttachmentDisable bool              `koanf:"attachment_disable"`
	Config            map[string]string `koanf:"config"`
}

type RawRule struct {
	Name       string   `koanf:"name"`
	Expression string   `koanf:"expression"`
	Endpoints  []string `koanf:"endpoints"`
}

var RawDefault = struct {
	TimeFormat       string `koanf:"time_format"`
	MaxPayloadSize   string `koanf:"max_payload_size"`
	DataDirectory    string `koanf:"data_directory"`
	PythonExecutable string `koanf:"python_executable"`
	SMTPPort         uint16 `koanf:"smtp.port"`
	HTTPPort         uint16 `koanf:"http.port"`
	// IMAPPort         uint16 `koanf:"imap.port"`
}{
	TimeFormat:       TimeFormat12H,
	MaxPayloadSize:   "25 MB",
	DataDirectory:    "smtpbridge_data",
	PythonExecutable: "python3",
	SMTPPort:         1025,
	HTTPPort:         8080,
	// IMAPPort:         10143,
}

var flagFlatKeys map[string]string = map[string]string{
	"time-format":       "time_format",
	"data-directory":    "data_directory",
	"python-executable": "python_executable",
}

func WithFlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.String("config", "", flagUsageString("", "Path to config file."))
	flags.String("time-format", "", flagUsageString(TimeFormat12H, fmt.Sprintf("Format for time display (%s/%s).", TimeFormat12H, TimeFormat24H)))
	flags.String("data-directory", "", flagUsageString(RawDefault.DataDirectory, "Path to data directory."))
	flags.String("python-executable", "", flagUsageString(RawDefault.PythonExecutable, "Python executable."))
	flags.Bool("debug", false, flagUsageBool(false, "Run in debug mode."))

	flags.Bool("smtp-disable", false, flagUsageBool(false, "Disable SMTP server."))
	flags.String("smtp-host", "", flagUsageString("", "SMTP host address to listen on."))
	flags.Int("smtp-port", 0, flagUsageInt(int(RawDefault.SMTPPort), "SMTP port to listen on."))

	flags.Bool("http-disable", false, flagUsageBool(false, "Disable HTTP server."))
	flags.String("http-host", "", flagUsageString("", "HTTP host address to listen on."))
	flags.Int("http-port", 0, flagUsageInt(int(RawDefault.HTTPPort), "HTTP port to listen on."))
	flags.Int("http-url", 0, flagUsageString("", "HTTP public URL (e.g. http://127.0.0.1:8080)."))

	// flags.Bool("imap-disable", false, flagUsageBool(false, "Disable IMAP server."))
	// flags.String("imap-host", "", flagUsageString("", "IMAP host address to listen on."))
	// flags.Int("imap-port", 0, flagUsageInt(int(RawDefault.IMAPPort), "HTTP port to listen on."))

	return flags
}

func (p Parser) Parse(raw Raw) (Config, error) {
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

	csrfSecretPath := path.Join(dataDirectory, "csrf.secret")

	sessionSecretPath := path.Join(dataDirectory, "session.secret")

	appriseScriptPath := path.Join(dataDirectory, "apprise_script.py")
	err := senders.AppriseWriteScript(appriseScriptPath)
	if err != nil {
		return Config{}, err
	}

	sessionsDirectory := path.Join(dataDirectory, "sessions")
	if err := os.MkdirAll(sessionsDirectory, 0755); err != nil {
		return Config{}, err
	}

	attachmentsDirectory := path.Join(dataDirectory, "attachments")
	if err := os.MkdirAll(attachmentsDirectory, 0755); err != nil {
		return Config{}, err
	}

	maxBytesForEachPayload, err := bytes.Parse(raw.MaxPayloadSize)
	if err != nil {
		return Config{}, err
	}

	var config *models.Config
	{
		var envelopeCount *int
		if raw.RetentionEnvelopeCount != "" {
			count, err := strconv.Atoi(raw.RetentionEnvelopeCount)
			if err != nil {
				return Config{}, err
			}

			envelopeCount = &count
		}
		var attachmentsSize *int64
		if raw.RetentionAttachmentSize != "" {
			size, err := bytes.Parse(raw.RetentionAttachmentSize)
			if err != nil {
				return Config{}, err
			}
			attachmentsSize = &size
		}
		var envelopeAge *time.Duration
		if raw.RetentionEnvelopeAge != "" {
			age, err := time.ParseDuration(raw.RetentionEnvelopeAge)
			if err != nil {
				return Config{}, err
			}
			envelopeAge = &age
		}

		config = &models.Config{
			RetentionPolicy: models.ConfigRetentionPolicy{
				EnvelopeCount:  envelopeCount,
				AttachmentSize: attachmentsSize,
				EnvelopeAge:    envelopeAge,
				MinAge:         5 * time.Minute,
			},
			AuthSMTP: auth.New(
				raw.SMTPUsername,
				raw.SMTPPassword,
			),
			AuthHTTP: auth.New(
				raw.HTTPUsername,
				raw.HTTPPassword,
			),
		}
	}

	endpointFactory := endpoint.NewFactory(raw.PythonExecutable, appriseScriptPath, endpoint.NewFuncMap(endpoint.CreateFuncMap{
		URL: raw.HTTPURL,
	}))
	var endpoints []models.Endpoint
	{
		for key, value := range raw.Endpoints {
			e, err := endpoint.NewInternal(endpointFactory, endpoint.CreateEndpoint{
				Name:              value.Name,
				AttachmentDisable: value.AttachmentDisable,
				TextDisable:       value.TextDisable,
				TitleTemplate:     value.TitleTemplate,
				BodyTemplate:      value.BodyTemplate,
				Kind:              value.Kind,
				Config:            value.Config,
			}, key)
			if err != nil {
				return Config{}, err
			}

			endpoints = append(endpoints, e)
		}
	}

	ruleToEndpoints := make(map[string][]string)
	var rules []models.Rule
	{
		for key, value := range raw.Rules {
			r, err := rule.NewInternal(models.DTORuleCreate{
				Name:       value.Name,
				Expression: value.Expression,
			}, key)
			if err != nil {
				return Config{}, err
			}

			ruleToEndpoints[key] = value.Endpoints
			rules = append(rules, r)
		}

		// Special case where if the keys of rules and endpoints match then we should assume the user wants them to be connected
		for endKey := range raw.Endpoints {
			for ruleKey := range raw.Rules {
				if ruleKey == endKey {
					ruleToEndpoints[ruleKey] = append(ruleToEndpoints[ruleKey], endKey)
				}
			}
		}
	}

	var timeHourFormat string
	switch raw.TimeFormat {
	case TimeFormat12H:
		timeHourFormat = helpers.TimeHourFormat12
	case TimeFormat24H:
		timeHourFormat = helpers.TimeHourFormat24
	default:
		return Config{}, fmt.Errorf("invalid time format: %s", raw.TimeFormat)
	}

	httpAddress := raw.HTTPHost + ":" + strconv.Itoa(int(raw.HTTPPort))

	smtpAddress := raw.SMTPHost + ":" + strconv.Itoa(int(raw.SMTPPort))

	// imapAddress := raw.IMAPHost + ":" + strconv.Itoa(int(raw.IMAPPort))

	return Config{
		Debug:                raw.Debug,
		TimeHourFormat:       timeHourFormat,
		DatabasePath:         databasePath,
		CSRFSecretPath:       csrfSecretPath,
		SessionSecretPath:    sessionSecretPath,
		SessionsDirectory:    sessionsDirectory,
		AttachmentsDirectory: attachmentsDirectory,
		HTTPDisable:          raw.HTTPDisable,
		HTTPAddress:          httpAddress,
		HTTPPort:             raw.HTTPPort,
		HTTPBodyLimit:        maxBytesForEachPayload,
		HTTPURL:              raw.HTTPURL,
		SMTPDisable:          raw.SMTPDisable,
		SMTPAddress:          smtpAddress,
		SMTPMaxMessageBytes:  maxBytesForEachPayload,
		// IMAPDisable:            raw.IMAPDisable,
		// IMAPAddress:            imapAddress,
		Config:                 config,
		EndpointFactory:        endpointFactory,
		InternalEndpoints:      endpoints,
		InternalRules:          rules,
		InternalRuleToEndpoint: ruleToEndpoints,
	}, nil
}

type Parser struct {
	k *koanf.Koanf
}

func NewParser(flags *flag.FlagSet) (Parser, error) {
	var k = koanf.New(".")

	// Load defaults
	k.Load(structs.ProviderWithDelim(RawDefault, "koanf", "."), nil)

	if envConfig := os.Getenv("SMTPBRIDGE_CONFIG_YAML"); envConfig == "" {
		// Load YAML file
		var configFile string
		if p := flags.Lookup("config"); p != nil && p.Value.String() != "" {
			// Config file from flag
			configFile = p.Value.String()
		} else {
			// Config file from default
			var err error
			configFile, err = resolve([]string{
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
				return Parser{}, err
			}
		}
		if configFile != "" {
			log.Info().Str("path", configFile).Msg("Reading config from file")
			if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
				return Parser{}, err
			}
		}
	} else {
		// Load YAML env
		if envConfig != "" {
			log.Info().Str("env", "SMTPBRIDGE_CONFIG_YAML").Msgf("Reading config from environment")
			k.Load(rawbytes.Provider([]byte(envConfig)), yaml.Parser())
		}
	}

	// Load flag
	k.Load(basicflag.ProviderWithValue(flags, "-", func(key, value string) (string, interface{}) {
		if value == "" || value == "0" || value == "false" {
			return "", nil
		}
		if remap, ok := flagFlatKeys[key]; ok {
			return remap, ok
		}
		return key, value
	}), nil)

	return Parser{k: k}, nil
}

func (p Parser) Read() (Raw, error) {
	raw := Raw{}
	err := p.k.UnmarshalWithConf("", &raw, koanf.UnmarshalConf{Tag: "koanf", FlatPaths: true})
	if err != nil {
		return Raw{}, err
	}

	p.k.Unmarshal("endpoints", &raw.Endpoints)
	p.k.Unmarshal("rules", &raw.Rules)

	return raw, nil
}
