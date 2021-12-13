package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/go-smtpbridge/app"
)

func Factory(endpointType string, config map[string]string) (app.EndpointPort, error) {
	switch endpointType {
	case "telegram":
		token, ok := config["token"]
		if !ok {
			return nil, fmt.Errorf("%s token not found", endpointType)
		}
		chatID, ok := config["chat_id"]
		if !ok {
			return nil, fmt.Errorf("%s chat_id not found", endpointType)
		}

		return NewTelegram(token, chatID), nil
	}

	return nil, fmt.Errorf("%s endpoint not supported", endpointType)
}
