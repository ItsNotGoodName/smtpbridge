package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

func Factory(endpointType string, config map[string]string) (app.EndpointPort, error) {
	switch endpointType {
	case "telegram":
		return telegramFactory(config)
	}

	return nil, fmt.Errorf("%s endpoint not supported", endpointType)
}

func telegramFactory(config map[string]string) (*Telegram, error) {
	token, ok := config["token"]
	if !ok {
		return nil, fmt.Errorf("telegram token not found")
	}
	chatID, ok := config["chat_id"]
	if !ok {
		return nil, fmt.Errorf("telegram chat_id not found")
	}
	return NewTelegram(token, chatID), nil
}
