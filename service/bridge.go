package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Bridge struct {
	bridges []*domain.Bridge
}

func NewBridge(cfg *config.Config, endpointREPO domain.EndpointRepositoryPort) *Bridge {
	var bridges []*domain.Bridge

	for _, bridge := range cfg.Bridges {
		for _, endpoint := range bridge.Endpoints {
			if _, err := endpointREPO.Get(endpoint); err != nil {
				log.Fatalln("service.NewBridge:", err)
			}
		}

		filters := make([]domain.Filter, len(bridge.Filters))
		for i := range bridge.Filters {
			filters[i] = domain.NewFilter(bridge.Filters[i].To, bridge.Filters[i].From, bridge.Filters[i].ToRegex, bridge.Filters[i].FromRegex)
		}

		bridges = append(bridges, domain.NewBridge(bridge.Name, bridge.Endpoints, bridge.OnlyText, bridge.OnlyAttachments, filters))
	}

	return &Bridge{bridges: bridges}
}

func (b *Bridge) ListByMessage(msg *domain.Message) []*domain.Bridge {
	var bridges []*domain.Bridge
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, bridge)
	}

	return bridges
}
