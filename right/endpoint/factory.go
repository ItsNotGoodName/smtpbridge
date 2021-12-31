package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

const (
	NameTelegram = "telegram"
	NameMock     = "mock"
)

func factory(endpointType string, config map[string]string) (domain.EndpointPort, error) {
	switch endpointType {
	case NameTelegram:
		return telegramFactory(config)
	case NameMock:
		return NewMock()
	}

	return nil, fmt.Errorf("%v: %s", domain.ErrEndpointInvalidType, endpointType)
}

func telegramFactory(config map[string]string) (*Telegram, error) {
	token, ok := config["token"]
	if !ok {
		return nil, fmt.Errorf("%v: %s: token not found", domain.ErrEndpointInvalidConfig, NameTelegram)
	}

	chatID, ok := config["chat_id"]
	if !ok {
		return nil, fmt.Errorf("%v: %s: chat_id not found", domain.ErrEndpointInvalidConfig, NameTelegram)
	}

	return NewTelegram(token, chatID), nil
}
