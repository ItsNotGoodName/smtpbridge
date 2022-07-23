package endpoint

import "github.com/ItsNotGoodName/smtpbridge/core"

const (
	TypeConsole  string = "console"
	TypeTelegram string = "telegram"
)

func newSender(endpointType string, config Config) (Sender, error) {
	if endpointType == TypeConsole {
		return &Console{}, nil
	}

	if endpointType == TypeTelegram {
		if err := config.Require([]string{"token", "chat_id"}); err != nil {
			return nil, err
		}
		return NewTelegram(config["token"], config["chat_id"]), nil
	}

	return nil, core.ErrEndpointTypeInvalid
}

type EndpointService struct {
	store Store
}

func NewEndpointService(store Store) *EndpointService {
	return &EndpointService{
		store: store,
	}
}

func (es *EndpointService) CreateEndpoint(name string, endpointType string, config Config) error {
	sender, err := newSender(endpointType, config)
	if err != nil {
		return err
	}

	return es.store.CreateEndpoint(NewEndpoint(name, endpointType, sender))
}

func (es *EndpointService) GetEndpoint(name string) (Endpoint, error) {
	return es.store.GetEndpoint(name)
}

func (es *EndpointService) ListEndpoint() []Endpoint {
	return es.store.ListEndpoint()
}
