package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
	"github.com/containrrr/shoutrrr"
)

var Schema models.EndpointSchema = models.EndpointSchema{
	{
		Name: "Console",
		Kind: "console",
	},
	{
		Name: "Telegram",
		Kind: "telegram",
		Fields: []models.EndpointSchemaField{
			{
				Name:        "Token",
				Description: "",
				Key:         "token",
			},
			{
				Name:        "Chat ID",
				Description: "",
				Key:         "chat_id",
			},
		},
	},
	{
		Name: "Shoutrrr",
		Kind: "shoutrrr",
		Fields: []models.EndpointSchemaField{
			{
				Name:        "URLs",
				Key:         "urls",
				Description: "List of URLs.",
				Multiline:   true,
			},
		},
	},
	{
		Name: "Apprise",
		Kind: "apprise",
		Fields: []models.EndpointSchemaField{
			{
				Name:        "URLs",
				Key:         "urls",
				Description: "List of URLs.",
				Multiline:   true,
			},
		},
	},
}

var errInvalidSenderKind = fmt.Errorf("invalid sender kind")

func (s Factory) build(kind string, config models.EndpointConfig) (Sender, error) {
	switch kind {
	case "console":
		return senders.NewConsole(), nil
	case "telegram":
		token := config.Str("token")
		if token == "" {
			return nil, fmt.Errorf("token empty")
		}

		chatID := config.Str("chat_id")
		if chatID == "" {
			return nil, fmt.Errorf("chat_id empty")
		}

		return senders.NewTelegram(token, chatID), nil
	case "shoutrrr":
		router, err := shoutrrr.CreateSender(config.StrSlice("urls")...)
		if err != nil {
			return nil, err
		}

		return senders.NewShoutrrr(router), nil
	case "apprise":
		return senders.NewApprise(s.pythonExecutable, s.appriseScriptPath, config.StrSlice("urls")), nil
	default:
		return nil, errInvalidSenderKind
	}
}
