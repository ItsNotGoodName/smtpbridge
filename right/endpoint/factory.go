package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

const (
	NameTelegram = "telegram"
	NameMock     = "mock"
)

func factory(endpointType string, config map[string]string) (core.EndpointPort, error) {
	switch endpointType {
	case NameTelegram:
		return telegramFactory(config)
	case NameMock:
		return NewMock()
	}

	return nil, fmt.Errorf("%v: %s", core.ErrEndpointInvalidType, endpointType)
}

func telegramFactory(config map[string]string) (*Telegram, error) {
	token, ok := config["token"]
	if !ok {
		return nil, fmt.Errorf("%v: %s: token not found", core.ErrEndpointInvalidConfig, NameTelegram)
	}

	chatID, ok := config["chat_id"]
	if !ok {
		return nil, fmt.Errorf("%v: %s: chat_id not found", core.ErrEndpointInvalidConfig, NameTelegram)
	}

	return NewTelegram(token, chatID), nil
}
