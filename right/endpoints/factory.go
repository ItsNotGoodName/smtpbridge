package endpoints

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
)

const (
	TypeTelegram = "telegram"
	TypeMock     = "mock"
)

func factory(endpointType string, config map[string]string) (endpoint.Endpoint, error) {
	switch endpointType {
	case TypeTelegram:
		return telegramFactory(config)
	case TypeMock:
		return NewMock()
	}

	return nil, fmt.Errorf("%v: %s", endpoint.ErrInvalidType, endpointType)
}

func telegramFactory(config map[string]string) (*Telegram, error) {
	token, ok := config["token"]
	if !ok {
		return nil, fmt.Errorf("%v: %w: token not found", TypeTelegram, endpoint.ErrInvalidConfig)
	}

	chatID, ok := config["chat_id"]
	if !ok {
		return nil, fmt.Errorf("%v: %w: chat_id not found", TypeTelegram, endpoint.ErrInvalidConfig)
	}

	return NewTelegram(token, chatID), nil
}
