package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/senders"
	"github.com/containrrr/shoutrrr"
	securejoin "github.com/cyphar/filepath-securejoin"
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
				Key:         "token",
				Name:        "Token",
				Description: "Bot token from BotFather.",
				Example:     "0123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			{
				Key:         "chat_id",
				Name:        "Chat ID",
				Description: "Channel/chat/group ID.",
				Example:     "1000000000",
			},
		},
	},
	{
		Name: "Shoutrrr",
		Kind: "shoutrrr",
		Fields: []models.EndpointSchemaField{
			{
				Key:         "urls",
				Name:        "URLs",
				Description: "List of Shoutrrr URLs.",
				Example:     "telegram://token@telegram?chats=channel-1[,chat-id-1,...]",
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
				Description: "List of Apprise URLs.",
				Example:     "tgram://bottoken/ChatID",
				Multiline:   true,
			},
		},
	},
	{
		Name: "Script",
		Kind: "script",
		Fields: []models.EndpointSchemaField{
			{
				Key:         "file",
				Name:        "File",
				Description: "Name of script file located in the script directory.",
				Example:     "my-script.py",
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
	case "script":
		scriptPath, err := securejoin.SecureJoin(s.scriptDirectory, config.Str("file"))
		if err != nil {
			return nil, err
		}

		return senders.NewScript(scriptPath), nil
	default:
		return nil, errInvalidSenderKind
	}
}
