package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Bridge struct {
	bridges []*core.Bridge
}

func NewBridge(cfg *config.Config, endpointREPO core.EndpointRepositoryPort) *Bridge {
	var bridges []*core.Bridge

	for _, bridge := range cfg.Bridges {
		for _, endpoint := range bridge.Endpoints {
			if _, err := endpointREPO.Get(endpoint); err != nil {
				log.Fatalln("service.NewBridge:", err)
			}
		}

		filters := make([]core.Filter, len(bridge.Filters))
		for i := range bridge.Filters {
			filters[i] = core.NewFilter(bridge.Filters[i].To, bridge.Filters[i].From, bridge.Filters[i].ToRegex, bridge.Filters[i].FromRegex)
		}

		bridges = append(bridges, core.NewBridge(bridge.Name, bridge.Endpoints, bridge.OnlyText, bridge.OnlyAttachments, filters))
	}

	return &Bridge{bridges: bridges}
}

func (b *Bridge) ListByMessage(msg *core.Message) []*core.Bridge {
	var bridges []*core.Bridge
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, bridge)
	}

	return bridges
}
