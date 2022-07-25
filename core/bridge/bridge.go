package bridge

type (
	Bridge struct {
		Filter    Filter
		Endpoints []string
	}

	CreateBridgeRequest struct {
		Name          string
		From          string
		FromRegex     string
		To            string
		ToRegex       string
		MatchTemplate string
		Endpoints     []string
	}

	Service interface {
		CreateBridge(req *CreateBridgeRequest) error
		ListBridge() []Bridge
	}
)

func NewBridge(filter Filter, endpoints []string) Bridge {
	return Bridge{
		Filter:    filter,
		Endpoints: endpoints,
	}
}
