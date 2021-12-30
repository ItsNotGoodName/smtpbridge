package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

func factory(endpointType string, config map[string]string) (domain.EndpointPort, error) {
	switch endpointType {
	case "telegram":
		return telegramFactory(config)
	case "mock":
		return NewMock()
	}

	return nil, fmt.Errorf("%s: %v", endpointType, domain.ErrEndpointInvalidType)
}

func telegramFactory(config map[string]string) (*Telegram, error) {
	token, ok := config["token"]
	if !ok {
		return nil, fmt.Errorf("telegram token not found: %v", domain.ErrEndpointInvalidConfig)
	}

	chatID, ok := config["chat_id"]
	if !ok {
		return nil, fmt.Errorf("telegram chat_id not found: %v", domain.ErrEndpointInvalidConfig)
	}

	return NewTelegram(token, chatID), nil
}
