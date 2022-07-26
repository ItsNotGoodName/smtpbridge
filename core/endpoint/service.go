package endpoint

import "github.com/ItsNotGoodName/smtpbridge/core"

type EndpointService struct {
	store Store
}

const (
	TypeConsole  string = "console"
	TypeTelegram string = "telegram"
)

func newSender(endpointType string, config Config) (Sender, error) {
	// Console
	if endpointType == TypeConsole {
		return &Console{}, nil
	}

	// Telegram
	if endpointType == TypeTelegram {
		if err := config.Require([]string{"token", "chat_id"}); err != nil {
			return nil, err
		}
		return NewTelegram(config["token"], config["chat_id"]), nil
	}

	return nil, core.ErrEndpointTypeInvalid
}

func NewEndpointService(store Store) *EndpointService {
	return &EndpointService{
		store: store,
	}
}

func (es *EndpointService) CreateEndpoint(req CreateEndpointRequest) error {
	sender, err := newSender(req.Type, req.Config)
	if err != nil {
		return err
	}

	end, err := NewEndpoint(req.Name, req.Type, req.TextDisable, req.TextTemplate, req.AttachmentsDisable, sender)
	if err != nil {
		return err
	}

	return es.store.CreateEndpoint(end)
}

func (es *EndpointService) GetEndpoint(name string) (Endpoint, error) {
	return es.store.GetEndpoint(name)
}

func (es *EndpointService) ListEndpoint() []Endpoint {
	return es.store.ListEndpoint()
}
