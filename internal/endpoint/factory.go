package endpoint

import (
	"fmt"
	"html/template"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
	"github.com/containrrr/shoutrrr"
)

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
		return Endpoint{}, err
	}

	titleTemplate, err := template.New("").Funcs(s.funcMap).Parse(e.TitleTemplate)
	if err != nil {
		return Endpoint{}, err
	}

	bodyTemplate, err := template.New("").Funcs(s.funcMap).Parse(e.BodyTemplate)
	if err != nil {
		return Endpoint{}, err
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

func (s Factory) build(kind string, config models.EndpointConfig) (Sender, error) {
	switch kind {
	case "console":
		return senders.NewConsole(), nil
	case "telegram":
		token := config.Str("token")
		if token == "" {
			return nil, fmt.Errorf("telegram: 'token' empty")
		}

		chatID := config.Str("chat_id")
		if chatID == "" {
			return nil, fmt.Errorf("telegram: 'chat_id' empty")
		}

		return senders.NewTelegram(token, chatID), nil
	case "shoutrrr":
		router, err := shoutrrr.CreateSender(config.StrSlice("urls")...)
		if err != nil {
			return nil, fmt.Errorf("shoutrrr: %w", err)
		}

		return senders.NewShoutrrr(router), nil
	case "apprise":
		return senders.NewApprise(s.pythonExecutable, s.appriseScriptPath, config.StrSlice("urls")), nil
	default:
		return nil, fmt.Errorf("invalid sender kind: %s", kind)
	}
}
