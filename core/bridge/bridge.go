package bridge

import "github.com/ItsNotGoodName/smtpbridge/core/envelope"

type (
	Bridge struct {
		Filter    []Filter
		Endpoints []string
	}

	CreateBridgeRequest struct {
		Filters   []CreateFilterRequest
		Endpoints []string
	}

	CreateFilterRequest struct {
		From          string
		FromRegex     string
		To            string
		ToRegex       string
		MatchTemplate string
	}

	Service interface {
		CreateBridge(req *CreateBridgeRequest) error
		ListBridge() []Bridge
	}
)

func NewBridge(filter []Filter, endpoints []string) Bridge {
	return Bridge{
		Filter:    filter,
		Endpoints: endpoints,
	}
}

func (b Bridge) Match(env *envelope.Envelope) (int, bool) {
	for i, filter := range b.Filter {
		if filter.Match(env) {
			return i, true
		}
	}

	return -1, false
}
